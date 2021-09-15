package types

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for columns should fill any gaps in keys and Items with "", nil respectively
func TestColumnInsert(t *testing.T)  {
	col := Column{Name: "hi", Dtype: StringType, Items: []interface{}{"hi", "wow"}, keys: []string{"hio", "woo"}}
	col.insert(4, "yooo", "yeah")
	expectedItems := []interface{}{"hi", "wow", nil, nil, "yeah"}
	expectedKeys := []string{"hio", "woo", "", "", "yooo"}

	if !utils.AreSliceEqual(expectedItems, col.Items) {
		t.Fatalf("Items expected: %v, got %v", expectedItems, col.Items)
	}

	if !utils.AreStringSliceEqual(expectedKeys, col.keys) {
		t.Fatalf("Keys expected: %v, got %v", expectedKeys, col.keys)
	}
}