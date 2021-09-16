package types

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// ToSlice converts an ordered map into a slice
func TestToSlice(t *testing.T)  {
	type testRecord struct {
		input OrderedMap;
		expected []interface{}
	}

	testData := []testRecord{
		{
			input: OrderedMap{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			expected: []interface{}{"hi", "hello", "yoohoo", "salut"},
		},
		{
			input: OrderedMap{0: 4, 1: 2, 2: "yoohoo", 3: "salut"},
			expected: []interface{}{4, 2, "yoohoo", "salut"},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice()
		if !utils.AreSliceEqual(got, tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}