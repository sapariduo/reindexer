package reindexer

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/sapariduo/reindexer"
	"github.com/sapariduo/reindexer/dsl"
)

type sortDistinctOptions struct {
	SortIndexes     []string
	DistinctIndexes []string
	TestComposite   bool
}

type IndexesTestCase struct {
	Name      string
	Namespace string
	Options   sortDistinctOptions
	Item      interface{}
}

// TestItem common test case
type TestItem struct {
	Prices        []*TestJoinItem `reindex:"prices,,joined"`
	Pricesx       []*TestJoinItem `reindex:"pricesx,,joined"`
	ID            int             `reindex:"id,-"`
	Genre         int64           `reindex:"genre,tree"`
	Year          int             `reindex:"year,tree"`
	Packages      []int           `reindex:"packages,hash"`
	Name          string          `reindex:"name,tree"`
	Countries     []string        `reindex:"countries,tree"`
	Age           int             `reindex:"age,hash"`
	AgeLimit      int64           `json:"age_limit" reindex:"age_limit,hash,sparse"`
	CompanyName   string          `json:"company_name" reindex:"company_name,hash,sparse"`
	Address       string          `json:"address"`
	PostalCode    int             `json:"postal_code"`
	EmptyInt      int             `json:"empty_int,omitempty"`
	Description   string          `reindex:"description,fuzzytext"`
	Rate          float64         `reindex:"rate,tree"`
	ExchangeRate  float64         `json:"exchange_rate"`
	PollutionRate float32         `json:"pollution_rate"`
	IsDeleted     bool            `reindex:"isdeleted,-"`
	Actor         Actor           `reindex:"actor"`
	PricesIDs     []int           `reindex:"price_id"`
	LocationID    string          `reindex:"location"`
	EndTime       int             `reindex:"end_time,-"`
	StartTime     int             `reindex:"start_time,tree"`
	Tmp           string          `reindex:"tmp,-"`
	_             struct{}        `reindex:"id+tmp,,composite,pk"`
	_             struct{}        `reindex:"age+genre,,composite"`
	_             struct{}        `reindex:"location+rate,,composite"`
}

// TestItemIDOnly test case for non-indexed fields
type TestItemIDOnly struct {
	Prices        []*TestJoinItem `reindex:"prices,,joined"`
	Pricesx       []*TestJoinItem `reindex:"pricesx,,joined"`
	ID            int             `reindex:"id,-"`
	Genre         int64           `json:"genre"`
	Year          int             `json:"year"`
	Packages      []int           `json:"packages"`
	Name          string          `json:"name"`
	Countries     []string        `json:"countries"`
	Age           int             `json:"age"`
	AgeLimit      int64           `json:"age_limit"`    // reindex:"age_limit,hash,sparse"`
	CompanyName   string          `json:"company_name"` // reindex:"company_name,hash,sparse"`
	Address       string          `json:"address"`
	PostalCode    int             `json:"postal_code"`
	Description   string          `json:"description"`
	Rate          float64         `json:"rate"`
	ExchangeRate  float64         `json:"exchange_rate"`
	PollutionRate float32         `json:"pollution_rate"`
	IsDeleted     bool            `json:"isdeleted"`
	Actor         Actor           `reindex:"actor"`
	PricesIDs     int             `json:"price_id"`
	LocationID    string          `json:"location"`
	EndTime       int             `json:"end_time"`
	StartTime     int             `json:"start_time"`
	Tmp           string          `reindex:"tmp,-"`
	_             struct{}        `reindex:"id+tmp,,composite,pk"`
	_             struct{}        `reindex:"age+genre,,composite"`
	_             struct{}        `reindex:"location+rate,,composite"`
}

// TestItemWithSparse test case for sparse indexes
type TestItemWithSparse struct {
	Prices        []*TestJoinItem `reindex:"prices,,joined"`
	Pricesx       []*TestJoinItem `reindex:"pricesx,,joined"`
	ID            int             `reindex:"id,-"`
	Genre         int64           `reindex:"genre,tree,sparse"`
	Year          int             `reindex:"year,tree,sparse"`
	Packages      []int           `reindex:"packages,hash,sparse"`
	Name          string          `reindex:"name,tree,sparse"`
	Countries     []string        `reindex:"countries,tree,sparse"`
	Age           int             `reindex:"age,hash,sparse"`
	AgeLimit      int64           `json:"age_limit" reindex:"age_limit,hash,sparse"`
	CompanyName   string          `json:"company_name" reindex:"company_name,hash,sparse"`
	Address       string          `json:"address"`
	PostalCode    int             `json:"postal_code"`
	Description   string          `reindex:"description,fuzzytext"`
	Rate          float64         `reindex:"rate,tree,sparse"`
	ExchangeRate  float64         `json:"exchange_rate"`
	PollutionRate float32         `json:"pollution_rate"`
	IsDeleted     bool            `reindex:"isdeleted,-"`
	Actor         Actor           `reindex:"actor"`
	PricesIDs     []int           `reindex:"price_id"`
	LocationID    string          `reindex:"location"`
	EndTime       int             `reindex:"end_time,-"`
	StartTime     int             `reindex:"start_time,tree"`
	Tmp           string          `reindex:"tmp,-"`
	_             struct{}        `reindex:"id+tmp,,composite,pk"`
	_             struct{}        `reindex:"age+genre,,composite"`
	_             struct{}        `reindex:"location+rate,,composite"`
}

type TestItemSimple struct {
	ID    int    `reindex:"id,,pk"`
	Year  int    `reindex:"year,tree"`
	Name  string `reindex:"name"`
	Phone string
}

type TestItemCustom struct {
	Actor             Actor `reindex:"actor"`
	Genre             int64
	Name              string `reindex:"name"`
	CustomUniqueField int
	ID                int `reindex:"id,,pk"`
	Year              int `reindex:"year,tree"`
}

type TestItemSimpleCmplxPK struct {
	ID    int32    `reindex:"id,-"`
	Year  int32    `reindex:"year,tree"`
	Name  string   `reindex:"name"`
	SubID string   `reindex:"subid,-"`
	_     struct{} `reindex:"id+subid,,composite,pk"`
}

type objectType struct {
	Name string `reindex:"name"`
	Age  int    `reindex:"age"`
	Rate int    `reindex:"rate"`
	Hash int64  `reindex:"hash"`
}

