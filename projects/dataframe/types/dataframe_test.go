package types

import (
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
	expectedCols = []string{"first name", "last name", "age", "location"}
	keys = []string{"John_Doe", "Jane_Doe", "Paul_Doe", "Richard_Roe", "Reyna_Roe", "Ruth_Roe"}
)

// fromArray should create a dataframe from an array of maps
func TestFromArray(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := df.getColNames()
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}
}

// fromMap should create a dataframe from a map of maps
func TestFromMap(t *testing.T)  {
	df, err := FromMap(dataMap, primaryFields)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := df.getColNames()
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}
}

// Insert should insert more records to the dataframe, overwriting any of the same key
func TestInsert(t *testing.T)  {
	df := Dataframe{pkFields: primaryFields}

	// insert thrice, but still have the same data
	df.Insert(dataArray)
	df.Insert(dataArray)
	df.Insert(dataArray)

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := df.getColNames()
	if !utils.AreStringSliceEqual(colNames, expectedCols){
		t.Fatalf("cols expected: %v, got: %v", expectedCols, colNames)
	}

	for _, col := range df.cols {
		if !utils.AreStringSliceEqual(col.keys, keys){
			t.Errorf("col '%s' keys expected: %v, got %v", col.Name, keys, col.keys)
		}

		expectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name)
		if !utils.AreSliceEqual(col.Items, expectedItems){
			t.Errorf("col '%s' items expected: %v, got %v", col.Name, expectedItems, col.Items)
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
	allCols := append(expectedCols, "address")
	allKeys := append(keys, "Roy_Roe", "David_Doe")

	df := Dataframe{pkFields: primaryFields}

	// Insert the two sets of records
	df.Insert(dataArray)
	df.Insert(extraData)

	if !utils.AreStringSliceEqual(df.pkFields, primaryFields){
		t.Errorf("pkFields expected: %v, got %v", primaryFields, df.pkFields)
	}

	colNames := df.getColNames()
	if !utils.AreStringSliceEqual(colNames, allCols){
		t.Fatalf("cols expected: %v, got: %v", allCols, colNames)
	}

	for _, col := range df.cols {
		if !utils.AreStringSliceEqual(col.keys, allKeys){
			t.Errorf("col '%s' keys expected: %v, got %v", col.Name, allKeys, col.keys)
		}

		initialExpectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name)
		extraExpectedItems := utils.ExtractFieldFromMapList(extraData, col.Name)
		expectedItems := append(initialExpectedItems, extraExpectedItems...)

		if !utils.AreSliceEqual(col.Items, expectedItems){
			t.Errorf("col '%s' items expected: %v, got %v", col.Name, expectedItems, col.Items)
		}
	}	
}