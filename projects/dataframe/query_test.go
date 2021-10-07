package dataframe

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/internal"
)

// AND combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any false on a given index,
// the final list has a false, else it has true
func TestAND(t *testing.T)  {
	testData := []internal.FilterType{
		{ true, true, true, true, true },
		{ true, false, true, true, true, false},
		{ true, false, false, true, },
		{ true, true, false, true, },	
	}

	expected := internal.FilterType{ true, false, false, true, false, false}
	got := AND(testData...)

	for i, expectedValue := range expected {
		if expectedValue != got[i] {
				t.Fatalf("on index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}	
}

// OR combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any true on a given index,
// the final list has a true in that index, else it has false,
// But then gives all consitituent arrays the OR version of all of them such that
// to allow for all 
func TestOR(t *testing.T)  {
	testData := []internal.FilterType{
		{ true, true, true, false, true },
		{ true, false, true, false, true, false},
		{ true, false, false, false, },
		{ true, true, false, false, },	
	}

	expected := internal.FilterType{ true, true, true, false, true, false}
	got := OR(testData...)

	for i, expectedValue := range expected {
		if expectedValue != got[i] {
				t.Fatalf("on index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}
}

// NOT combines a list of booleans into a list of booleans
// that index-wise, if there is any true on a given index,
// the final list has a false in that index, else it has true (it inverts the booleans)
func TestNOT(t *testing.T)  {
	testData := internal.FilterType{ true, true, false, true, true}

	expected := internal.FilterType{ !true, !true, !false, !true, !true}
	got := NOT(testData)

	for i, expectedValue := range expected {
		if expectedValue != got[i] {
			t.Fatalf("on index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}	
}
