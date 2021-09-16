package types

import (
	"testing"
)

// AND combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any false on a given index,
// the final list has a false, else it has true
func TestAND(t *testing.T)  {
	testData := []Filter{
		{
			"foo": { true, true, true, true, true},
			"bar": { true, false, true, true, true, false},
		},
		{
			"foo": { true, false, false, true, },
			"bar": { true, true, false, true, },
		},		
	}

	expected := Filter{
		"foo": { true, false, false, true, false},
		"bar": { true, false, false, true, false, false},
	}
	got := AND(testData...)

	for key, expectedArray := range expected {
		for i, expectedValue := range expectedArray {
			if expectedValue != got[key][i] {
				t.Fatalf("for key %s, on index %d expected: %v, got: %v", key, i, expectedValue, got[key][i])
			}
		}	
	}
}

// OR combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any true on a given index,
// the final list has a true in that index, else it has false
func TestOR(t *testing.T)  {
	testData := []Filter{
		{
			"foo": { true, true, false, true, true},
			"bar": { true, false, true, true, true, false},
		},
		{
			"foo": { true, false, false, true, },
			"bar": { true, false, false, true, },
		},		
	}

	expected := Filter{
		"foo": { true, true, false, true, true},
		"bar": { true, false, true, true, true, false},
	}
	got := OR(testData...)

	for key, expectedArray := range expected {
		for i, expectedValue := range expectedArray {
			if expectedValue != got[key][i] {
				t.Fatalf("for key %s, on index %d expected: %v, got: %v", key, i, expectedValue, got[key][i])
			}
		}	
	}
}

// NOT combines a list of booleans into a list of booleans
// that index-wise, if there is any true on a given index,
// the final list has a false in that index, else it has true (it inverts the booleans)
func TestNOT(t *testing.T)  {
	testData := Filter{
		"foo": { true, true, false, true, true},
		"bar": { true, false, true, true, true, false},
	}

	expected := Filter{
		"foo": { !true, !true, !false, !true, !true},
		"bar": { !true, !false, !true, !true, !true, !false},
	}
	got := NOT(testData)

	for key, expectedArray := range expected {
		for i, expectedValue := range expectedArray {
			if expectedValue != got[key][i] {
				t.Fatalf("for key %s, on index %d expected: %v, got: %v", key, i, expectedValue, got[key][i])
			}
		}	
	}
}