type TestItemObjectArray struct {
	ID        int          `reindex:"id,,pk"`
	Code      int64        `reindex:"code"`
	IsEnabled bool         `reindex:"is_enabled"`
	Desc      string       `reindex:"desc"`
	Objects   []objectType `reindex:"objects"`
	MainObj   objectType   `reindex:"main_obj"`
	Size      int          `reindex:"size"`
	Hash      int64        `reindex:"hash"`
	Salt      int          `reindex:"salt"`
	IsUpdated bool         `reindex:"is_updated"`
}

type TestItemNestedPK struct {
	PrimaryID int      `reindex:"primary_id"`
	Nested    TestItem `json:"nested"`
}

func init() {
	tnamespaces["test_items"] = TestItem{}
	tnamespaces["test_items_id_only"] = TestItemIDOnly{}
	tnamespaces["test_items_with_sparse"] = TestItemWithSparse{}

	tnamespaces["test_items_simple"] = TestItemSimple{}
	tnamespaces["test_items_simple_cmplx_pk"] = TestItemSimpleCmplxPK{}
	tnamespaces["test_items_not"] = TestItemSimple{}
	tnamespaces["test_items_delete_query"] = TestItem{}
	tnamespaces["test_items_update_query"] = TestItem{}
}

func FillTestItemsForNot() {
	tx := newTestTx(DB, "test_items_not")

	if err := tx.Upsert(&TestItemSimple{
		ID:   1,
		Year: 2001,
		Name: "blabla",
	}); err != nil {
		panic(err)
	}
	if err := tx.Upsert(&TestItemSimple{
		ID:   2,
		Year: 2002,
		Name: "sss",
	}); err != nil {
		panic(err)
	}
	tx.MustCommit()

}

func newTestItem(id int, pkgsCount int) interface{} {
	startTime := rand.Int() % 50000
	return &TestItem{
		ID:            mkID(id),
		Year:          rand.Int()%50 + 2000,
		Genre:         int64(rand.Int() % 50),
		Name:          randString(),
		Age:           rand.Int() % 5,
		AgeLimit:      int64(rand.Int()%10 + 40),
		CompanyName:   randString(),
		Address:       randString(),
		PostalCode:    randPostalCode(),
		Description:   randString(),
		Packages:      randIntArr(10, 10000, 50),
		Rate:          float64(rand.Int()%100) / 10,
		ExchangeRate:  rand.Float64(),
		PollutionRate: rand.Float32(),
		IsDeleted:     rand.Int()%2 != 0,
		PricesIDs:     randIntArr(10, 7000, 50),
		LocationID:    randLocation(),
		StartTime:     startTime,
		EndTime:       startTime + (rand.Int()%5)*1000,
		Actor: Actor{
			Name: randString(),
		},
	}
}

func newTestItemIDOnly(id int, pkgsCount int) interface{} {
	startTime := rand.Int() % 50000
	return &TestItemIDOnly{
		ID:            mkID(id),
		Year:          rand.Int()%50 + 2000,
		Genre:         int64(rand.Int() % 50),
		Name:          randString(),
		Age:           rand.Int() % 5,
		AgeLimit:      int64(rand.Int()%10 + 40),
		CompanyName:   randString(),
		Address:       randString(),
		PostalCode:    randPostalCode(),
		Description:   randString(),
		Packages:      randIntArr(pkgsCount, 10000, 50),
		Rate:          float64(rand.Int()%100) / 10,
		ExchangeRate:  rand.Float64(),
		PollutionRate: rand.Float32(),
		IsDeleted:     rand.Int()%2 != 0,
		PricesIDs:     rand.Int() % 100, //randIntArr(10, 7000, 50),
		LocationID:    randLocation(),
		StartTime:     startTime,
		EndTime:       startTime + (rand.Int()%5)*1000,
		Actor: Actor{
			Name: randString(),
		},
	}
}

func newTestItemWithSparse(id int, pkgsCount int) interface{} {
	startTime := rand.Int() % 50000
	return &TestItemWithSparse{
		ID:            mkID(id),
		Year:          rand.Int()%50 + 2000,
		Genre:         int64(rand.Int() % 50),
		Name:          randString(),
		Age:           rand.Int() % 5,
		AgeLimit:      int64(rand.Int()%10 + 40),
		CompanyName:   randString(),
		Address:       randString(),
		PostalCode:    randPostalCode(),
		Description:   randString(),
		Packages:      randIntArr(pkgsCount, 10000, 50),
		Rate:          float64(rand.Int()%100) / 10,
		ExchangeRate:  rand.Float64(),
		PollutionRate: rand.Float32(),
		IsDeleted:     rand.Int()%2 != 0,
		PricesIDs:     randIntArr(10, 7000, 50),
		LocationID:    randLocation(),
		StartTime:     startTime,
		EndTime:       startTime + (rand.Int()%5)*1000,
		Actor: Actor{
			Name: randString(),
		},
	}
}

func newTestItemNestedPK(id int, pkgsCount int) *TestItemNestedPK {
	return &TestItemNestedPK{
		PrimaryID: mkID(id),
		Nested:    *newTestItem(id+10, pkgsCount).(*TestItem),
	}
}

func newTestItemObjectArray(id int, arrSize int) *TestItemObjectArray {
	arr := make([]objectType, 0, arrSize)

	for i := 0; i < arrSize; i++ {
		arr = append(arr, objectType{
			Name: randString(),
			Age:  rand.Int() % 50,
			Rate: rand.Int() % 10,
			Hash: rand.Int63() % 1000000 >> 1,
		})
	}

	return &TestItemObjectArray{
		ID:        id,
		Code:      rand.Int63() % 10000 >> 1,
		IsEnabled: (rand.Int() % 2) == 0,
		Desc:      randString(),
		Objects:   arr,
		MainObj: objectType{
			Name: randString(),
			Age:  rand.Int() % 50,
			Rate: rand.Int() % 10,
			Hash: rand.Int63() % 1000000 >> 1,
		},
		Size:      rand.Int() % 200000,
		Hash:      rand.Int63() % 1000000 >> 1,
		Salt:      rand.Int() % 100000,
		IsUpdated: (rand.Int() % 2) == 0,
	}
}

func FillTestItemsTx(start int, count int, pkgsCount int, tx *txTest) {
	for i := 0; i < count; i++ {
		testItem := newTestItem(start+i, pkgsCount)
		if err := tx.Upsert(testItem); err != nil {
			panic(err)
		}
	}
}

