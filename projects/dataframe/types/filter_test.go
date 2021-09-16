package types

import "testing"

// Coalesce combines the Filter (map of bool arrays) into a bool array
// using AND logic such that
// say, age > 7 and name = "John" returns only true when both are true
func TestCoalesce(t *testing.T)  {
	testData := Filter{
		"foo": { true, true, true, true, true},
		"bar": { true, false, true, true, true, false},
	}

	expected := []bool{true, false, true, true, true, false}
	got := testData.Coalesce()

	if len(expected) != len(got) {
		t.Fatalf("expected length: %d; got %d", len(expected), len(got))
	}

	for i, expectedValue := range expected {
		if expectedValue != got[i] {
			t.Fatalf("On index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}	
}
