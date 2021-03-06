package types

import (
	"fmt"
	"log"
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
func TestDataframe_Insert(t *testing.T)  {
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
func TestDataframe_InsertNonExistingCols(t *testing.T)  {
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

// ToArray should convert the data into an array. If no string args are passed,
// the values have all the fields
func TestDataframe_ToArray(t *testing.T)  {
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

// ToArray should convert the data into an array. If string args are passed,
// the values have the specified fields only
func TestDataframe_ToArrayWithArgs(t *testing.T)  {
	fields := []string{"age", "location"}
	excludedFields := []string{"last name", "first name"}

	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	records, err := df.ToArray(fields...)
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

		for _, excludedField := range excludedFields {
			if _, exists := record[excludedField]; exists {
				t.Fatalf("excluded field %v has been included in \n %v", excludedField, record)
			}
		}
	}
}

// Delete should delete any records that fulfill a given condition
func TestDataframe_Delete(t *testing.T)  {
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
func TestDataframe_DeleteReinsert(t *testing.T)  {
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
func TestDataframe_Update(t *testing.T)  {
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

// Select should be able to query data allowing for selection of fields,
// sorting, grouping, filtering, applying etc.
func TestDataframe_Select(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		q *query;
		expected []map[string]interface{};
	}

	testTable := []testRecord{
		{
			// select will ignore columns like 'date' that don't exist in the dataframe
			q: df.Select("age", "first name", "last name", "date").Apply(
				df.Col("age").Tx(func(v interface{}) interface{} {return v.(int) * 8}),
				df.Col("first name").Tx(func(v interface{}) interface{} { return fmt.Sprintf("name is %s", v) }),
			), 
			expected: []map[string]interface{}{
				{"first name": "name is John", "last name": "Doe", "age":  8*30, },
				{"first name": "name is Jane", "last name": "Doe", "age": 8*50, },
				{"first name": "name is Paul", "last name": "Doe", "age": 8*19, },
				{"first name": "name is Richard", "last name": "Roe", "age": 8*34, },
				{"first name": "name is Reyna", "last name": "Roe", "age": 8*45, },
				{"first name": "name is Ruth", "last name": "Roe", "age": 8*60, },
			},
		},
		{
			q: df.Select("age", "first name", "last name", "location").SortBy(
				df.Col("last name").Order(ASC),
                df.Col("age").Order(DESC),                
            ), 
			expected: []map[string]interface{}{
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
				{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
			},
		},
		{
			// all columns that are not part of the GroupBy will be ignored in the select as they make no sense
			// select will also ignore any columns in the groupby that were not passed in the list of selects
			q: df.Select("age", "last name", "first name").GroupBy("last name").Agg(
                df.Col("age").Agg(MEAN),
				// even a custom agggregate functions are possible
                df.Col("location").Agg(func(arr []interface{}) interface{}{return "random"}),
            ), 
			expected: []map[string]interface{}{
				{"last name": "Doe", "age": float64(33) },
				{"last name": "Roe", "age": float64(139) / 3},
			},
		},
		{
			// Passing no fields in Select returns all columns
			q: df.Select().Where(
				AND(
					OR(
						df.Col("age").LessThan(20),
						df.Col("last name").IsLike(regexp.MustCompile("^(?i)roe$")),
					),
					df.Col("location").Equals("Kampala"),
				),
			),
			expected: []map[string]interface{}{
				{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
			},
		},
		{
			q: df.Select("age", "last name").Where(
				df.Col("age").GreaterOrEquals(30),
			).GroupBy("last name").Agg(
				df.Col("age").Agg(SUM),
			).SortBy(
				df.Col("age").Order(DESC),
			).Apply(
				df.Col("age").Tx(func(v interface{}) interface{} {return fmt.Sprintf("total: %v", v)}),
			),
			expected: []map[string]interface{}{
				{"last name": "Roe", "age": "total: 139",},
				{"last name": "Doe", "age": "total: 80",},
			},
		},
	}

	for loop, tr := range testTable {
		df.Clear()

		df.Insert(dataArray)
		if err != nil {
			t.Fatalf("df error is: %s", err)
		}

		records, err := tr.q.Execute()
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
func TestDataframe_Clear(t *testing.T)  {
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

// Copy should make a new Dataframe that resembles the dataframe but
// has no reference to the items of the previous Dataframe
func TestDataframe_Copy(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	newDf, err := df.Copy()
	if err != nil {
		t.Fatalf("df copy error is: %s", err)
	}

	if newDf == df {
		t.Fatalf("expected %p not to equal %p", newDf, df)
	}

	if !utils.AreStringSliceEqual(df.pkFields, newDf.pkFields){
		t.Fatalf("new df pkFields expected: %v, got %v", df.pkFields, newDf.pkFields)
	}

	oldCols := utils.SortStringSlice(df.ColumnNames(), utils.ASC)
	newCols := utils.SortStringSlice(newDf.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(oldCols, newCols){
		t.Fatalf("new df column names expected: %v, got %v", oldCols, newCols)
	}

	if !utils.AreStringSliceEqual(df.Keys(), newDf.Keys()){
		t.Fatalf("new df keys expected: %v, got %v", df.Keys(), newDf.Keys())
	}

	for key, col := range df.cols {
		newDfCol := newDf.cols[key]
		if newDfCol == col {
			t.Fatalf("expected col '%s' of address %p not to equal %p", key, newDfCol, col)
		}
	}

	newDfRecords, err := newDf.ToArray()
	if err != nil {
		t.Fatalf("newDf ToArray error is: %s", err)
	}

	oldRecords, err := df.ToArray()
	if err != nil {
		t.Fatalf("df ToArray error is: %s", err)
	}

	for i, record := range dataArray {
		for field, expected := range record {
			newDfValue := newDfRecords[i][field]
			oldDfValue := oldRecords[i][field]

			if expected != oldDfValue {
				t.Fatalf("Old Df: the record %d expected %v, got %v", i, expected, oldDfValue)
			}

			if expected != newDfValue {
				t.Fatalf("New Df: the record %d expected %v, got %v", i, expected, newDfValue)
			}
		}
	}
}

// Merge combines into the given dataframe, the dataframes passed, overwriting any records that
// have the same primary key value
func TestDataframe_Merge(t *testing.T)  {
	df1, err := FromArray(dataArray[:1], primaryFields)
	if err != nil {
		t.Fatalf("df1 error is: %s", err)
	}

	df2, err := FromArray(dataArray[1:3], primaryFields)
	if err != nil {
		t.Fatalf("df2 error is: %s", err)
	}

	df3, err := FromArray(dataArray[3:], primaryFields)
	if err != nil {
		t.Fatalf("df3 error is: %s", err)
	}

	df4, err := FromArray(dataArray[:1], primaryFields)
	if err != nil {
		t.Fatalf("df4 error is: %s", err)
	}

	err = df1.Merge(df2, df3, df4)
	if err != nil {
		t.Fatalf("Merge error %s", err)
	}

	if !utils.AreStringSliceEqual(df1.pkFields, primaryFields){
		t.Fatalf("pkFields expected: %v, got %v", primaryFields, df1.pkFields)
	}

	colNames := utils.SortStringSlice(df1.ColumnNames(), utils.ASC)
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}

	if !utils.AreStringSliceEqual(keys, df1.Keys()) {
		t.Fatalf("keys expected: %v, got: %v", keys, df1.Keys())
	}
}

// The PrettyPrintRecords method prints out the records in a pretty format
func ExampleDataframe_PrettyPrintRecords()  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		log.Fatalf("df error is: %s", err)
	}
	
	df.PrettyPrintRecords()
	// Output:
	// [
	// 	{
	// 		"age": 30,
	// 		"first name": "John",
	// 		"last name": "Doe",
	// 		"location": "Kampala"
	// 	},
	// 	{
	// 		"age": 50,
	// 		"first name": "Jane",
	// 		"last name": "Doe",
	// 		"location": "Lusaka"
	// 	},
	// 	{
	// 		"age": 19,
	// 		"first name": "Paul",
	// 		"last name": "Doe",
	// 		"location": "Kampala"
	// 	},
	// 	{
	// 		"age": 34,
	// 		"first name": "Richard",
	// 		"last name": "Roe",
	// 		"location": "Nairobi"
	// 	},
	// 	{
	// 		"age": 45,
	// 		"first name": "Reyna",
	// 		"last name": "Roe",
	// 		"location": "Nairobi"
	// 	},
	// 	{
	// 		"age": 60,
	// 		"first name": "Ruth",
	// 		"last name": "Roe",
	// 		"location": "Kampala"
	// 	}
	// ]
}
