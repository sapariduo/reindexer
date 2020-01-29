package ft

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sapariduo/reindexer"
	_ "github.com/sapariduo/reindexer/bindings/builtin"
)

func doRankingTest(t *testing.T, indexType string) {

	totalSearchQuality := 0.0

	rx := reindexer.NewReindex(*dsn)
	defer rx.Close()

	testCases := ParseRankingTestCases()
	for i, testC := range testCases {
		// Need to make a local copy to use it in the closure below
		testCase := testC

		namespace := "ft_rank_" + strconv.Itoa(i)
		createReindexDbInstance(rx, namespace, indexType)
		fillReindexWithData(rx, namespace, testCase.AllDocuments)

		t.Run("Running test case: "+testCase.Description, func(t *testing.T) {
			for _, q := range testCase.Queries {
				// Need to make a local copy to use it in the closure below
				query := q
				t.Run("should match: "+query, func(t *testing.T) {
					dbItems, err := rx.Query(namespace).
						WhereString("text_field", reindexer.EQ, query, "").
						Exec().
						FetchAll()

					assert.NoError(t, err)
					expected := testCase.ExpectedDocuments
					actual := dbItemsToSliceOfDocuments(dbItems)
					assert.ElementsMatch(t, expected, actual)

					quality := calculateQuality(testCase.ExpectedDocuments, dbItemsToSliceOfDocuments(dbItems), testCase.AnyOrderClasses)
					quality = quality * 100
					totalSearchQuality += quality
					fmt.Printf("\n%s\nQuality - %.2f %%\n", testCase.Description, quality)
				})
			}
		})
	}

	totalSearchQuality = totalSearchQuality / float64(len(testCases))
	fmt.Printf("\nTotal Quality - %.2f %%\n", totalSearchQuality)
}

func dbItemsToSliceOfDocuments(dbItems []interface{}) []string {
	result := []string{}

	for _, dbItem := range dbItems {
		item := dbItem.(*TextItem)
		result = append(result, item.TextField)
	}

	return result
}

// This function calculates the quality of full text search based on
// Levenstein distance between itemsInExpectedOrder and itemsInActualOrder
// It also accepts a slice of equivalentClassesOfItems which represents classes of items that are equivalent to each other with respect to order
func calculateQuality(itemsInExpectedOrder []string, itemsInActualOrder []string, equivalentClassesOfItems [][]string) float64 {
	result := 0.0

	equivalentItemsMap := getEquivalentItemsMap(equivalentClassesOfItems)

	for i, expectedItem := range itemsInExpectedOrder {
		if expectedItem == itemsInActualOrder[i] {
			result += 1
		} else if elemIsInSlice(expectedItem, equivalentItemsMap[expectedItem]) {
			result += 1
		}
	}

	return result / float64(len(itemsInExpectedOrder))
}

func getEquivalentItemsMap(equivalentClassesOfItems [][]string) map[string][]string {
	result := map[string][]string{}

	for _, items := range equivalentClassesOfItems {
		for _, item := range items {
			result[item] = items
		}
	}

	return result
}

func elemIsInSlice(elem string, slice []string) bool {
	for _, sliceElem := range slice {
		if sliceElem == elem {
			return true
		}
	}
	return false
}

func TestFTFastRankingTest(t *testing.T) {
	doRankingTest(t, "text")
}
func TestFTFuzzyRankingTest(t *testing.T) {
	doRankingTest(t, "fuzzytext")
}
