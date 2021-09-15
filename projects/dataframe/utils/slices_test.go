package utils

import "testing"

// MergeSlices should combine two slices into one slice
func TestMergeSlices(t *testing.T)  {
	a := []interface{}{"foo", "bar", "hi"}
	b := []interface{}{"heyya", "woohoo"}

	ab := MergeSlices(a, b)
	expected := []interface{}{"foo", "bar", "hi", "heyya", "woohoo"}
	for i, item := range expected {
		if ab[i] != item {
			t.Fatalf("ab: index %d, expected: '%s', got '%s'", i, item, ab[i])
		}
	}

	ba := MergeSlices(b, a)
	expected = []interface{}{"heyya", "woohoo", "foo", "bar", "hi"}
	for i, item := range expected {
		if ba[i] != item {
			t.Fatalf("ba: index %d, expected: '%s', got '%s'", i, item, ba[i])
		}
	}
}

// SliceEquals should check if two slices are equal, and return an error if they are not
func TestSliceEquals(t *testing.T)  {
	a := []interface{}{"foo", "bar", "hi"}
	b := []interface{}{"heyya", "woohoo"}
	c := []interface{}{"heyya", "woohoo"}

	err := SliceEquals(c, b)
	if err != nil {
		t.Fatalf("c should be equal to b: error: %s", err)
	}

	err = SliceEquals(a, b)
	if err == nil {
		t.Fatal("a should not be equal to b")
	}
}