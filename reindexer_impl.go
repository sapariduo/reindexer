package reindexer

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/sapariduo/reindexer/bindings"
	"github.com/sapariduo/reindexer/cjson"
	"github.com/sapariduo/reindexer/dsl"
)

type reindexerNamespace struct {
	cacheItems    map[int]cacheItem
	cacheLock     sync.RWMutex
	joined        map[string][]int
	indexes       []bindings.IndexDef
	rtype         reflect.Type
	deepCopyIface bool
	name          string
	opts          NamespaceOptions
	cjsonState    cjson.State
	nsHash        int
	opened        bool
}

// reindexerImpl The reindxer state struct
type reindexerImpl struct {
	lock          sync.RWMutex
	ns            map[string]*reindexerNamespace
	storagePath   string
	binding       bindings.RawBinding
	debugLevels   map[string]int
	nsHashCounter int
	status        error
}

type cacheItem struct {
	item interface{}
	// version of item, for cachins
	version int
}

// NewReindexImpl Create new instanse of Reindexer DB
// Returns pointer to created instance
func newReindexImpl(dsn string, options ...interface{}) *reindexerImpl {

	if dsn == "builtin" {
		dsn += "://"
	}

	u, err := url.Parse(dsn)
	if err != nil {
		panic(fmt.Errorf("Can't parse DB DSN '%s'", dsn))
	}

	binding := bindings.GetBinding(u.Scheme)
	if binding == nil {
		panic(fmt.Errorf("Reindex binding '%s' is not available, can't create DB", u.Scheme))
	}

	binding = binding.Clone()
	rx := &reindexerImpl{
		ns:      make(map[string]*reindexerNamespace, 100),
		binding: binding,
	}

	if err = binding.Init(u, options...); err != nil {
		rx.status = err
	}

	if changing, ok := binding.(bindings.RawBindingChanging); ok {
		changing.OnChangeCallback(rx.resetCaches)
	}

	rx.registerNamespaceImpl(NamespacesNamespaceName, &NamespaceOptions{}, NamespaceDescription{})
	rx.registerNamespaceImpl(PerfstatsNamespaceName, &NamespaceOptions{}, NamespacePerfStat{})
	rx.registerNamespaceImpl(MemstatsNamespaceName, &NamespaceOptions{}, NamespaceMemStat{})
	rx.registerNamespaceImpl(QueriesperfstatsNamespaceName, &NamespaceOptions{}, QueryPerfStat{})
	rx.registerNamespaceImpl(ConfigNamespaceName, &NamespaceOptions{}, DBConfigItem{})

	return rx
}

// getStatus will return current db status
func (db *reindexerImpl) getStatus(ctx context.Context) bindings.Status {
	status := db.binding.Status(ctx)
	status.Err = db.status
	return status
}

// setLogger sets logger interface for output reindexer logs
func (db *reindexerImpl) setLogger(log Logger) {
	if log != nil {
		logger = log
		db.binding.EnableLogger(log)
	} else {
		logger = &nullLogger{}
		db.binding.DisableLogger()
	}
}

// ping checks connection with reindexer
func (db *reindexerImpl) ping(ctx context.Context) error {
	return db.binding.Ping(ctx)
}

func (db *reindexerImpl) close() {
	if err := db.binding.Finalize(); err != nil {
		panic(err)
	}
}

