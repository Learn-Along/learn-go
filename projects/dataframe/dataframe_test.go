package dataframe

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/internal"
	"github.com/learn-along/learn-go/projects/dataframe/internal/utils"
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

	for _, col := range df.cols {
		expectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name())
		var gotItems []interface{}

		switch items := col.Items().(type) {
		case []bool:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []int:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []string:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []float64:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}		
		}

		if !utils.AreSliceEqual(gotItems, expectedItems){
			t.Fatalf("col '%s' items expected: %v, got %v", col.Name(), expectedItems, col.Items())
		}
	}
}

func BenchmarkFromArray(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		FromArray(dataArray, primaryFields)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 10,265 ns/op	    	| 2248 B/op	      		 | 51 allocs/op  		 | x  	   |
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

	for _, col := range df.cols {
		expectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name())
		var gotItems []interface{}

		switch items := col.Items().(type) {
		case []bool:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			sort.Slice(gotItems, func(i, j int) bool {
				return gotItems[i].(bool) && gotItems[j].(bool)
			})
			sort.Slice(expectedItems, func(i, j int) bool {
				return expectedItems[i].(bool) && expectedItems[j].(bool)
			})

		case []int:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			sort.Slice(gotItems, func(i, j int) bool {
				return gotItems[i].(int) < gotItems[j].(int)
			})
			sort.Slice(expectedItems, func(i, j int) bool {
				return expectedItems[i].(int) < expectedItems[j].(int)
			})

		case []string:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			sort.Slice(gotItems, func(i, j int) bool {
				return gotItems[i].(string) < gotItems[j].(string)
			})
			sort.Slice(expectedItems, func(i, j int) bool {
				return expectedItems[i].(string) < expectedItems[j].(string)
			})

		case []float64:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}	
			
			sort.Slice(gotItems, func(i, j int) bool {
				return gotItems[i].(float64) < gotItems[j].(float64)
			})
			sort.Slice(expectedItems, func(i, j int) bool {
				return expectedItems[i].(float64) < expectedItems[j].(float64)
			})
		}

		if !utils.AreSliceEqual(gotItems, expectedItems){
			t.Fatalf("col '%s' items expected: %v, got %v", col.Name(), expectedItems, col.Items())
		}
	}
}

func BenchmarkFromMap(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		FromMap(dataMap, primaryFields)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None					    	| 8,609 ns/op	    	| 2244 B/op	      		 | 51 allocs/op		 	 | x  	   |
}

// Insert should Insert more records to the dataframe, overwriting any of the same key
func TestDataframe_Insert(t *testing.T)  {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]internal.Column{},
		index: map[interface{}]int{},
	}

	// Insert thrice, but still have the same data due to the primary keys...treat this like a db
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
		expectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name())
		var gotItems []interface{}

		switch items := col.Items().(type) {
		case []bool:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []int:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []string:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
		case []float64:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}		
		}

		if !utils.AreSliceEqual(gotItems, expectedItems){
			t.Fatalf("col '%s' items expected: %v, got %v", col.Name(), expectedItems, col.Items())
		}
	}
}

