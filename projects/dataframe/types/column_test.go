package types

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for columns should fill any gaps in keys and Items with "", nil respectively
func TestColumnInsert(t *testing.T)  {
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{0: "hi", 1: "wow"}}
	col.insert(4, "yeah")
	expectedItems := []interface{}{"hi", "wow", nil, nil, "yeah"}

	if !utils.AreSliceEqual(expectedItems, col.Items()) {
		t.Fatalf("items expected: %v, got %v", expectedItems, col.Items())
	}
}