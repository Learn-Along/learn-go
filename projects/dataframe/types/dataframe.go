package types

import "fmt"

type Dataframe struct {
	cols map[string]*Column;
	pkFields []string;
}

// Constructs a Dataframe from an array of maps and returns a pointer to it
func FromArray(records []map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{pkFields: primaryFields, cols: map[string]*Column{}}

	for i, record := range records {
		key, err := createKey(record, primaryFields)
		if err != nil {
			return nil, fmt.Errorf("failed to create key for %v using field %v", record, primaryFields)
		}

		for fieldName, value := range record {
			col := df.Col(fieldName)			
			col.insert(i, key, value)
		}
	}

	return &df, nil
}

// Constructs a Dataframe from a map of maps and returns a pointer to it
func FromMap(records map[interface{}]map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{pkFields: primaryFields, cols: map[string]*Column{}}

	i := 0
	for _, record := range records {
		key, err := createKey(record, primaryFields)
		if err != nil {
			return nil, fmt.Errorf("failed to create key for %v using field %v", record, primaryFields)
		}

		for fieldName, value := range record {
			col := df.Col(fieldName)			
			col.insert(i, key, value)
		}

		i++
	}

	return &df, nil
}

// Creates a Key to be used to identify the given record
func createKey(record map[string]interface{}, primaryFields []string) (string, error)  {
	return "", nil
}


// Gets the pointer to a given column, or creates it if it does not exist
func (d *Dataframe) Col(name string) *Column {
	col := d.cols[name]

	if col == nil {
		newCol := Column{Name: name, Items: []interface{}{}, keys: []string{}, Dtype: ObjectType}
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
