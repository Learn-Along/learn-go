package types

import (
	"testing"
)

// fromArray should create a dataframe from an array of maps
func TestFromArray(t *testing.T)  {
	data := []map[string]interface{}{
		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	primaryFields := []string{"first name", "last name"}
	expectedCols := []string{"first name", "last name", "age", "location"}

	df, err := FromArray(data, primaryFields)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	for i := 0; i < len(primaryFields); i++ {
		if df.pkFields[i] != primaryFields[i] {
			t.Errorf("pkField %d expected value: '%s', got '%s'", i, primaryFields[i], df.pkFields[i])
		}
	}

	for _, col := range expectedCols {
		if dfCol, ok := df.cols[col]; !ok {
			t.Fatalf("Col '%s' does not exist on the dataframe", col)
		} else if dfCol.Name != col {
			t.Fatalf("Col.Name '%s' does not match the column name: '%s'", dfCol.Name, col)
		}
	}
}

// fromMap should create a dataframe from a map of maps
func TestFromMap(t *testing.T)  {
	data := map[interface{}]map[string]interface{}{
		"John Doe": {"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		"Jane Doe": {"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		"Paul Doe": {"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		"Richard Roe": {"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		"Reyna Roe": {"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		"Ruth Roe": {"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	primaryFields := []string{"first name", "last name"}
	expectedCols := []string{"first name", "last name", "age", "location"}

	df, err := FromMap(data, primaryFields)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	for i := 0; i < len(primaryFields); i++ {
		if df.pkFields[i] != primaryFields[i] {
			t.Errorf("pkField %d expected value: '%s', got '%s'", i, primaryFields[i], df.pkFields[i])
		}
	}

	for _, col := range expectedCols {
		if dfCol, ok := df.cols[col]; !ok {
			t.Fatalf("Col '%s' does not exist on the dataframe", col)
		} else if dfCol.Name != col {
			t.Fatalf("Col.Name '%s' does not match the column name: '%s'", dfCol.Name, col)
		}
	}
}

// Insert should insert more records to the dataframe, overwriting any of the same key
func TestInsert(t *testing.T)  {
	data := []map[string]interface{}{
		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
		{"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
		{"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
		{"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
	}
	primaryFields := []string{"first name", "last name"}
	expectedCols := []string{"first name", "last name", "age", "location"}	
	keys := []string{"John_Doe", "Jane_Doe", "Paul_Doe", "Richard_Roe", "Reyna_Roe", "Ruth_Roe"}

	df := Dataframe{pkFields: primaryFields}

	// insert thrice, but still have the same data
	df.Insert(data)
	df.Insert(data)
	df.Insert(data)

	for i := 0; i < len(primaryFields); i++ {
		if df.pkFields[i] != primaryFields[i] {
			t.Errorf("pkField %d expected value: '%s', got '%s'", i, primaryFields[i], df.pkFields[i])
		}
	}

	for _, col := range expectedCols {
		// test that the column exists
		dfCol, ok := df.cols[col]
		if !ok {
			t.Fatalf("Col '%s' does not exist on the dataframe", col)
		}

		// ensure that column has all keys
		for i, key := range keys {
			if dfCol.keys[i] != key {
				t.Fatalf("Col '%s' had key: '%s' instead of '%s'", col, dfCol.keys[i], key)
			}
		}

		// test that the items in that column exist
		for pos := 0; pos < len(data); pos++ {
			value := dfCol.items[pos]
			if value.value != data[pos][dfCol.Name] {
				t.Fatalf("Value of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, data[pos][dfCol.Name], value.value)
			}

			if value.pk != keys[pos] {
				t.Fatalf("Pk of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, keys[pos], value.pk)
			}
		}		
	}
}

// Insert should add the new records at the end of the dtaframe,
// while initializing the values for the non-existing columns to nil or its equivalent
// for the other prexisting values
func TestInsertNonExistingCols(t *testing.T)  {

	firstData := []map[string]interface{}{
		{"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
		{"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
		{"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
	}
	nextData := []map[string]interface{}{
		{"first name": "Richard", "last name": "Roe", "address": "Nairobi" },
		{"first name": "Reyna", "last name": "Roe", "address": "Nairobi" },
		{"first name": "Ruth", "last name": "Roe", "address": "Kampala" },
	}
	primaryFields := []string{"first name", "last name"}
	firstDataCols := []string{"first name", "last name", "age", "location"}
	nextDataCols := []string{"first name", "last name", "address"}	
	firstKeys := []string{"John_Doe", "Jane_Doe", "Paul_Doe"}
	nextKeys := []string{"Richard_Roe", "Reyna_Roe", "Ruth_Roe"}

	df := Dataframe{pkFields: primaryFields}

	// Insert the two sets of records
	df.Insert(firstData)
	df.Insert(nextData)

	for i := 0; i < len(primaryFields); i++ {
		if df.pkFields[i] != primaryFields[i] {
			t.Errorf("pkField %d expected value: '%s', got '%s'", i, primaryFields[i], df.pkFields[i])
		}
	}

	// ensure that all items have the first data columns 
	// // ensure that the nextData items have nil in place columns that didnot exist in them
	for _, col := range firstDataCols {
		// test that the column exists
		dfCol, ok := df.cols[col]
		if !ok {
			t.Fatalf("Col '%s' does not exist on the dataframe", col)
		}

		// ensure that column has all first data keys
		for i, key := range firstKeys {
			if dfCol.keys[i] != key {
				t.Fatalf("first col '%s' had key: '%s' instead of '%s'", col, dfCol.keys[i], key)
			}
		}

		// ensure that column has all first data keys
		for i, key := range nextKeys {
			if dfCol.keys[i] != key {
				t.Fatalf("next col '%s' had key: '%s' instead of '%s'", col, dfCol.keys[i], key)
			}
		}

		// test that the first items in that column exist
		for pos := 0; pos < len(firstData); pos++ {
			value := dfCol.items[pos]
			if value.value != firstData[pos][dfCol.Name] {
				t.Fatalf("Value of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, firstData[pos][dfCol.Name], value.value)
			}

			if value.pk != firstKeys[pos] {
				t.Fatalf("Pk of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, firstKeys[pos], value.pk)
			}
		}	

		// test that the next items in that column exist
		for pos := len(firstData); pos < len(nextData); pos++ {
			value := dfCol.items[pos]
			posInNextItems := pos - len(firstData)
			if value.value != nextData[posInNextItems][dfCol.Name] {
				t.Fatalf("Value of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, nextData[posInNextItems][dfCol.Name], value.value)
			}

			if value.pk != nextKeys[posInNextItems] {
				t.Fatalf("Pk of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, nextKeys[posInNextItems], value.pk)
			}
		}

	}

	// ensure that all items have the next data columns
	// // ensure that the firstData items have nil in place columns that didnot exist in them
	for _, col := range nextDataCols {
		// test that the column exists
		dfCol, ok := df.cols[col]
		if !ok {
			t.Fatalf("Col '%s' does not exist on the dataframe", col)
		}

		// ensure that column has all first data keys
		for i, key := range firstKeys {
			if dfCol.keys[i] != key {
				t.Fatalf("first col '%s' had key: '%s' instead of '%s'", col, dfCol.keys[i], key)
			}
		}

		// ensure that column has all first data keys
		for i, key := range nextKeys {
			if dfCol.keys[i] != key {
				t.Fatalf("next col '%s' had key: '%s' instead of '%s'", col, dfCol.keys[i], key)
			}
		}

		// test that the first items in that column exist
		for pos := 0; pos < len(firstData); pos++ {
			value := dfCol.items[pos]
			if value.value != firstData[pos][dfCol.Name] {
				t.Fatalf("Value of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, firstData[pos][dfCol.Name], value.value)
			}

			if value.pk != firstKeys[pos] {
				t.Fatalf("Pk of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, firstKeys[pos], value.pk)
			}
		}	

		// test that the next items in that column exist
		for pos := len(firstData); pos < len(nextData); pos++ {
			value := dfCol.items[pos]
			posInNextItems := pos - len(firstData)
			if value.value != nextData[posInNextItems][dfCol.Name] {
				t.Fatalf("Value of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, nextData[posInNextItems][dfCol.Name], value.value)
			}

			if value.pk != nextKeys[posInNextItems] {
				t.Fatalf("Pk of item %d for field '%s' expected: '%s'; got '%s'", pos, dfCol.Name, nextKeys[posInNextItems], value.pk)
			}
		}
	}	
}