func FillTestItems(ns string, start int, count int, pkgsCount int) {
	tx := newTestTx(DB, ns)
	FillTestItemsTx(start, count, pkgsCount, tx)
	tx.MustCommit()
}

type testItemsCreator func(int, int) interface{}

func FillTestItemsTxWithFunc(start, count, pkgsCount int, tx *txTest, fn testItemsCreator) {
	for i := 0; i < count; i++ {
		testItem := fn(start+i, pkgsCount)
		if err := tx.Upsert(testItem); err != nil {
			panic(err)
		}
	}
}

func FillTestItemsWithFunc(ns string, start int, count int, pkgsCount int, fn testItemsCreator) {
	tx := newTestTx(DB, ns)
	FillTestItemsTxWithFunc(start, count, pkgsCount, tx, fn)
	tx.MustCommit()
}

func TestQueries(t *testing.T) {
	log.Printf("Seeding indexed test items")

	FillTestItemsWithFunc("test_items", 0, 2500, 20, newTestItem)
	FillTestItemsWithFunc("test_items", 2500, 2500, 0, newTestItem)

	log.Println("Seeding non-indexed test items ...")

	FillTestItemsWithFunc("test_items_id_only", 0, 2500, 20, newTestItemIDOnly)
	FillTestItemsWithFunc("test_items_id_only", 2500, 2500, 0, newTestItemIDOnly)

	log.Println("Seeding sparse indexed test items ...")

	FillTestItemsWithFunc("test_items_with_sparse", 0, 2500, 20, newTestItemWithSparse)
	FillTestItemsWithFunc("test_items_with_sparse", 2500, 2500, 0, newTestItemWithSparse)

	//// FillTestItems("test_items", 0, 2500, 20)
	//// FillTestItems("test_items", 2500, 2500, 0)
	FillTestItemsForNot()

	if err := DB.CloseNamespace("test_items"); err != nil {
		panic(err)
	}

	if err := DB.OpenNamespace("test_items", reindexer.DefaultNamespaceOptions(), TestItem{}); err != nil {
		panic(err)
	}

	if err := DB.CloseNamespace("test_items_id_only"); err != nil {
		panic(err)
	}

	if err := DB.OpenNamespace("test_items_id_only", reindexer.DefaultNamespaceOptions(), TestItemIDOnly{}); err != nil {
		panic(err)
	}

	if err := DB.CloseNamespace("test_items_with_sparse"); err != nil {
		panic(err)
	}

	if err := DB.OpenNamespace("test_items_with_sparse", reindexer.DefaultNamespaceOptions(), TestItemWithSparse{}); err != nil {
		panic(err)
	}

	CheckTestItemsJsonQueries()

	CheckAggregateQueries()

	log.Println("Start queries tests ...")

	CheckTestItemsQueries()
	CheckTestItemsSQLQueries()
	CheckTestItemsDSLQueries()

	// Delete test
	tx := newTestTx(DB, "test_items")
	for i := 0; i < 4000; i++ {
		if err := tx.Delete(TestItem{ID: mkID(i)}); err != nil {
			panic(err)
		}
	}
	// Check insert after delete
	FillTestItemsTx(0, 500, 0, tx)
	//Check second update
	FillTestItemsTx(0, 1000, 5, tx)

	for i := 0; i < 5000; i++ {
		tx.Delete(TestItem{ID: mkID(i)})
	}

	// Stress test delete & update & insert
	for i := 0; i < 5000; i++ {
		tx.Delete(TestItem{ID: mkID(rand.Int() % 500)})
		FillTestItemsTx(rand.Int()%500, 1, 0, tx)
		tx.Delete(TestItem{ID: mkID(rand.Int() % 500)})
		FillTestItemsTx(rand.Int()%500, 1, 10, tx)
		if (i % 1000) == 0 {
			tx.Commit()
			tx = newTestTx(DB, "test_items")
		}
	}

	tx.Commit()

	FillTestItems("test_items", 3000, 1000, 0)
	FillTestItems("test_items", 4000, 500, 20)
	CheckTestItemsQueries()
	CheckTestItemsSQLQueries()
	CheckTestItemsDSLQueries()

}

type CompositeFacetResultItem struct {
	CompanyName string
	Rate        float64
	Count       int
}
type CompositeFacetResult []CompositeFacetResultItem

func (f CompositeFacetResult) Len() int      { return len(f) }
func (f CompositeFacetResult) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f CompositeFacetResult) Less(i, j int) bool {
	if f[i].Count == f[j].Count {
		if f[i].CompanyName == f[j].CompanyName {
			return f[i].Rate < f[j].Rate
		} else {
			return f[i].CompanyName > f[j].CompanyName
		}
	}
	return f[i].Count < f[j].Count
}