// openNamespace Open or create new namespace and indexes based on passed struct.
// IndexDef fields of struct are marked by `reindex:` tag
func (db *reindexerImpl) openNamespace(ctx context.Context, namespace string, opts *NamespaceOptions, s interface{}) (err error) {
	namespace = strings.ToLower(namespace)
	if err = db.registerNamespaceImpl(namespace, opts, s); err != nil {
		panic(err)
	}

	ns, err := db.getNS(namespace)
	if err != nil {
		return err
	}

	for retry := 0; retry < 2; retry++ {
		if err = db.binding.OpenNamespace(ctx, namespace, opts.enableStorage, opts.dropOnFileFormatError); err != nil {
			break
		}

		for _, indexDef := range ns.indexes {
			if err = db.binding.AddIndex(ctx, namespace, indexDef); err != nil {
				break
			}
		}

		if err != nil {
			rerr, ok := err.(bindings.Error)
			if ok && rerr.Code() == bindings.ErrConflict && opts.dropOnIndexesConflict {
				db.binding.DropNamespace(ctx, namespace)
				continue
			}
			db.binding.CloseNamespace(ctx, namespace)
			break
		}

		break
	}

	return err
}

// RegisterNamespace Register go type against namespace. There are no data and indexes changes will be performed
func (db *reindexerImpl) registerNamespace(namespace string, opts *NamespaceOptions, s interface{}) (err error) {
	namespace = strings.ToLower(namespace)
	return db.registerNamespaceImpl(namespace, opts, s)
}

// registerNamespace Register go type against namespace. There are no data and indexes changes will be performed
func (db *reindexerImpl) registerNamespaceImpl(namespace string, opts *NamespaceOptions, s interface{}) (err error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	namespace = strings.ToLower(namespace)

	db.lock.Lock()
	defer db.lock.Unlock()

	oldNs, ok := db.ns[namespace]
	if ok {
		// Ns exists, and have different type
		if oldNs.rtype.Name() != t.Name() {
			return errNsExists
		}
		// Ns exists, and have the same type.
		return nil
	}
	haveDeepCopy := false

	if !opts.disableObjCache {
		var copier DeepCopy
		copier, haveDeepCopy = reflect.New(t).Interface().(DeepCopy)
		if haveDeepCopy {
			cpy := copier.DeepCopy()
			cpyType := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(cpy)).Interface())
			if cpyType != reflect.TypeOf(s) {
				return ErrDeepCopyType
			}
		}
	}

	ns := &reindexerNamespace{
		cacheItems:    make(map[int]cacheItem, 100),
		rtype:         t,
		name:          namespace,
		joined:        make(map[string][]int),
		opts:          *opts,
		cjsonState:    cjson.NewState(),
		deepCopyIface: haveDeepCopy,
		nsHash:        db.nsHashCounter,
		opened:        false,
	}

	validator := cjson.Validator{}

	if err = validator.Validate(s); err != nil {
		return err
	}
	if ns.indexes, err = parseIndex(namespace, ns.rtype, &ns.joined); err != nil {
		return err
	}

	db.nsHashCounter++
	db.ns[namespace] = ns
	return nil
}

// dropNamespace - drop whole namespace from DB
func (db *reindexerImpl) dropNamespace(ctx context.Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	db.lock.Lock()
	delete(db.ns, namespace)
	db.lock.Unlock()

	return db.binding.DropNamespace(ctx, namespace)
}

// truncateNamespace - delete all items from namespace
func (db *reindexerImpl) truncateNamespace(ctx context.Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	return db.binding.TruncateNamespace(ctx, namespace)
}

// closeNamespace - close namespace, but keep storage
func (db *reindexerImpl) closeNamespace(ctx context.Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	db.lock.Lock()
	delete(db.ns, namespace)
	db.lock.Unlock()

	return db.binding.CloseNamespace(ctx, namespace)
}

// upsert (Insert or Update) item to index
// Item must be the same type as item passed to OpenNamespace, or []byte with json
func (db *reindexerImpl) upsert(ctx context.Context, namespace string, item interface{}, precepts ...string) error {
	_, err := db.modifyItem(ctx, namespace, nil, item, nil, modeUpsert, precepts...)
	return err
}

// insert item to namespace by PK
// Item must be the same type as item passed to OpenNamespace, or []byte with json data
// Return 0, if no item was inserted, 1 if item was inserted
func (db *reindexerImpl) insert(ctx context.Context, namespace string, item interface{}, precepts ...string) (int, error) {
	return db.modifyItem(ctx, namespace, nil, item, nil, modeInsert, precepts...)
}

