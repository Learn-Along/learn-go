package types

import (
	"regexp"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

var (
	dataArray = []map[string]interface{}{
		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	dataMap = map[interface{}]map[string]interface{}{
		"John Doe": {"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		"Jane Doe": {"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		"Paul Doe": {"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		"Richard Roe": {"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		"Reyna Roe": {"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		"Ruth Roe": {"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	primaryFields = []string{"first name", "last name"}
	expectedCols = utils.SortStringSlice([]string{"first name", "last name", "age", "location"}, utils.ASC)
	keys = []string{"John_Doe", "Jane_Doe", "Paul_Doe", "Richard_Roe", "Reyna_Roe", "Ruth_Roe"}
)

// fromArray should create a dataframe from an array of maps
func TestFromArray(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("error is: %s", err)
	}

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Fatalf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := utils.SortStringSlice(df.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}

	if !utils.AreStringSliceEqual(keys, df.Keys()) {
		t.Fatalf("keys expected: %v, got: %v", keys, df.Keys())
	}
}

// fromMap should create a dataframe from a map of maps
func TestFromMap(t *testing.T)  {
	df, err := FromMap(dataMap, primaryFields)
	if err != nil {
		t.Fatalf("error is: %s", err)
	}

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := utils.SortStringSlice(df.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}

	// since the map has disorganized order, we will sort them out first
	expectedKeys := utils.SortStringSlice(keys, utils.ASC)
	sortedKeys := utils.SortStringSlice(df.Keys(), utils.ASC)
	if !utils.AreStringSliceEqual(expectedKeys, sortedKeys) {
		t.Fatalf("keys expected: %v, got: %v", expectedKeys, sortedKeys)
	}
}

// Insert should insert more records to the dataframe, overwriting any of the same key
func TestInsert(t *testing.T)  {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]*Column{},
		index: map[string]int{},
		pks: OrderedMap{},
	}

	// insert thrice, but still have the same data due to the primary keys...treat this like a db
	df.Insert(dataArray)
	df.Insert(dataArray)
	df.Insert(dataArray)

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := utils.SortStringSlice(df.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}

	if !utils.AreStringSliceEqual(keys, df.Keys()) {
		t.Fatalf("keys expected: %v, got: %v", keys, df.Keys())
	}

	for _, col := range df.cols {
		expectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name)
		if !utils.AreSliceEqual(col.Items(), expectedItems){
			t.Fatalf("col '%s' items expected: %v, got %v", col.Name, expectedItems, col.Items())
		}
	}
}

// Insert should add the new records at the end of the dtaframe,
// while initializing the values for the non-existing columns to nil or its equivalent
// for the other prexisting values
func TestInsertNonExistingCols(t *testing.T)  {
	extraData := []map[string]interface{}{
		{"first name": "Roy", "last name": "Roe", "address": "Nairobi" },
		{"first name": "David", "last name": "Doe", "address": "Nairobi" },
	}
	allCols := utils.SortStringSlice(append(expectedCols, "address"), utils.ASC)
	allKeys := append(keys, "Roy_Roe", "David_Doe")

	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]*Column{},
		index: map[string]int{},
		pks: OrderedMap{},
	}

	// Insert the two sets of records
	df.Insert(dataArray)
	df.Insert(extraData)

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := utils.SortStringSlice(df.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(colNames, allCols){
		t.Fatalf("cols expected: %v, got: %v", allCols, colNames)
	}

	if !utils.AreStringSliceEqual(allKeys, df.Keys()) {
		t.Fatalf("keys expected: %v, got: %v", keys, df.Keys())
	}

	for _, col := range df.cols {
		initialExpectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name)
		extraExpectedItems := utils.ExtractFieldFromMapList(extraData, col.Name)
		expectedItems := append(initialExpectedItems, extraExpectedItems...)

		if !utils.AreSliceEqual(col.Items(), expectedItems){
			t.Errorf("col '%s' items expected: %v, got %v", col.Name, expectedItems, col.Items())
		}
	}	
}

// ToArray should convert the data into an array
func TestToArray(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	records, err := df.ToArray()
	if err != nil {
		t.Fatalf("error on ToArray is: %s", err)
	}

	if len(records) != len(dataArray) {
		t.Fatalf("expected number of records: %d, got %d", len(records), len(dataArray))
	}

	for i, record := range records {
		for field, value := range record {
			expected := dataArray[i][field]
			if expected != value {
				t.Fatalf("the record %d expected %v, got %v", i, expected, value)
			}
		}
	}
}

// Delete should delete any records that fulfill a given condition
func TestDelete(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter Filter;
		expected []map[string]interface{};
	}

	// FIXME: Reinserting records causes the tests to fail
	testTable := []testRecord{
		{
			filter: df.Col("age").GreaterThan(33), 
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
			},
		},
		{
			filter: df.Col("last name").IsLike(regexp.MustCompile("oe$")), 
			expected: []map[string]interface{}{},
		},
		{
			filter: df.Col("last name").IsLike(regexp.MustCompile("D")), 
			expected: []map[string]interface{}{
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
		// {
		// 	filter: AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)), 
		// 	expected: []map[string]interface{}{
		// 		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		// 		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		// 		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		// 		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		// 		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		// 	},
		// },
		// {
		// 	filter: OR(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(45)), 
		// 	expected: []map[string]interface{}{
		// 		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		// 		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		// 	},
		// },
		// {
		// 	filter: NOT(df.Col("location").Equals("Kampala")), 
		// 	expected: []map[string]interface{}{
		// 		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		// 		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		// 		{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
		// 	},
		// },
	}

	for loop, tr := range testTable {
		df.Insert(dataArray)
		if err != nil {
			t.Fatalf("df error is: %s", err)
		}

		err = df.Delete(tr.filter)
		if err != nil {
			t.Fatalf("df delete error is: %s", err)
		}

		records, err := df.ToArray()
		if err != nil {
			t.Fatalf("error on ToArray is: %s", err)
		}

		if len(records) != len(tr.expected) {
			t.Fatalf("loop %d, expected number of records: %d, got %d", loop, len(tr.expected), len(records))
		}

		for i, record := range records {
			for field, value := range record {
				expectedValue := tr.expected[i][field]
				if expectedValue != value {
					t.Fatalf("loop %d, the record %d expected %v, got %v", loop, i, expectedValue, value)
				}
			}
		}		
	}
}