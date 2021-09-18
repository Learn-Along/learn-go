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
	noOfExpectedCols = len(expectedCols)
	keys = []string{"John_Doe", "Jane_Doe", "Paul_Doe", "Richard_Roe", "Reyna_Roe", "Ruth_Roe"}
	noOfExpectedKeys = len(keys)
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
		index: map[interface{}]int{},
		// pks: orderedMapType{},
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
		index: map[interface{}]int{},
		// pks: orderedMapType{},
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
		filter filterType;
		expected []map[string]interface{};
	}

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
		{
			filter: AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)), 
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
			},
		},
		{
			filter: OR(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(45)), 
			expected: []map[string]interface{}{
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
			},
		},
		{
			filter: NOT(df.Col("location").Equals("Kampala")), 
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
	}

	for loop, tr := range testTable {
		df.Clear()

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
					t.Fatalf("loop %d, the record %d expected %v, got %v, \n records: %v", loop, i, expectedValue, value, records)
				}
			}
		}		
	}
}

// Insert, delete, insert should update only those records that don't exist
func TestDeleteReinsert(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter filterType;
		onReinsert []map[string]interface{};
	}

	testTable := []testRecord{
		{
			filter: df.Col("age").GreaterThan(33), 
			onReinsert: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
		{
			filter: df.Col("last name").IsLike(regexp.MustCompile("D")), 
			onReinsert: []map[string]interface{}{
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
			},
		},
		{
			filter: AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)), 
			onReinsert: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
	}

	for loop, tr := range testTable {
		df.Clear()

		df.Insert(dataArray)
		if err != nil {
			t.Fatalf("df error is: %s", err)
		}

		err = df.Delete(tr.filter)
		if err != nil {
			t.Fatalf("df delete error is: %s", err)
		}

		// reinsert 
		df.Insert(dataArray)
		if err != nil {
			t.Fatalf("df error is: %s", err)
		}

		records, err := df.ToArray()
		if err != nil {
			t.Fatalf("error on ToArray is: %s", err)
		}

		if len(records) != len(tr.onReinsert) {
			t.Fatalf("loop %d, expected number of records: %d, got %d", loop, len(tr.onReinsert), len(records))
		}

		for i, record := range records {
			for field, value := range record {
				expectedValue := tr.onReinsert[i][field]
				if expectedValue != value {
					t.Fatalf("loop %d, the record %d expected %v, got %v, \n records: %v", loop, i, expectedValue, value, records)
				}
			}
		}		
	}
}

// Update should update any records that fulfill a given condition,
// however, the primary keys should not be touched
// and any unknown columns are just added to all records, defaulting to nil for the rest
func TestUpdate(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter filterType;
		newData map[string]interface{};
		expected []map[string]interface{};
	}

	testTable := []testRecord{
		{
			filter: df.Col("age").LessOrEquals(33), 
			newData: map[string]interface{}{"location": "Kapchorwa", "new field": "yay", "age": 16},
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 16, "location": "Kapchorwa", "new field": "yay" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka", "new field": nil },
				{"first name": "Paul", "last name": "Doe", "age": 16, "location": "Kapchorwa", "new field": "yay" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi", "new field": nil },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi", "new field": nil },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala", "new field": nil },
			},
		},
		{
			filter: df.Col("last name").IsLike(regexp.MustCompile("oe$")), 
			newData: map[string]interface{}{"first name": "Hen", "age": 20,},
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 20, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 20, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 20, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 20, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 20, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 20, "location": "Kampala" },
			},
		},
		{
			filter: df.Col("last name").IsLike(regexp.MustCompile("D")), 
			newData: map[string]interface{}{"location": "Bujumbura"},
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Bujumbura" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Bujumbura" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Bujumbura" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
		{
			filter: AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)),
			newData: map[string]interface{}{"age": 87}, 
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 87, "location": "Kampala" },
			},
		},
		{
			filter: OR(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(45)),
			newData: map[string]interface{}{"last name": "Rigobertha", "age": 73}, 
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 73, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 73, "location": "Lusaka" },
				{"first name": "Paul", "last name": "Doe", "age": 73, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Ruth", "last name": "Roe", "age": 73, "location": "Kampala" },
			},
		},
		{
			filter: NOT(df.Col("location").Equals("Kampala")), 
			newData: map[string]interface{}{"location": "Nebbi"},
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Nebbi" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nebbi" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nebbi" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
	}

	for loop, tr := range testTable {
		df.Clear()

		df.Insert(dataArray)
		if err != nil {
			t.Fatalf("df error is: %s", err)
		}

		err = df.Update(tr.filter, tr.newData)
		if err != nil {
			t.Fatalf("df update error is: %s", err)
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
					t.Fatalf("loop %d, the record %d expected %v, got %v, \n records: %v", loop, i, expectedValue, value, records)
				}
			}
		}		
	}
}

// Clear should clear all the cols, index and pks
func TestClear(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Fatalf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	noOfColumns := len(df.ColumnNames())
	if noOfColumns != noOfExpectedCols {
		t.Fatalf("number of cols expected: %v, got: %v", noOfExpectedCols, noOfColumns)
	}

	noOfKeys := len(df.Keys())
	if noOfKeys != noOfExpectedKeys  {
		t.Fatalf("numer of keys expected: %v, got: %v", noOfExpectedKeys, noOfKeys)
	}

	indices := len(df.index)
	if indices != noOfKeys {
		t.Fatalf("number of indices expected: %v; got: %v", noOfKeys, indices)
	}

	df.Clear()

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Fatalf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	noOfColumns = len(df.ColumnNames())
	if noOfColumns != 0 {
		t.Fatalf("number of cols expected: %v, got: %v", 0, noOfColumns)
	}

	noOfKeys = len(df.Keys())
	if noOfKeys != 0  {
		t.Fatalf("number of keys expected: %v, got: %v", 0, noOfKeys)
	}

	indices = len(df.index)
	if indices != 0 {
		t.Fatalf("number of indices expected: %v; got: %v", 0, indices)
	}
}