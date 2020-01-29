package basic

import (
	"flag"
	"fmt"
	"strconv"

	"sort"

	"github.com/sapariduo/reindexer"
	_ "github.com/sapariduo/reindexer/bindings/builtin"
	_ "github.com/sapariduo/reindexer/bindings/cproto"
	. "github.com/sapariduo/reindexer/test/ft/specs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dsn = flag.String("dsn", "builtin://", "reindex db dsn")

var _ = Describe("", func() {
	flag.Parse()

	for i, testC := range ParseBasicTestCases() {
		// Need to make a local copy to use it in the closure below
		testCase := testC

		namespace := "ft_" + strconv.Itoa(i)
		reindexDB := createReindexDbInstance(namespace)
		fillReindexWithData(reindexDB, namespace, testCase.AllDocuments)
		Context("Running test case: "+testCase.Description, func() {
			for _, validQ := range testCase.ValidQueries {
				// Need to make a local copy to use it in the closure below
				validQuery := validQ
				It("should match: "+validQuery, func() {
					dbItems, err := reindexDB.Query(namespace).
						WhereString("text_field", reindexer.EQ, validQuery, "").
						Exec().
						FetchAll()

					Expect(err).To(BeNil())
					expected := testCase.ExpectedDocuments
					actual := dbItemsToSliceOfDocuments(dbItems)
					sort.Strings(expected)
					sort.Strings(actual)
					var exists = struct{}{}

					actualMap := make(map[string]struct{})
					for _, s := range actual {
						actualMap[s] = exists
					}
					for _, s := range expected {
						if _, ok := actualMap[s]; !ok {
							panic(GINKGO_PANIC)
						}
					}

				})
			}
			for _, invalidQ := range testCase.InvalidQueries {
				// Need to make a local copy to use it in the closure below
				invalidQuery := invalidQ
				It("shouldn't match: "+invalidQuery, func() {
					dbItems, err := reindexDB.Query(namespace).
						WhereString("text_field", reindexer.EQ, invalidQuery).
						Exec().
						FetchAll()

					Expect(err).To(BeNil())
					for _, document := range dbItemsToSliceOfDocuments(dbItems) {
						Expect(testCase.ExpectedDocuments).NotTo(ContainElement(document))
					}
				})
			}
		})
	}
})

type TextItem struct {
	ID        int    `reindex:"id,,pk"`
	TextField string `reindex:"text_field,fuzzytext"`
}

func createReindexDbInstance(namespace string) *reindexer.Reindexer {
	reindexDB := reindexer.NewReindex(*dsn)
	err := reindexDB.OpenNamespace(namespace, reindexer.DefaultNamespaceOptions(), TextItem{})
	if err != nil {
		panic(fmt.Errorf("Couldn't create namespace: "+namespace, err))
	}

	// Disable non exact searchers, disable stop word dictionat
	ftConfig := reindexer.DefaultFtFuzzyConfig()
	ftConfig.StopWords = []string{}
	ftConfig.Stemmers = []string{}
	ftConfig.EnableKbLayout = false
	ftConfig.EnableTranslit = false
	err = reindexDB.ConfigureIndex(namespace, "text_field", ftConfig)
	if err != nil {
		panic(fmt.Errorf("Couldn't set full text index config: "+namespace, err))

	}

	return reindexDB
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

func dbItemsToSliceOfDocuments(dbItems []interface{}) []string {
	result := []string{}

	for _, dbItem := range dbItems {
		item := dbItem.(*TextItem)
		result = append(result, item.TextField)
	}

	return result
}
