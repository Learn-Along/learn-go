package types

import (
	"fmt"
	"strings"
)

type Dataframe struct {
	cols map[string]*Column;
	pkFields []string;
	index map[string]int;
	pks []string;
}

// Constructs a Dataframe from an array of maps and returns a pointer to it
func FromArray(records []map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]*Column{},
		index: map[string]int{},
		pks: []string{},
	}

	for _, record := range records {
		err := df.insertRecord(record)
		if err != nil {
			return nil, err
		}
	}

	finalLength := len(df.index)
	df.fillUpCols(finalLength, nil)

	return &df, nil
}

// Constructs a Dataframe from a map of maps and returns a pointer to it
func FromMap(records map[interface{}]map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]*Column{},
		index: map[string]int{},
		pks: []string{},
	}

	for _, record := range records {
		err := df.insertRecord(record)
		if err != nil {
			return nil, err
		}
	}

	finalLength := len(df.index)
	df.fillUpCols(finalLength, nil)

	return &df, nil
}

// Creates a Key to be used to identify the given record
func createKey(record map[string]interface{}, primaryFields []string) (string, error)  {
	key := ""
	separator := "_"

	for _, pkField := range primaryFields {
		if value, ok := record[pkField]; ok {
			key += fmt.Sprintf("%s_", value)
		} else {
			return "", fmt.Errorf("key error: %s in record %v", pkField, record)
		}
	}
	
	return strings.TrimRight(key, separator), nil
}


// Gets the pointer to a given column, or creates it if it does not exist
func (d *Dataframe) Col(name string) *Column {
	col := d.cols[name]

	if col == nil {
		newCol := Column{Name: name, Items: []interface{}{}, Dtype: ObjectType}
		d.cols[name] = &newCol 
		return &newCol
	}

	return col
}

// Utility to return all column names
func (d *Dataframe) getColNames() []string {
	names := []string{}
	for _, col := range d.cols {
		names = append(names, col.Name)
	}

	return names
}

// Inserts items passed as a list of maps into the Dataframe,
// It will overwrite any record whose primary field values match with the new records
func (d *Dataframe) Insert(records []map[string]interface{}) error {
	for _, record := range records {
		err := d.insertRecord(record)
		if err != nil {
			// FIXME: This should probably rollback; might need to make snapshots
			return err
		}
	}	

	finalLength := len(d.index)
	d.fillUpCols(finalLength, nil)

	return nil
}

// Deletes the items that fulfill the filter
func (d *Dataframe) Delete(filter Filter) error {
	return nil
}

// Updates the items that fulfill the given filter with the new value
func (d *Dataframe) Update(filter Filter, value map[string]interface{}) error  {
	return nil
}

// Selects a given number of fields, and returns a Query instance of the same
func (d *Dataframe) Select(fields ...string) *Query {
	return nil
}

// Merges the dataframe df to d
func (d *Dataframe) Merge(df *Dataframe) error {
	return nil
}

// Copies the dataframe and returns the new copy
func (d *Dataframe) Copy() (Dataframe, error) {
	return Dataframe{}, nil
}

func (d *Dataframe) ToArray() ([]map[string]interface{}, error) {
	return nil, nil
}

// Frees the memory held by the dataframe by nilling all pointers in it
func (d *Dataframe) Free() error {
	return nil	
}

// Inserts a single record
func (d *Dataframe) insertRecord(record map[string]interface{}) error {
	key, err := createKey(record, d.pkFields)
	if err != nil {
		return fmt.Errorf("failed to create key for %v using field %v", record, d.pkFields)
	}

	row, ok := d.index[key]; 
	if !ok {
		row = len(d.index)
		d.index[key] = row
		d.pks = append(d.pks, key)
	}		

	for fieldName, value := range record {
		col := d.Col(fieldName)			
		col.insert(row, value)
	}

	return nil
}

// Fills up the columns with the given value to reach a given length for all columns
func (d *Dataframe) fillUpCols(finalLength int, value interface{})  {
	for _, col := range d.cols {
		colLength := len(col.Items)
		
		for i := colLength; i < finalLength; i++ {
			col.insert(i, value)
		}
	}
}