func BenchmarkDataframe_Insert(b *testing.B)  {
	df, err := FromArray(dataArray[:1], primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Insert(dataArray)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 7,575 ns/op	    	| 1792 B/op	      		 | 49 allocs/op			 | x 	   |
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
		cols: map[string]internal.Column{},
		index: map[interface{}]int{},
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
		initialExpectedItems := utils.ExtractFieldFromMapList(dataArray, col.Name())
		extraExpectedItems := utils.ExtractFieldFromMapList(extraData, col.Name())
		expectedItems := append(initialExpectedItems, extraExpectedItems...)
		var gotItems []interface{}

		switch items := col.Items().(type) {
		case []bool:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			// replace the nils with false
			for i, v := range expectedItems {
				if v == nil {
					expectedItems[i] = false
				}
			}
		case []int:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			// replace the nils with false
			for i, v := range expectedItems {
				if v == nil {
					expectedItems[i] = 0
				}
			}
			
		case []string:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}

			// replace the nils with false
			for i, v := range expectedItems {
				if v == nil {
					expectedItems[i] = ""
				}
			}

		case []float64:
			for _, v := range items {
				gotItems = append(gotItems, v)
			}
			
			// replace the nils with false
			for i, v := range expectedItems {
				if v == nil {
					expectedItems[i] = 0
				}
			}
		}

		if !utils.AreSliceEqual(gotItems, expectedItems){
			t.Fatalf("col '%s' items expected: %v, got %v", col.Name(), expectedItems, col.Items())
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

func BenchmarkDataframe_ToArray(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.ToArray()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 3815 ns/op	    	| 2456 B/op	      		 | 34 allocs/op			 | x	   |
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

func BenchmarkDataframe_ToArrayWithArgs(b *testing.B)  {
	fields := []string{"age", "location"}
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.ToArray(fields...)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 2884 ns/op	    	| 2264 B/op	      		 | 22 allocs/op			 |  	   |
}

// Delete should delete any records that fulfill a given condition
func TestDataframe_Delete(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter internal.FilterType;
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

// Insert, delete, Insert should update only those records that don't exist
func TestDataframe_DeleteReinsert(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter internal.FilterType;
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

func BenchmarkDataframe_Delete_GreaterThan(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkDelete(df, df.Col("age").GreaterThan(33), b)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 11431 ns/op	    	| 3232 B/op	      		 |76 allocs/op		 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Delete_IsLike(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkDelete(df, df.Col("last name").IsLike(regexp.MustCompile("oe$")), b)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 9229 ns/op	   		| 2176 B/op	      		 | 60 allocs/op		 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Delete_AND(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkDelete(df, AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)), b)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 11980 ns/op	    	| 3504 B/op	      		 | 79 allocs/op		 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Delete_OR(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkDelete(df, OR(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(45)), b)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 11395 ns/op	        | 3232 B/op	     		 | 76 allocs/op			 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Delete_NOT(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkDelete(df, NOT(df.Col("location").Equals("Kampala")), b)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 12073 ns/op	    	| 3312 B/op	      		 | 77 allocs/op			 | x 	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

// Update should update any records that fulfill a given condition,
// however, the primary keys should not be touched
// and any unknown columns are just added to all records, defaulting to nil values of the types of the rest
func TestDataframe_Update(t *testing.T)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		t.Fatalf("df error is: %s", err)
	}

	type testRecord struct {
		filter internal.FilterType;
		newData map[string]interface{};
		expected []map[string]interface{};
	}

	testTable := []testRecord{
		{
			filter: df.Col("age").LessOrEquals(33), 
			newData: map[string]interface{}{"location": "Kapchorwa", "new field": "yay", "age": 16},
			expected: []map[string]interface{}{
				{"first name": "John", "last name": "Doe", "age": 16, "location": "Kapchorwa", "new field": "yay" },
				{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka", "new field": "" },
				{"first name": "Paul", "last name": "Doe", "age": 16, "location": "Kapchorwa", "new field": "yay" },
				{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi", "new field": "" },
				{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi", "new field": "" },
				{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala", "new field": "" },
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

func BenchmarkDataframe_Update_GreaterThan(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkUpdate(
		df, 
		df.Col("age").LessOrEquals(33), 
		map[string]interface{}{"location": "Kapchorwa", "new field": "yay", "age": 16},
		b,
	)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 9,656 ns/op	    	| 2,496 B/op	      	 | 60 allocs/op	 	 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Update_IsLike(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkUpdate(
		df, 
		df.Col("last name").IsLike(regexp.MustCompile("oe$")),
		map[string]interface{}{"first name": "Hen", "age": 20,},
		b,
	)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 9408 ns/op	    	| 2240 B/op	      		 | 58 allocs/op	 	 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Update_AND(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkUpdate(
		df, 
		AND(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(33)), 
		map[string]interface{}{"age": 87},
		b,
	)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 8850 ns/op	    	| 2240 B/op      		 | 58 allocs/op		 	 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Update_OR(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkUpdate(
		df, 
		OR(df.Col("location").Equals("Kampala"), df.Col("age").GreaterThan(45)), 
		map[string]interface{}{"last name": "Rigobertha", "age": 73},
		b,
	)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 9,249 ns/op	    	| 2240 B/op	      		 | 58 allocs/op			 | x	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
}

func BenchmarkDataframe_Update_NOT(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	benchmarkUpdate(
		df,
		NOT(df.Col("location").Equals("Kampala")), 
		map[string]interface{}{"location": "Nebbi"},
		b,
	)

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 9,032 ns/op	    	| 2240 B/op	      		 | 58 allocs/op			 | x 	   |
	// | minus that for Insert 			| - 7,575 ns/op	    	| -1792 B/op	      	 | -49 allocs/op		 |		   |
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
                df.Col("location").Agg(func(arr internal.ItemSlice) internal.Item {return "random"}),
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
				df.Col("last name").Tx(func(v interface{}) interface{} {return fmt.Sprintf("last name: %v", v)}),
			),
			expected: []map[string]interface{}{
				{"last name": "last name: Roe", "age": 139,},
				{"last name": "last name: Doe", "age": 80,},
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

func BenchmarkDataframe_Select_Apply(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Select("age", "first name", "last name", "date").Apply(
			df.Col("age").Tx(func(v interface{}) interface{} {return v.(int) * 8}),
			df.Col("first name").Tx(func(v interface{}) interface{} { return fmt.Sprintf("name is %s", v) }),
		).Execute()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 25453 ns/op	   		| 10385 B/op	     	 | 194 allocs/op		 | x	   |
}

func BenchmarkDataframe_Select_Sortby(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Select("age", "first name", "last name", "location").SortBy(
			df.Col("last name").Order(ASC),
			df.Col("age").Order(DESC),                
		).Execute()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 36987 ns/op	    	| 14490 B/op	     	 | 259 allocs/op		 | x	   |
}

func BenchmarkDataframe_Select_Groupby(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Select("age", "last name", "first name").GroupBy("last name").Agg(
			df.Col("age").Agg(MEAN),
			df.Col("location").Agg(func(arr internal.ItemSlice) internal.Item {return "random"}),
		).Execute()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 26228 ns/op	    	| 9922 B/op	     		 | 186 allocs/op		 | x	   |
}

func BenchmarkDataframe_Select_Where(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Select().Where(
			AND(
				OR(
					df.Col("age").LessThan(20),
					df.Col("last name").IsLike(regexp.MustCompile("^(?i)roe$")),
				),
				df.Col("location").Equals("Kampala"),
			),
		).Execute()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 25300 ns/op	    	| 9922 B/op	     		 | 186 allocs/op		 | x	   |

}

func BenchmarkDataframe_Select_All_Combined(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Select("age", "last name").Where(
			df.Col("age").GreaterOrEquals(30),
		).GroupBy("last name").Agg(
			df.Col("age").Agg(SUM),
		).SortBy(
			df.Col("age").Order(DESC),
		).Apply(
			df.Col("age").Tx(func(v interface{}) interface{} {return fmt.Sprintf("total: %v", v)}),
		).Execute()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 46572 ns/op	   		| 20243 B/op	     	 | 354 allocs/op		 | x	   |

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

func BenchmarkDataframe_Clear(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Clear()
		df.Insert(dataArray)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 7372 ns/op	    	| 1632 B/op	      		 | 48 allocs/op		 	 | x	   |
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

func BenchmarkDataframe_Copy(b *testing.B)  {
	df, err := FromArray(dataArray, primaryFields)
	if err != nil {
		b.Fatalf("error creating df: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df.Copy()
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 11326 ns/op	    	| 4704 B/op	      		 | 85 allocs/op	 	 	 | x	   |
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

func BenchmarkDataframe_Merge(b *testing.B)  {
	df1, err := FromArray(dataArray[:1], primaryFields)
	if err != nil {
		b.Fatalf("df1 error is: %s", err)
	}

	df2, err := FromArray(dataArray[1:3], primaryFields)
	if err != nil {
		b.Fatalf("df2 error is: %s", err)
	}

	for i := 0; i < b.N; i++ {
		df1.Merge(df2)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 5636 ns/op	    	| 2272 B/op	      		 | 44 allocs/op			 | x	   |
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

func benchmarkDelete(df *Dataframe, filter internal.FilterType, b *testing.B)  {
	for i := 0; i < b.N; i++ {
		df.Delete(filter)
		df.Insert(dataArray)
	}
}

func benchmarkUpdate(df *Dataframe, filter internal.FilterType, data map[string]interface{}, b *testing.B)  {
	for i := 0; i < b.N; i++ {
		df.Update(filter, data)
		df.Insert(dataArray)
	}
}