func CheckAggregateQueries() {

	facetLimit := 100
	facetOffset := 10
	ctx, cancel := context.WithCancel(context.Background())
	q := DB.Query("test_items")
	q.AggregateAvg("year")
	q.AggregateSum("YEAR")
	q.AggregateFacet("age")
	q.AggregateFacet("name")
	q.AggregateMin("age")
	q.AggregateMax("age")
	q.AggregateFacet("company_name", "rate").Limit(facetLimit).Offset(facetOffset).Sort("count", false).Sort("company_name", true).Sort("rate", false)
	q.AggregateFacet("packages")
	it := q.ExecCtx(ctx)
	cancel()
	if it.Error() != nil {
		panic(it.Error())
	}

	qcheck := DB.Query("test_items")
	res, err := qcheck.Exec().FetchAll()
	if err != nil {
		panic(err)
	}

	aggregations := it.AggResults()
	if len(aggregations) != 8 {
		panic(fmt.Errorf("%d != 8", len(aggregations)))
	}

	var sum float64
	ageFacet := make(map[int]int, 0)
	nameFacet := make(map[string]int, 0)
	packagesFacet := make(map[int]int, 0)
	type CompositeFacetItem struct {
		CompanyName string
		Rate        float64
	}
	compositeFacet := make(map[CompositeFacetItem]int, 0)
	ageMin, ageMax := 100000000, -10000000

	for _, it := range res {
		testItem := it.(*TestItem)
		sum += float64(testItem.Year)
		ageFacet[testItem.Age]++
		nameFacet[testItem.Name]++
		for _, pack := range testItem.Packages {
			packagesFacet[pack]++
		}
		compositeFacet[CompositeFacetItem{testItem.CompanyName, testItem.Rate}]++
		if testItem.Age > ageMax {
			ageMax = testItem.Age
		}
		if testItem.Age < ageMin {
			ageMin = testItem.Age
		}
	}

	var compositeFacetResult CompositeFacetResult
	for k, v := range compositeFacet {
		compositeFacetResult = append(compositeFacetResult, CompositeFacetResultItem{k.CompanyName, k.Rate, v})
	}
	sort.Sort(compositeFacetResult)
	compositeFacetResult = compositeFacetResult[min(facetOffset, len(compositeFacetResult)):min(facetOffset+facetLimit, len(compositeFacetResult))]

	if sum != aggregations[1].Value {
		panic(fmt.Errorf("%f != %f", sum, aggregations[1].Value))
	}
	if sum/float64(len(res)) != aggregations[0].Value {
		panic(fmt.Errorf("%f != %f,len=%d", sum/float64(len(res)), aggregations[0].Value, len(res)))
	}

	if len(aggregations[2].Fields) != 1 {
		panic(fmt.Errorf("%d != 1", len(aggregations[2].Fields)))
	}
	if aggregations[2].Fields[0] != "age" {
		panic(fmt.Errorf("%s != %s", aggregations[2].Fields[0], "age"))
	}
	if len(aggregations[2].Facets) != len(ageFacet) {
		panic(fmt.Errorf("%d != %d", len(aggregations[2].Facets), len(ageFacet)))
	}
	for _, facet := range aggregations[2].Facets {
		if len(facet.Values) != 1 {
			panic(fmt.Errorf("%d != 1", len(facet.Values)))
		}
		intVal, _ := strconv.Atoi(facet.Values[0])
		if count, ok := ageFacet[intVal]; ok != true || count != facet.Count {
			panic(fmt.Errorf("facet '%s' val '%s': %d != %d", aggregations[2].Fields[0], facet.Values[0], count, facet.Count))
		}
	}

	if len(aggregations[3].Fields) != 1 {
		panic(fmt.Errorf("%d != 1", len(aggregations[3].Fields)))
	}
	if aggregations[3].Fields[0] != "name" {
		panic(fmt.Errorf("%s != %s", aggregations[3].Fields[0], "name"))
	}
	for _, facet := range aggregations[3].Facets {
		if len(facet.Values) != 1 {
			panic(fmt.Errorf("%d != 1", len(facet.Values)))
		}
		if count, ok := nameFacet[facet.Values[0]]; ok != true || count != facet.Count {
			panic(fmt.Errorf("facet '%s' val '%s': %d != %d", aggregations[3].Fields[0], facet.Values[0], count, facet.Count))
		}
	}
	if ageMin != int(aggregations[4].Value) {
		panic(fmt.Errorf("%d != %f", ageMin, aggregations[4].Value))
	}
	if ageMax != int(aggregations[5].Value) {
		panic(fmt.Errorf("%d != %f", ageMax, aggregations[5].Value))
	}
	if len(aggregations[6].Fields) != 2 {
		panic(fmt.Errorf("%d != 1", len(aggregations[6].Fields)))
	}
	if aggregations[6].Fields[0] != "company_name" {
		panic(fmt.Errorf("%s != %s", aggregations[6].Fields[0], "company_name"))
	}
	if aggregations[6].Fields[1] != "rate" {
		panic(fmt.Errorf("%s != %s", aggregations[6].Fields[1], "rate"))
	}
	if len(compositeFacetResult) != len(aggregations[6].Facets) {
		panic(fmt.Errorf("Composite facet sizes differ: %d != %d", len(compositeFacetResult), len(aggregations[6].Facets)))
	}
	for i := 0; i < len(compositeFacetResult); i++ {
		if len(aggregations[6].Facets[i].Values) != 2 {
			panic(fmt.Errorf("%d != 2", len(aggregations[6].Facets[i].Values)))
		}
		rate, err := strconv.ParseFloat(aggregations[6].Facets[i].Values[1], 64)
		if err != nil {
			panic(err)
		}
		if compositeFacetResult[i].CompanyName != aggregations[6].Facets[i].Values[0] || compositeFacetResult[i].Rate != rate ||
			compositeFacetResult[i].Count != aggregations[6].Facets[i].Count {
			panic(fmt.Errorf("Facet 'company_name', 'rate' #%d {'%s', '%s': %d} != {'%s', '%f': %d}", i,
				aggregations[6].Facets[i].Values[0], aggregations[6].Facets[i].Values[1], aggregations[6].Facets[i].Count,
				compositeFacetResult[i].CompanyName, compositeFacetResult[i].Rate, compositeFacetResult[i].Count))
		}
	}
	if len(aggregations[7].Fields) != 1 {
		panic(fmt.Errorf("%d != 1", len(aggregations[7].Fields)))
	}
	if aggregations[7].Fields[0] != "packages" {
		panic(fmt.Errorf("%s != %s", aggregations[7].Fields[0], "packages"))
	}
	for _, facet := range aggregations[7].Facets {
		if len(facet.Values) != 1 {
			panic(fmt.Errorf("%d != 1", len(facet.Values)))
		}
		value, err := strconv.Atoi(facet.Values[0])
		if err != nil {
			panic(err)
		}
		if count, ok := packagesFacet[value]; ok != true || count != facet.Count {
			panic(fmt.Errorf("facet '%s' val '%s' (%d): %d != %d", aggregations[7].Fields[0], facet.Values[0], value, count, facet.Count))
		}
	}
}

func CheckTestItemsJsonQueries() {
	ctx, cancel := context.WithCancel(context.Background())
	json, _ := DB.Query("test_items").Select("ID", "Genre").Limit(3).ReqTotal("total_count").ExecToJson("test_items").FetchAll()
	//	fmt.Println(string(json))
	_ = json

	json2, _ := DB.Query("test_items").Select("ID", "Genre").Limit(3).ReqTotal("total_count").GetJsonCtx(ctx)
	//	fmt.Println(string(json2))
	cancel()
	_ = json2
	// TODO

}

func makeLikePattern(s string) string {
	runes := make([]rune, len(s))
	i := 0
	for _, rune := range s {
		if rand.Int()%4 == 0 {
			runes[i] = '_'
		} else {
			runes[i] = rune
		}
		i++
	}
	var result string
	if rand.Int()%4 == 0 {
		result += "%"
	}
	current := 0
	next := rand.Int() % (len(s) + 1)
	last := next
	for current < len(s) {
		if current < next {
			result += string(runes[current:next])
			last = next
			current = rand.Int()%(len(s)-last+1) + last
		}
		next = rand.Int()%(len(s)-current+1) + current
		if current > last || rand.Int()%4 == 0 {
			result += "%"
		}
	}
	if rand.Int()%4 == 0 {
		result += "%"
	}
	return result
}

