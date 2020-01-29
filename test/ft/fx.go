package ft

import (
	"flag"
	"fmt"
	"testing"

	"github.com/sapariduo/reindexer"
)

var dsn = flag.String("dsn", "builtin://", "reindex db dsn")
var dsnSlave = flag.String("dsnslave", "", "reindex slave db dsn")
var slaveCount = flag.Int("slavecount", 1, "reindex slave db count")
var benchmarkSeedCount = flag.Int("seedcount", 500000, "count of items for benchmark seed")
var benchmarkSeedCPU = flag.Int("seedcpu", 1, "number threads of for seeding")

func TestMain(m *testing.M) {
	flag.Parse()
}

type TextItem struct {
	ID        int `reindex:"id,,pk"`
	TextField string
}

func createReindexDbInstance(rx *reindexer.Reindexer, namespace string, indexType string) {
	err := rx.OpenNamespace(namespace, reindexer.DefaultNamespaceOptions(), TextItem{})
	if err != nil {
		panic(fmt.Errorf("Couldn't create namespace: "+namespace, err))
	}

	var config interface{}
	if indexType == "fuzzytext" {
		// Disable non exact searchers, disable stop word dictionat
		cfg := reindexer.DefaultFtFuzzyConfig()
		cfg.StopWords = []string{}
		cfg.Stemmers = []string{}
		cfg.EnableKbLayout = false
		cfg.EnableTranslit = false
		config = cfg
	} else {
		cfg := reindexer.DefaultFtFastConfig()
		cfg.StopWords = []string{}
		cfg.Stemmers = []string{}
		cfg.EnableKbLayout = false
		cfg.EnableTranslit = false
		config = cfg
	}

	rx.DropIndex(namespace, "text_field")
	err = rx.AddIndex(namespace, reindexer.IndexDef{
		Name:      "text_field",
		JSONPaths: []string{"TextField"},
		Config:    config,
		IndexType: indexType,
		FieldType: "string",
	})

	if err != nil {
		panic(fmt.Errorf("Couldn't set full text index config %s : %s", namespace, err.Error()))
	}
}

func fillReindexWithData(reindexDB *reindexer.Reindexer, namespace string, documents []string) {
	nextId := 1
	for _, document := range documents {
		item := TextItem{
			ID:        nextId,
			TextField: document,
		}
		if _, err := reindexDB.Insert(namespace, &item); err != nil {
			panic(err)
		}
		nextId++
	}
}
