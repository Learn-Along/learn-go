package utils

import "testing"

// AreStringSliceEqual should check if two slices of string arrays are equal, and return an false if they are not
func TestAreStringSliceEqual(t *testing.T)  {
	a := []string{"foo", "bar", "hi"}
	b := []string{"heyya", "woohoo"}
	c := []string{"heyya", "woohoo"}

	if !AreStringSliceEqual(c, b) {
		t.Fatal("c should be equal to b")
	}

	if AreStringSliceEqual(a, b) {
		t.Fatal("a should not be equal to b")
	}
}

// AreSliceEqual should check if two slices are equal, and return an false if they are not
func TestAreSliceEqual(t *testing.T)  {
	a := []interface{}{"foo", "bar", "hi"}
	b := []interface{}{"heyya", "woohoo"}
	c := []interface{}{"heyya", "woohoo"}

	if !AreSliceEqual(c, b) {
		t.Fatal("c should be equal to b")
	}

	if AreSliceEqual(a, b) {
		t.Fatal("a should not be equal to b")
	}
}

// ExtractFieldFromMapList should extract a given field from a list of maps and return
// a list fo the values for that given field
func TestExtractFieldFromMapList(t *testing.T)  {
	data := []map[string]interface{}{
		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	testTable := map[string][]interface{}{
		"first name": {"John", "Jane", "Paul", "Richard", "Reyna", "Ruth",},
		"last name": {"Doe", "Doe", "Doe", "Roe", "Roe", "Roe"},
		"age": {30, 50, 19, 34, 45, 60},
		"location": {"Kampala", "Lusaka", "Kampala", "Nairobi", "Nairobi", "Kampala"},
	}

	for field, expected := range testTable {
		got := ExtractFieldFromMapList(data, field)

		if !AreSliceEqual(expected, got) {
			t.Fatalf("for field: '%s', expected: %v, got %v", field, expected, got)
		}
	}
}

// SortStringSlice sorts a slice of strings in either ascending or descending order, immutably
func TestSortStringSlice(t *testing.T)  {
	a := []string{"foo", "bar", "hi", "heyya", "woohoo"}

	expected := []string{"bar", "foo", "heyya", "hi", "woohoo"}
	got := SortStringSlice(a, ASC)
	if !AreStringSliceEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}

	expected = []string{"woohoo", "hi", "heyya", "foo", "bar",}
	got = SortStringSlice(a, DESC)
	if !AreStringSliceEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// ConvertToStringSlice converts an []interface{} to []string
func TestConvertToStringSlice(t *testing.T)  {
	type testRecord struct {
		slice []interface{};
		ignoreNil bool;
		expected []string;
	}

	testData := []testRecord{
		{
			slice: []interface{}{"hi", "hello", "yoohoo", "salut"},
			ignoreNil: false,
			expected: []string{"hi", "hello", "yoohoo", "salut"},
		},
		{
			slice: []interface{}{"hi", "hello", "salut"},
			ignoreNil: true,
			expected: []string{"hi", "hello", "salut"},
		},
		{
			slice: []interface{}{"hi", "hello", nil, "salut"},
			ignoreNil: true,
			expected: []string{"hi", "hello", "salut"},
		},
		{
			slice: []interface{}{"hi", "hello", nil, "salut"},
			ignoreNil: false,
			expected: []string{"hi", "hello", "", "salut"},
		},
	}

	for _, tr := range testData {
		got := ConvertToStringSlice(tr.slice, tr.ignoreNil)
		if !AreStringSliceEqual(got, tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}