func callQueriesSequence(namespace, distinct, sort string, desc, testComposite bool) {
	// Take items with single condition
	newTestQuery(DB, namespace).Where("genre", reindexer.EQ, rand.Int()%50).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("rate", reindexer.EQ, float32(rand.Int()%100)/10).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("age_limit", reindexer.EQ, int64((rand.Int()%10)+40)).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("postal_code", reindexer.EQ, randPostalCode()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("company_name", reindexer.EQ, randString()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("address", reindexer.EQ, randString()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("exchange_rate", reindexer.EQ, rand.Float64()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("pollution_rate", reindexer.EQ, rand.Float32()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("genre", reindexer.GT, rand.Int()%50).Distinct(distinct).Sort(sort, desc).Debug(reindexer.TRACE).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.GT, randString()).Distinct(distinct).Offset(21).Limit(50).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("rate", reindexer.GT, float32(rand.Int()%100)/10).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("age_limit", reindexer.GT, int64(40)).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("postal_code", reindexer.GT, randPostalCode()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("company_name", reindexer.GT, randString()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("address", reindexer.GT, randString()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("exchange_rate", reindexer.GT, rand.Float64()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("pollution_rate", reindexer.GT, rand.Float32()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("genre", reindexer.LT, rand.Int()%50).Distinct(distinct).Sort(sort, desc).Limit(100).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.LT, randString()).Offset(10).Limit(200).ExecAndVerify()
	newTestQuery(DB, namespace).Where("rate", reindexer.LT, float32(rand.Int()%100)/10).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("age_limit", reindexer.LT, int64(50)).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("postal_code", reindexer.LT, randPostalCode()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("company_name", reindexer.LT, randString()).ExecAndVerify()
	newTestQuery(DB, namespace).Where("address", reindexer.LT, randString()).ExecAndVerify()
	newTestQuery(DB, namespace).Where("exchange_rate", reindexer.LT, rand.Float64()).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("pollution_rate", reindexer.LT, rand.Float32()).Distinct(distinct).Sort(sort, desc).Limit(500).ExecAndVerify()

	newTestQuery(DB, namespace).Where("genre", reindexer.RANGE, []int{rand.Int() % 100, rand.Int() % 100}).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.RANGE, []string{randString(), randString()}).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("rate", reindexer.RANGE, []float32{float32(rand.Int()%100) / 10, float32(rand.Int()%100) / 10}).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("age_limit", reindexer.RANGE, []int64{40, 50}).ExecAndVerify()
	newTestQuery(DB, namespace).Where("company_name", reindexer.RANGE, []string{randString(), randString()}).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("name", reindexer.LIKE, makeLikePattern(randString())).ExecAndVerify()

	newTestQuery(DB, namespace).Where("packages", reindexer.SET, randIntArr(10, 10000, 50)).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("packages", reindexer.EMPTY, 0).Distinct(distinct).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("packages", reindexer.ANY, 0).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("isdeleted", reindexer.EQ, true).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).Sort("address", false).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).Sort("postal_code", true).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).Sort("age_limit", false).Sort(sort, desc).ExecAndVerify()
	newTestQuery(DB, namespace).Where("packages", reindexer.GE, 5).Where("price_id", reindexer.GE, 100).EqualPosition("packages", "price_id").ExecAndVerify()

	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).
		Sort("year", true).
		Sort("name", false).
		Sort("genre", true).
		Sort("age", false).
		Sort("age_limit", true).
		ExecAndVerify()

	newTestQuery(DB, namespace).Where("name", reindexer.EQ, randString()).Distinct(distinct).
		Sort("year", true).
		Sort("name", false).
		Sort("genre", true).
		Sort("age", false).
		Sort("age_limit", true).
		Offset(4).
		Limit(44).
		ExecAndVerify()

	// Complex queires
	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5). // composite index age+genre
		Where("age", reindexer.EQ, 3).
		Where("year", reindexer.GT, 2010).
		Where("age_limit", reindexer.LE, int64(50)).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("year", reindexer.GT, 2002).
		Where("genre", reindexer.EQ, 4). // composite index age+genre, and extra and + or conditions
		Where("age", reindexer.EQ, 3).
		Where("age_limit", reindexer.GE, int64(40)).
		Where("isdeleted", reindexer.EQ, true).Or().Where("year", reindexer.GT, 2001).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).Debug(reindexer.TRACE).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age", reindexer.SET, []int{1, 2, 3, 4}).
		Where("id", reindexer.EQ, mkID(rand.Int()%5000)).
		Where("tmp", reindexer.EQ, ""). // composite pk with store index
		Where("isdeleted", reindexer.EQ, true).Or().Where("year", reindexer.GT, 2001).Debug(reindexer.TRACE).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("genre", reindexer.SET, []int{5, 1, 7}).
		Where("year", reindexer.LT, 2010).Or().Where("genre", reindexer.EQ, 3).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).Or().Where("packages", reindexer.EMPTY, 0).Debug(reindexer.TRACE).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("genre", reindexer.SET, []int{5, 1, 7}).
		Where("year", reindexer.LT, 2010).Or().Where("packages", reindexer.ANY, 0).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).Debug(reindexer.TRACE).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("genre", reindexer.EQ, 5).Or().Where("genre", reindexer.EQ, 6).
		Where("year", reindexer.RANGE, []int{2001, 2020}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("year", reindexer.RANGE, []int{2001, 2020}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("actor.name", reindexer.EQ, randString()).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Not().Where("genre", reindexer.EQ, 5).
		Where("year", reindexer.RANGE, []int{2001, 2020}).
		Where("age_limit", reindexer.RANGE, []int64{40, 50}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		Not().Where("year", reindexer.RANGE, []int{2001, 2020}).
		Where("age_limit", reindexer.RANGE, []int64{40, 50}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Where("genre", reindexer.EQ, 5).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Or().Where("genre", reindexer.EQ, 5).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Not().OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Or().Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		Or().OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Or().Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		Or().OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		Not().Where("year", reindexer.RANGE, []int{2001, 2010}).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Or().Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Where("genre", reindexer.EQ, 5).
		Or().OpenBracket().
		Where("genre", reindexer.EQ, 4).
		Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		CloseBracket().
		Not().Where("year", reindexer.RANGE, []int{2001, 2010}).
		OpenBracket().
		Where("age_limit", reindexer.RANGE, []int64{40, 45}).
		Or().Where("packages", reindexer.SET, randIntArr(5, 10000, 50)).
		Not().OpenBracket().
		Or().OpenBracket().
		OpenBracket().
		Where("age", reindexer.SET, []int{1, 2, 3, 4}).
		Or().Where("id", reindexer.EQ, mkID(rand.Int()%5000)).
		Not().Where("tmp", reindexer.EQ, ""). // composite pk with store index
		CloseBracket().
		Or().OpenBracket().
		Where("tmp", reindexer.EQ, ""). // composite pk with store index
		Not().Where("isdeleted", reindexer.EQ, true).Or().Where("year", reindexer.GT, 2001).Debug(reindexer.TRACE).
		CloseBracket().
		CloseBracket().
		CloseBracket().
		CloseBracket().
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().Debug(reindexer.TRACE).
		Not().Where("genre", reindexer.EQ, 10).
		ExecAndVerify()

	newTestQuery(DB, "TEST_ITEMS_NOT").ReqTotal().
		Where("NAME", reindexer.EQ, "blabla").
		ExecAndVerify()

	newTestQuery(DB, "test_items_not").ReqTotal().
		Where("year", reindexer.EQ, 2002).
		ExecAndVerify()

	newTestQuery(DB, "test_items_not").ReqTotal().
		Where("YEAR", reindexer.EQ, 2002).
		Not().Where("name", reindexer.EQ, "blabla").
		ExecAndVerify()

	newTestQuery(DB, "TEST_ITEMS_NOT").ReqTotal().
		Where("name", reindexer.EQ, "blabla").
		Not().Where("year", reindexer.EQ, 2002).
		ExecAndVerify()

	newTestQuery(DB, "test_items_not").ReqTotal().
		Where("name", reindexer.EQ, "blabla").
		Not().Where("year", reindexer.EQ, 2001).
		ExecAndVerify()

	newTestQuery(DB, "test_items_not").ReqTotal().
		Where("year", reindexer.EQ, 2002).
		Not().Where("name", reindexer.EQ, "sss").
		ExecAndVerify()

	newTestQuery(DB, namespace).Where("end_time", reindexer.GT, 10000).Not().Where("genre", reindexer.EQ, 10).Distinct(distinct).Sort(sort, desc).ExecAndVerify()

	if !testComposite {
		return
	}
	compositeValues := []interface{}{[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)}}

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.EQ, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.LE, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.LT, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.GT, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.GE, compositeValues).
		ExecAndVerify()

	compositeValues = []interface{}{[]interface{}{randLocation(), float64(rand.Int()%100) / 10}}

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("location+rate", reindexer.EQ, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("location+rate", reindexer.GT, compositeValues).
		ExecAndVerify()

	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("location+rate", reindexer.LT, compositeValues).
		ExecAndVerify()

	compositeValues = []interface{}{
		[]interface{}{randLocation(), float64(rand.Int()%100) / 10},
		[]interface{}{randLocation(), float64(rand.Int()%100) / 10},
	}
	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("location+rate", reindexer.RANGE, compositeValues).
		ExecAndVerify()

	compositeValues = []interface{}{
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
		[]interface{}{rand.Int() % 10, int64(rand.Int() % 50)},
	}
	newTestQuery(DB, namespace).Distinct(distinct).Sort(sort, desc).ReqTotal().
		Where("age+genre", reindexer.SET, compositeValues).
		ExecAndVerify()
}

func CheckTestItemsQueries() {
	testCases := []IndexesTestCase{
		{
			Name:      "TEST WITH COMMON INDEXES",
			Namespace: "test_items",
			Options: sortDistinctOptions{
				SortIndexes:     []string{"", "NAME", "YEAR", "RATE", "RATE + (GENRE - 40) * ISDELETED"},
				DistinctIndexes: []string{"", "YEAR", "RATE"},
				TestComposite:   true,
			},
			Item: TestItem{},
		},
		{
			Name:      "TEST WITH ID ONLY INDEX",
			Namespace: "test_items_id_only",
			Options: sortDistinctOptions{
				SortIndexes:     []string{"", "name", "year", "rate", "rate + (genre - 40) * isdeleted"},
				DistinctIndexes: []string{"", "year", "rate"},
				TestComposite:   false,
			},
			Item: TestItemIDOnly{},
		},
		{
			Name:      "TEST WITH SPARSE INDEXES",
			Namespace: "test_items_with_sparse",
			Options: sortDistinctOptions{
				SortIndexes:     []string{"", "NAME", "YEAR", "RATE", "-ID + (END_TIME - START_TIME) / 100"},
				DistinctIndexes: []string{"", "YEAR", "RATE"},
				TestComposite:   false,
			},
			Item: TestItemWithSparse{},
		},
	}

	for _, testCase := range testCases {
		log.Println(testCase.Name)
		for _, desc := range []bool{true, false} {
			for _, sort := range testCase.Options.SortIndexes {
				for _, distinct := range testCase.Options.DistinctIndexes {
					log.Printf("\tDISTINCT '%s' SORT '%s' DESC %v\n", distinct, sort, desc)
					// Just take all items from namespace
					newTestQuery(DB, testCase.Namespace).Distinct(distinct).Sort(sort, desc).Limit(1).ExecAndVerify()
					callQueriesSequence(testCase.Namespace, distinct, sort, desc, testCase.Options.TestComposite)
				}
			}
		}
	}
}

func CheckTestItemsSQLQueries() {
	companyName := randString()
	if res, err := DB.ExecSQL("SELECT ID, company_name FROM test_ITEMS WHERE COMPANY_name > '" + companyName + "' ORDER BY YEAr DESC LIMIT 10000000 ; ").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("company_name", reindexer.GT, companyName).Sort("year", true).Verify(res, false)
	}

	if res, err := DB.ExecSQL("SELECT ID, postal_code FROM test_items WHERE postal_code > 121355 ORDER BY year DESC LIMIT 10000000 ; ").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("postal_code", reindexer.GT, 121355).Sort("year", true).Verify(res, false)
	}

	if res, err := DB.ExecSQL("SELECT ID,Year,Genre,age_limit FROM test_items WHERE YEAR > '2016' AND genre IN ('1',2,'3') and age_limit IN(40,'50',42,'47') ORDER BY year DESC LIMIT 10000000 ; ").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("year", reindexer.GT, 2016).Where("GENRE", reindexer.SET, []int{1, 2, 3}).Where("age_limit", reindexer.SET, []int64{40, 50, 42, 47}).Sort("YEAR", true).Verify(res, false)
	}

	if res, err := DB.ExecSQL("SELECT * FROM test_items WHERE YEAR <= '2016' OR genre < 5 or AGE_LIMIT >= 40 ORDER BY YEAR ASC ").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("year", reindexer.LE, 2016).Or().Where("genre", reindexer.LT, 5).Or().Where("age_limit", reindexer.GE, int64(40)).Sort("year", false).Verify(res, true)
	}

	if res, err := DB.ExecSQL("SELECT count(*), * FROM test_items WHERE year >= '2016' OR rate = '1.1' OR year RANGE (2010,2014) or AGE_LIMIT <= 50").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("year", reindexer.GE, 2016).Or().Where("RATE", reindexer.EQ, 1.1).Or().Where("YEAR", reindexer.RANGE, []int{2010, 2014}).Or().Where("age_limit", reindexer.LE, int64(50)).Verify(res, true)
	}

	likePattern := makeLikePattern(randString())
	if res, err := DB.ExecSQL("SELECT count(*), * FROM test_items WHERE year >= '2016' OR rate = '1.1' OR company_name LIKE '" + likePattern + "' or AGE_LIMIT <= 50").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("year", reindexer.GE, 2016).Or().Where("RATE", reindexer.EQ, 1.1).Or().Where("company_name", reindexer.LIKE, likePattern).Or().Where("age_limit", reindexer.LE, int64(50)).Verify(res, true)
	}

	if res, err := DB.ExecSQL("SELECT ID,'Actor.Name' FROM test_items WHERE 'actor.name' > 'bde'  LIMIT 10000000").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("actor.name", reindexer.GT, []string{"bde"}).Verify(res, false)
	}

	if res, err := DB.ExecSQL("SELECT count(*), * FROM test_items WHERE year >= '2016' OR (rate = '1.1' OR company_name LIKE '" + likePattern + "') and AGE_LIMIT <= 50").FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").Where("year", reindexer.GE, 2016).Or().OpenBracket().Where("RATE", reindexer.EQ, 1.1).Or().Where("company_name", reindexer.LIKE, likePattern).CloseBracket().Where("age_limit", reindexer.LE, int64(50)).Verify(res, true)
	}
}

