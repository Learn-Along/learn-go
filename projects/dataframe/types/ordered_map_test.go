package types

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// ToSlice converts an ordered map into a slice, ignoring gaps in indices automatically
func TestToSlice(t *testing.T)  {
	type testRecord struct {
		input orderedMapType;
		expected []interface{}
	}

	testData := []testRecord{
		{
			input: orderedMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			expected: []interface{}{"hi", "hello", "yoohoo", "salut"},
		},
		{
			input: orderedMapType{0: 4, 1: 2, 2: "yoohoo", 3: "salut"},
			expected: []interface{}{4, 2, "yoohoo", "salut"},
		},
		{
			input: orderedMapType{0: 4, 2: "yoohoo", 3: "salut"},
			expected: []interface{}{4, "yoohoo", "salut"},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice()
		if !utils.AreSliceEqual(got, tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

// Defragmentize should reorder the ordered dict basing
// on the new indices passed to it
func TestDefragmentize(t *testing.T)  {
	type testRecord struct {
		_map orderedMapType;
		newOrder []int;
		expected []interface{}
	}

	testData := []testRecord{
		{
			_map: orderedMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			newOrder: []int{3, 1, 0},
			expected: []interface{}{"salut", "hello", "hi"},
		},
		{
			_map: orderedMapType{0: 4, 1: 2, 2: "yoohoo", 3: "salut"},
			newOrder: []int{1, 0, 3, 2},
			expected: []interface{}{2, 4, "salut", "yoohoo"},
		},
	}

	for _, tr := range testData {
		tr._map.Defragmentize(tr.newOrder)
		got := tr._map.ToSlice()

		if !utils.AreSliceEqual(got, tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}