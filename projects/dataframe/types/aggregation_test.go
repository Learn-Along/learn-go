package types

import "testing"

// MAX should return the maximum value (as a float64 value) in a given list of items
func TestMAX(t *testing.T)  {
	type testRecord struct {
		input []interface{};
		expected interface{}
	}

	testData := []testRecord{
		{
			input: []interface{}{"hi", "hello", "yoohoo", "salut"},
			expected: "yoohoo",
		},
		{
			input: []interface{}{1, 2.3, 4, 0.8},
			expected: 4.0,
		},
		{
			input: []interface{}{1, "hello", 5, "salut"},
			expected: nil,
		},
		{
			input: []interface{}{80, 78, 98, 4},
			expected: 98.0,
		},		{
			input: []interface{}{80.7, 78.6, 98.5, 98.509},
			expected: 98.509,
		},		{
			input: []interface{}{80, 78, 98, 4, nil, 7},
			expected: 98.0,
		},
		
	}

	for _, tr := range testData {
		got := MAX(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

// MIN should return the minimum value (as a float64 value) in a given list of items
func TestMIN(t *testing.T)  {
	type testRecord struct {
		input []interface{};
		expected interface{}
	}

	testData := []testRecord{
		{
			input: []interface{}{"hi", "hello", "yoohoo", "salut"},
			expected: "hello",
		},
		{
			input: []interface{}{1, 2.3, 4, 0.8},
			expected: 0.8,
		},
		{
			input: []interface{}{1, "hello", 5, "salut"},
			expected: nil,
		},
		{
			input: []interface{}{80, 78, 98, 4},
			expected: 4.0,
		},		{
			input: []interface{}{80.7, 78.6, 98.5, 98.509},
			expected: 78.6,
		},		{
			input: []interface{}{80, 78, 98, 4, nil, 7},
			expected: 4.0,
		},
		
	}

	for _, tr := range testData {
		got := MIN(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

// SUM should return the sum (as a float64 value) of the given list of items, but will return nil if 
// the values are not all numbers (nil values are treated as zero)
func TestSUM(t *testing.T)  {
	type testRecord struct {
		input []interface{};
		expected interface{}
	}

	testData := []testRecord{
		{
			input: []interface{}{"hi", "hello", "yoohoo", "salut"},
			expected: nil,
		},
		{
			input: []interface{}{1, 2.3, 4, 0.8},
			expected: 8.1,
		},
		{
			input: []interface{}{1, "hello", 5, "salut"},
			expected: nil,
		},
		{
			input: []interface{}{80, 78, 98, 4},
			expected: 260.0,
		},		{
			input: []interface{}{80.7, 78.6, 98.5, 98.509},
			expected: 356.309,
		},		{
			input: []interface{}{80, 78, 98, 4, nil, 7},
			expected: 267.0,
		},
		
	}

	for _, tr := range testData {
		got := SUM(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

// MEAN should return the mean (as a float64 value) of the given list of items.
// It returns nil if the values are not numbers (nil values are treated as zero)
func TestMEAN(t *testing.T)  {
	type testRecord struct {
		input []interface{};
		expected interface{}
	}

	testData := []testRecord{
		{
			input: []interface{}{"hi", "hello", "yoohoo", "salut"},
			expected: nil,
		},
		{
			input: []interface{}{1, 2.3, 4, 0.8},
			expected: 2.025,
		},
		{
			input: []interface{}{1, "hello", 5, "salut"},
			expected: nil,
		},
		{
			input: []interface{}{80, 78, 98, 4},
			expected: 65.0,
		},		{
			input: []interface{}{80.7, 78.6, 98.5, 98.509},
			expected: 89.07725,
		},		{
			input: []interface{}{80, 78, 98, 4, nil, 7},
			expected: 44.5,
		},
		
	}

	for _, tr := range testData {
		got := MEAN(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}