func CheckTestItemsDSLQueries() {
	likePattern := makeLikePattern(randString())
	d := dsl.DSL{
		Namespace: "TEST_ITEMS",
		Filters: []dsl.Filter{
			{
				Field: "YEAR",
				Cond:  "GT",
				Value: "2016",
			},
			{
				Field: "GENRE",
				Cond:  "SET",
				Value: []string{"1", "2", "3"},
			},
			{
				Field: "PACKAGES",
				Cond:  "ANY",
				Value: 0,
			},
			{
				Field: "countries",
				Cond:  "EMPTY",
				Value: 0,
			},
			{
				Field: "isdeleted",
				Cond:  "EQ",
				Value: true,
			},
			{
				Field: "company_name",
				Cond:  "LIKE",
				Value: likePattern,
			},
		},
		Sort: dsl.Sort{
			Field: "YEAR",
			Desc:  true,
		},
	}

	if q, err := DB.QueryFrom(d); err != nil {
		panic(err)
	} else if res, err := q.Exec().FetchAll(); err != nil {
		panic(err)
	} else {
		newTestQuery(DB, "test_items").
			Where("year", reindexer.GT, 2016).
			Where("genre", reindexer.SET, []int{1, 2, 3}).
			Where("packages", reindexer.ANY, 0).
			Where("countries", reindexer.EMPTY, 0).
			Where("isdeleted", reindexer.EQ, true).
			Where("company_name", reindexer.LIKE, likePattern).
			Sort("year", true).
			Verify(res, true)
	}
}

