package types

import (
	"fmt"
	"strings"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

type Dataframe struct {
	cols map[string]*Column;
	pkFields []string;
	index map[interface{}]int;
	pks OrderedMap;
}

// Constructs a Dataframe from an array of maps and returns a pointer to it
func FromArray(records []map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]*Column{},
		index: map[interface{}]int{},
		pks: OrderedMap{},
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
		index: map[interface{}]int{},
		pks: OrderedMap{},
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
		newCol := Column{Name: name, items: map[int]interface{}{}, Dtype: ObjectType}
		d.cols[name] = &newCol 
		return &newCol
	}

	return col
}

// Access method to return the keys in order
func (d *Dataframe) Keys() []string {
	return utils.ConvertToStringSlice(d.pks.ToSlice(), true)
}

// Returns the indices of the pks that have not been deleted, i.e. that have no nil
func (d *Dataframe) getNonNilPkIndices() []int {
	count := len(d.pks)
	indices := make([]int, count)

	counter := 0
	for i := 0; i < count && d.pks[i] != nil; i++ {
		indices[counter] = i
		counter++
	}

	return indices[:counter]

}

// access method to return all column names
func (d *Dataframe) ColumnNames() []string {
	count := len(d.cols)
	names := make([]string, count)

	i := 0
	for _, col := range d.cols {
		names[i] = col.Name
		i++
	}

	return names
}

// Returns the number of actual active items
func (d *Dataframe) Count() int {
	return len(d.index)
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

	finalLength := d.Count()
	d.fillUpCols(finalLength, nil)

	return nil
}

// reorders pks and indices and the cols
func (d *Dataframe) defragmentize()  {
	pkIndices := d.getNonNilPkIndices()

	for _, col := range d.cols {
		col.items.Defragmentize(pkIndices)
	}

	for newRow, oldRow := range pkIndices {
		pk := d.pks[oldRow]
		d.index[pk] = newRow
	}

	d.pks.Defragmentize(pkIndices)
}

// Deletes the items that fulfill the filters
func (d *Dataframe) Delete(filter Filter) error {
	colIndicesToDelete := []int{}
	count := d.Count()
	pkIndices := d.getNonNilPkIndices()

	for i, flag := range filter {
		// delete where flag is true and i is in range of the index
		if flag && i < count {
			// delete the pk from index
			pkIndex := pkIndices[i] 
			colIndicesToDelete = append(colIndicesToDelete, pkIndex)		
			
			pk := d.pks[pkIndex]
			delete(d.index, pk)
			// save nil in d.pks for index i	
			d.pks[pkIndex] = nil	
		}		
	}

	// delete the items in each col 
	for _, col := range d.cols {
		col.deleteMany(colIndicesToDelete)
	}

	// defragmentize the pks, index and cols 
	d.defragmentize()

	return nil
}

// Updates the items that fulfill the given filters with the new value
func (d *Dataframe) Update(filters []bool, value map[string]interface{}) error  {
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

// Converts that dataframe into a slice of records (maps)
func (d *Dataframe) ToArray() ([]map[string]interface{}, error) {
	count := d.Count()
	pkIndices := d.getNonNilPkIndices()
	data := make([]map[string]interface{}, count)
	

	for i, pkIndex := range pkIndices {
		record := map[string]interface{}{}

		for _, col := range d.cols {
			if i < len(col.items) {
				record[col.Name] = col.items[pkIndex]
			} else {
				record[col.Name] = nil
			}			
		}
		
		data[i] = record
	}

	return data, nil
}

// Clears all the data held by the dataframe except the primary key fields
func (d *Dataframe) Clear() {
	// clear the cols
	for k := range d.cols {
		delete(d.cols, k)
	}

	// clear the pks and index
	for i, k := range d.pks {
		delete(d.pks, i)
		delete(d.index, k)
	}	
}

// Inserts a single record
func (d *Dataframe) insertRecord(record map[string]interface{}) error {
	key, err := createKey(record, d.pkFields)
	if err != nil {
		return fmt.Errorf("failed to create key for %v using field %v", record, d.pkFields)
	}

	row, ok := d.index[key]; 
	if !ok {
		row = len(d.pks)
		d.index[key] = row
		d.pks[row] = key
	}		

	for fieldName, value := range record {
		col := d.Col(fieldName)			
		col.insert(row, value)
	}

	return nil
}

// Fills up the columns with the given value to reach a given length for all columns
func (d *Dataframe) fillUpCols(finalLength int, value interface{})  {
	pkIndices := d.getNonNilPkIndices()

	for _, col := range d.cols {
		colLength := len(col.items)
		
		for i := colLength; i < finalLength; i++ {
			pkIndex := pkIndices[i]
			col.insert(pkIndex, value)
		}
	}
}