// update item to namespace by PK
// Item must be the same type as item passed to OpenNamespace, or []byte with json data
// Return 0, if no item was updated, 1 if item was updated
func (db *reindexerImpl) update(ctx context.Context, namespace string, item interface{}, precepts ...string) (int, error) {
	return db.modifyItem(ctx, namespace, nil, item, nil, modeUpdate, precepts...)
}

// delete - remove single item from namespace by PK
// Item must be the same type as item passed to OpenNamespace, or []byte with json data
func (db *reindexerImpl) delete(ctx context.Context, namespace string, item interface{}, precepts ...string) error {
	_, err := db.modifyItem(ctx, namespace, nil, item, nil, modeDelete, precepts...)
	return err
}

// configureIndex - congigure index.
// config argument must be struct with index configuration
// Deprecated: Use UpdateIndex instead.
func (db *reindexerImpl) configureIndex(ctx context.Context, namespace, index string, config interface{}) error {
	nsDef, err := db.describeNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	index = strings.ToLower(index)
	for _, iDef := range nsDef.Indexes {
		if strings.ToLower(iDef.Name) == index {
			iDef.Config = config
			return db.binding.UpdateIndex(ctx, namespace, bindings.IndexDef(iDef.IndexDef))
		}
	}
	return fmt.Errorf("rq: Index '%s' not found in namespace %s", index, namespace)
}

// addIndex - add index.
func (db *reindexerImpl) addIndex(ctx context.Context, namespace string, indexDef ...IndexDef) error {
	for _, index := range indexDef {
		if err := db.binding.AddIndex(ctx, namespace, bindings.IndexDef(index)); err != nil {
			return err
		}
	}
	return nil
}

// updateIndex - update index.
func (db *reindexerImpl) updateIndex(ctx context.Context, namespace string, indexDef IndexDef) error {
	return db.binding.UpdateIndex(ctx, namespace, bindings.IndexDef(indexDef))
}

// dropIndex - drop index.
func (db *reindexerImpl) dropIndex(ctx context.Context, namespace, index string) error {
	return db.binding.DropIndex(ctx, namespace, index)
}

func loglevelToString(logLevel int) string {
	switch logLevel {
	case INFO:
		return "info"
	case TRACE:
		return "trace"
	case ERROR:
		return "error"
	case WARNING:
		return "warning"
	case 0:
		return "none"
	default:
		return ""
	}
}

// setDefaultQueryDebug sets default debug level for queries to namespaces
func (db *reindexerImpl) setDefaultQueryDebug(ctx context.Context, namespace string, level int) error {
	citem := &DBConfigItem{Type: "namespaces"}
	item, err := db.query(ConfigNamespaceName).WhereString("type", EQ, "namespaces").ExecCtx(ctx).FetchOne()
	if err != nil {
		return err
	}

	citem = item.(*DBConfigItem)

	found := false

	if citem.Namespaces == nil {
		namespaces := make([]DBNamespacesConfig, 0, 1)
		citem.Namespaces = &namespaces
	}

	for i := range *citem.Namespaces {
		if (*citem.Namespaces)[i].Namespace == namespace {
			(*citem.Namespaces)[i].LogLevel = loglevelToString(level)
			found = true
		}
	}
	if !found {
		*citem.Namespaces = append(*citem.Namespaces, DBNamespacesConfig{Namespace: namespace, LogLevel: loglevelToString(level)})
	}
	return db.upsert(ctx, ConfigNamespaceName, citem)
}

// query Create new Query for building request
func (db *reindexerImpl) query(namespace string) *Query {
	return newQuery(db, namespace, nil)
}

func (db *reindexerImpl) queryTx(namespace string, tx *Tx) *Query {
	return newQuery(db, namespace, tx)
}