func TestCanceledSelectQuery(t *testing.T) {
	if !strings.HasPrefix(DB.dsn, "builtin://") {
		return
	}

	FillTestItemsWithFunc("test_items", 0, 1000, 80, newTestItem)
	FillTestItemsWithFunc("test_items", 1000, 1000, 0, newTestItem)

	likePattern := makeLikePattern(randString())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := DB.WithContext(ctx).ExecSQL("SELECT count(*), * FROM test_items WHERE year >= '2016' OR (rate = '1.1' OR company_name LIKE '" + likePattern + "') and AGE_LIMIT <= 50").FetchAll()
	if err != context.Canceled {
		panic(fmt.Errorf("Canceled select request was executed"))
	}

	it := newTestQuery(DB, "test_items").Where("year", reindexer.GE, 2016).Or().OpenBracket().Where("RATE", reindexer.EQ, 1.1).Or().Where("company_name", reindexer.LIKE, likePattern).CloseBracket().Where("age_limit", reindexer.LE, int64(50)).ExecCtx(ctx)
	defer it.Close()
	if err != context.Canceled {
		panic(fmt.Errorf("Canceled select request was executed"))
	}
}

func TestDeleteQuery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := DB.UpsertCtx(ctx, "test_items_delete_query", newTestItem(1000, 5))
	cancel()
	if err != nil {
		panic(err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	count, err := DB.Query("test_items_delete_query").Where("id", reindexer.EQ, mkID(1000)).DeleteCtx(ctx)
	cancel()
	if err != nil {
		panic(err)
	}

	if count != 1 {
		panic(fmt.Errorf("Expected delete query return 1 item"))
	}

	ctx, cancel = context.WithCancel(context.Background())
	_, ok := DB.Query("test_items_delete_query").Where("id", reindexer.EQ, mkID(1000)).Get()
	cancel()
	if ok {
		panic(fmt.Errorf("Item was found after delete query, but will be deleted"))
	}

}
func TestCanceledDeleteQuery(t *testing.T) {
	if !strings.HasPrefix(DB.dsn, "builtin://") {
		return
	}

	err := DB.Upsert("test_items_delete_query", newTestItem(1000, 5))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = DB.Query("test_items_delete_query").Where("id", reindexer.EQ, mkID(1000)).DeleteCtx(ctx)
	if err != context.Canceled {
		panic(fmt.Errorf("Canceled delete request was executed"))
	}

	_, ok := DB.Query("test_items_delete_query").Where("id", reindexer.EQ, mkID(1000)).Get()
	if !ok {
		panic(fmt.Errorf("Item was deleted after canceled delete query"))
	}
}

func TestUpdateQuery(t *testing.T) {
	err := DB.Upsert("test_items_update_query", newTestItem(1000, 5))

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	it := DB.Query("test_items_update_query").Where("id", reindexer.EQ, mkID(1000)).
		Set("name", "hello").
		Set("empty_int", 1).
		Set("postal_code", 10).UpdateCtx(ctx)
	cancel()
	defer it.Close()
	if it.Error() != nil {
		panic(it.Error())
	}

	if it.Count() != 1 {
		panic(fmt.Errorf("Expected update query return 1 item"))
	}

	ctx, cancel = context.WithCancel(context.Background())
	f, ok := DB.Query("test_items_update_query").Where("id", reindexer.EQ, mkID(1000)).GetCtx(ctx)
	cancel()
	if !ok {
		panic(fmt.Errorf("Item was not found after update query"))
	}
	if f.(*TestItem).EmptyInt != 1 {
		panic(fmt.Errorf("Item have wrong value %d, shoud %d", f.(*TestItem).EmptyInt, 1))
	}
	if f.(*TestItem).PostalCode != 10 {
		panic(fmt.Errorf("Item have wrong value %d, shoud %d", f.(*TestItem).PostalCode, 10))
	}
	if f.(*TestItem).Name != "hello" {
		panic(fmt.Errorf("Item have wrong value %s, shoud %s", f.(*TestItem).Name, "hello"))
	}

}
func TestCanceledUpdateQuery(t *testing.T) {
	if !strings.HasPrefix(DB.dsn, "builtin://") {
		return
	}

	item := newTestItem(1000, 5)
	err := DB.Upsert("test_items_update_query", item)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = DB.UpsertCtx(ctx, "test_items_update_query", item)
	if err != context.Canceled {
		panic(fmt.Errorf("Canceled upsert request was executed"))
	}

	it := DB.Query("test_items_update_query").Where("id", reindexer.EQ, mkID(1000)).
		Set("name", "hello").
		Set("empty_int", 1).
		Set("postal_code", 10).UpdateCtx(ctx)
	defer it.Close()
	if it.Error() != context.Canceled {
		panic(fmt.Errorf("Canceled update request was executed"))
	}

	f, ok := DB.Query("test_items_update_query").Where("id", reindexer.EQ, mkID(1000)).Get()
	if !ok {
		panic(fmt.Errorf("Item was not found after update query"))
	}
	if f.(*TestItem).EmptyInt != item.(*TestItem).EmptyInt {
		panic(fmt.Errorf("Item have wrong value %d, should %d", f.(*TestItem).EmptyInt, item.(*TestItem).EmptyInt))
	}
	if f.(*TestItem).PostalCode != item.(*TestItem).PostalCode {
		panic(fmt.Errorf("Item have wrong value %d, should %d", f.(*TestItem).PostalCode, item.(*TestItem).PostalCode))
	}
	if f.(*TestItem).Name != item.(*TestItem).Name {
		panic(fmt.Errorf("Item have wrong value %s, should %s", f.(*TestItem).Name, item.(*TestItem).Name))
	}
}

func TestDeleteByPK(t *testing.T) {
	nsOpts := reindexer.DefaultNamespaceOptions()

	assertErrorMessage(t, DB.OpenNamespace("test_items_object_array", nsOpts, TestItemObjectArray{}), nil)
	for i := 1; i <= 30; i++ {
		assertErrorMessage(t, DB.Upsert("test_items_object_array", newTestItemObjectArray(i, rand.Int()%10)), nil)
	}

	assertErrorMessage(t, DB.MustBeginTx("test_items_object_array").Commit(), nil)
	assertErrorMessage(t, DB.CloseNamespace("test_items_object_array"), nil)
	assertErrorMessage(t, DB.OpenNamespace("test_items_object_array", nsOpts, TestItemObjectArray{}), nil)

	for i := 1; i <= 30; i++ {
		// specially create a complete item
		assertErrorMessage(t, DB.Delete("test_items_object_array", newTestItemObjectArray(i, rand.Int()%10)), nil)

		// check deletion result
		if item, found := DB.Query("test_items_object_array").WhereInt("id", reindexer.EQ, i).Get(); found {
			t.Fatalf("Item has not been deleted. < %+v > ", item)
		}
	}

	assertErrorMessage(t, DB.OpenNamespace("test_item_delete", nsOpts, TestItem{}), nil)
	for i := 1; i <= 30; i++ {
		assertErrorMessage(t, DB.Upsert("test_item_delete", newTestItem(i, rand.Int()%20)), nil)
	}

	assertErrorMessage(t, DB.MustBeginTx("test_item_delete").Commit(), nil)
	assertErrorMessage(t, DB.CloseNamespace("test_item_delete"), nil)
	assertErrorMessage(t, DB.OpenNamespace("test_item_delete", nsOpts, TestItem{}), nil)

	for i := 1; i <= 30; i++ {
		// specially create a complete item
		assertErrorMessage(t, DB.Delete("test_item_delete", newTestItem(i, rand.Int()%20)), nil)

		// check deletion result
		if item, found := DB.Query("test_item_delete").WhereInt("id", reindexer.EQ, mkID(i)).Get(); found {
			t.Fatalf("Item has not been deleted. < %+v > ", item)
		}
	}

	assertErrorMessage(t, DB.OpenNamespace("test_item_nested_pk", nsOpts, TestItemNestedPK{}), nil)
	for i := 1; i <= 30; i++ {
		assertErrorMessage(t, DB.Upsert("test_item_nested_pk", newTestItemNestedPK(i, rand.Int()%20)), nil)
	}

	assertErrorMessage(t, DB.MustBeginTx("test_item_nested_pk").Commit(), nil)
	assertErrorMessage(t, DB.CloseNamespace("test_item_nested_pk"), nil)
	assertErrorMessage(t, DB.OpenNamespace("test_item_nested_pk", nsOpts, TestItemNestedPK{}), nil)

	for i := 1; i <= 30; i++ {
		// specially create a complete item
		assertErrorMessage(t, DB.Delete("test_item_nested_pk", newTestItemNestedPK(i, rand.Int()%20)), nil)

		// check deletion result
		if item, found := DB.Query("test_item_nested_pk").WhereInt("id", reindexer.EQ, mkID(i)).Get(); found {
			t.Fatalf("Item has not been deleted. < %+v > ", item)
		}
	}
}
