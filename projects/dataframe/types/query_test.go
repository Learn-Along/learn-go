package types

import (
	"testing"
)

// AND combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any false on a given index,
// the final list has a false, else it has true
func TestAND(t *testing.T)  {
	testData := [][]bool{
		{ true, true, true, true, true},
		{ true, false, true, true, true, false},
		{ true, false, false, true, },
		{ true, true, false, true, },
	}

	expected := []bool{ true, false, false,	true, false, false}
	got := AND(testData...)

	for i, expectedValue := range expected {	
		if expectedValue != got[i] {
			t.Fatalf("on index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}
}

// OR combines a list of lists of booleans into a list of booleans
// that index-wise, if there is any true on a given index,
// the final list has a true in that index, else it has false
func TestOR(t *testing.T)  {
	testData := [][]bool{
		{ true, false, true, true, true, false},
		{ true, false, true, true, },
		{ true, false, false, false, },
		{ true, false, false, true, },
	}

	expected := []bool{ true, false, true,	true, true, false }
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
	testData := []bool{true, false, true, true, false, false, true }

	expected := []bool{!true, !false, !true, !true, !false, !false, !true }
	got := NOT(testData)

	for i, expectedValue := range expected {	
		if expectedValue != got[i] {
			t.Fatalf("on index %d expected: %v, got: %v", i, expectedValue, got[i])
		}
	}
}