// execSQL make query to database. Query is a SQL statement.
// Return Iterator.
func (db *reindexerImpl) execSQL(ctx context.Context, query string) *Iterator {
	namespace := getQueryNamespace(query)
	result, nsArray, err := db.prepareSQL(ctx, namespace, query, false)
	if err != nil {
		return errIterator(err)
	}
	iter := newIterator(ctx, nil, result, nsArray, nil, nil, nil)
	return iter
}

// execSQLToJSON make query to database. Query is a SQL statement.
// Return JSONIterator.
func (db *reindexerImpl) execSQLToJSON(ctx context.Context, query string) *JSONIterator {
	namespace := getQueryNamespace(query)
	result, _, err := db.prepareSQL(ctx, namespace, query, true)
	if err != nil {
		return errJSONIterator(err)
	}
	defer result.Free()
	json, jsonOffsets, explain, err := db.rawResultToJson(result.GetBuf(), namespace, "total", nil, nil)
	if err != nil {
		return errJSONIterator(err)
	}
	return newJSONIterator(ctx, nil, json, jsonOffsets, explain)
}

func getQueryNamespace(query string) string {
	// TODO: do not parse query string twice in go and cpp
	namespace := ""
	querySlice := strings.Fields(strings.ToLower(query))

	for i := range querySlice {
		if querySlice[i] == "from" && i+1 < len(querySlice) {
			namespace = querySlice[i+1]
			break
		}
	}
	return namespace
}

// beginTx - start update transaction
func (db *reindexerImpl) beginTx(ctx context.Context, namespace string) (*Tx, error) {
	return newTx(db, namespace, ctx)
}

// mustBeginTx - start update transaction, panic on error
func (db *reindexerImpl) mustBeginTx(ctx context.Context, namespace string) *Tx {
	tx, err := newTx(db, namespace, ctx)
	if err != nil {
		panic(err)
	}
	return tx
}

func (db *reindexerImpl) queryFrom(d dsl.DSL) (*Query, error) {
	if d.Namespace == "" {
		return nil, ErrEmptyNamespace
	}

	q := db.query(d.Namespace).Offset(d.Offset)

	if d.Explain {
		q.Explain()
	}

	if d.Limit > 0 {
		q.Limit(d.Limit)
	}

	if d.Distinct != "" {
		q.Distinct(d.Distinct)
	}
	if d.Sort.Field != "" {
		q.Sort(d.Sort.Field, d.Sort.Desc, d.Sort.Values...)
	}

	for _, filter := range d.Filters {
		if filter.Field == "" {
			return nil, ErrEmptyFieldName
		}
		if filter.Value == nil {
			continue
		}

		cond, err := GetCondType(filter.Cond)
		if err != nil {
			return nil, err
		}

		switch strings.ToUpper(filter.Op) {
		case "":
			q.Where(filter.Field, cond, filter.Value)
		case "NOT":
			q.Not().Where(filter.Field, cond, filter.Value)
		default:
			return nil, ErrOpInvalid
		}
	}

	return q, nil
}

// GetStats Get local thread reindexer usage stats
// Deprecated: Use SELECT * FROM '#perfstats' to get performance statistics.
func (db *reindexerImpl) getStats() bindings.Stats {
	log.Println("Deprecated function reindexer.GetStats call. Use SELECT * FROM '#perfstats' to get performance statistics")
	return bindings.Stats{}
}

// ResetStats Reset local thread reindexer usage stats
// Deprecated: no longer used.
func (db *reindexerImpl) resetStats() {
}

// enableStorage enables persistent storage of data
// Deprecated: storage path should be passed as DSN part to reindexer.NewReindex (""), e.g. reindexer.NewReindexer ("builtin:///tmp/reindex").
func (db *reindexerImpl) enableStorage(ctx context.Context, storagePath string) error {
	log.Println("Deprecated function reindexer.EnableStorage call")
	return db.binding.EnableStorage(ctx, storagePath)
}
