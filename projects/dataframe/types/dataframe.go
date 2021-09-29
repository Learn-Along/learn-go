package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

type Dataframe struct {
	cols map[string]Column;
	pkFields []string;
	index map[interface{}]int;
}

// Constructs a Dataframe from an array of maps and returns a pointer to it
func FromArray(records []map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]Column{},
		index: map[interface{}]int{},
	}

	// FIXME: what if we just generate the primary keys and the col items in one loop and just update
	// the created dataframe's cols and index. There would be no need for even calling normalizeCols
	for _, record := range records {
		err := df.insertRecord(record)
		if err != nil {
			return nil, err
		}
	}

	df.normalizeCols(nil)
	return &df, nil
}

// Constructs a Dataframe from a map of maps and returns a pointer to it
func FromMap(records map[interface{}]map[string]interface{}, primaryFields []string) (*Dataframe, error) {
	df := Dataframe{
		pkFields: primaryFields,
		cols: map[string]Column{},
		index: map[interface{}]int{},
	}

	// FIXME: what if we just generate the primary keys and the col items in one loop and just update
	// the created dataframe's cols and index. There would be no need for even calling normalizeCols
	for _, record := range records {
		err := df.insertRecord(record)
		if err != nil {
			return nil, err
		}
	}

	df.normalizeCols(nil)

	return &df, nil
}

// Inserts items passed as a list of maps into the Dataframe,
// It will overwrite any record whose primary field values match with the new records
func (d *Dataframe) Insert(records []map[string]interface{}) error {
	d.defragmentize()

	// FIXME:
	// To quicken this even further, we could transpose the matrix at this point 
	// and have slices corresponding to each column. These can then be bulk inserted into the columns.
	for _, record := range records {
		err := d.insertRecord(record)
		if err != nil {
			// FIXME: This should probably rollback; might need to make snapshots
			return err
		}
	}	

	d.normalizeCols(nil)
	return nil
}

// Deletes the items that fulfill the filters
func (d *Dataframe) Delete(filter filterType) error {
	count := d.Count()
	indicesToDelete := make([]int, count)
	pkIndices := d.getIndicesInOrder()
	keys := d.Keys()

	counter := 0
	for i, shouldDelete := range filter {
		// FIXME:
		// aside from the mutation delete(d.index,..) which might have race conditions,
		// these others could be done concurrently
		// as they don't affect themselves.
		if shouldDelete && i < count {
			indicesToDelete[counter] = pkIndices[i]
			counter++

			// FIXME:
			// remove this from here. Look for a bulk way of removing keys from a map quickly
			delete(d.index, keys[i])
		}		
	}

	// delete the items in each col 
	for _, col := range d.cols {
		// FIXME:
		// These again can be done concurrently since the data is saved in separate columns.
		col.deleteMany(indicesToDelete[:counter])
	}

	// defragmentize the pks, index and cols 
	d.defragmentize()

	return nil
}

// Updates the items that fulfill the given filters with the new value
func (d *Dataframe) Update(filter []bool, value map[string]interface{}) error  {
	count := d.Count()
	sizeOfValue := len(value)
	indicesToUpdate := make([]int, count)
	pkIndices := d.getIndicesInOrder()
	valueCopy := make(map[string]interface{}, sizeOfValue)
	pkFieldMap := d.getPkFieldMap()

	counter := 0
	for i, shouldUpdate := range filter {
		// FIXME: Concurrency should be possible here, possibly by ranging over 0 to len(filter)
		// The pkIndex could be pushed to a channel and another goroutine just updates that index
		if shouldUpdate && i < count {
			indicesToUpdate[counter] = pkIndices[i]		
			counter++	
		}		
	}

	for k, v := range value {
		// FIXME: Concurrency possible
		if _, ok := pkFieldMap[k]; !ok {
			valueCopy[k] = v
		}
	}

	// update only upto counter
	// This could a range over a channel instead...see FIXME at "for i, shouldUpdate := range filter"
	for _, pkIndex := range indicesToUpdate[:counter] {
		// FIXME: concurrenyc is possible for this inner loop
		for colName, v := range valueCopy {		
			d.Col(colName).insert(pkIndex, v)			
		}		
	}

	return nil
}

// Selects a given number of fields, and returns a query instance of the same
func (d *Dataframe) Select(fields ...string) *query {
	// Creates a new query with this df and one SELECT action in the ops list
	return &query{df: d, ops: []action{{_type: SELECT_ACTION, payload: fields}}}
}

// Merges the dataframes dfs to d
func (d *Dataframe) Merge(dfs ...*Dataframe) error {
	for _, df := range dfs {
		// FIXME: Is it possible to merge without having to change to row-wise structure first.
		// that is basing on the assumption that columnar is more efficient as we claimed
		records, err := df.ToArray()
		if err != nil {
			return err
		}

		err = d.Insert(records)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns the number of actual active items
func (d *Dataframe) Count() int {
	return len(d.index)
}

// Copies the dataframe and returns the new copy
func (d *Dataframe) Copy() (*Dataframe, error) {
	// FIXME: Why not just copy the columns, the index, and the pkFields
	// Having to call ToArray is very inefficient 
	records, err := d.ToArray()
	if err != nil {
		return nil, err
	}

	return FromArray(records, d.pkFields)
}

// Converts that dataframe into a slice of records (maps). If selectedFields is a non-empty slice 
// the fields are limited only to the passed fields
func (d *Dataframe) ToArray(selectedFields ...string) ([]map[string]interface{}, error) {
	count := d.Count()
	pkIndices := d.getIndicesInOrder()
	data := make([]map[string]interface{}, count)
	cols := map[string]Column{}

	for _, field := range selectedFields {
		// FIXME: concurrency possible
		if val, ok := d.cols[field]; ok {
			cols[field] = val
		}
	}

	if len(cols) == 0 {
		cols = d.cols
	}

	for i, pkIndex := range pkIndices {
		record := map[string]interface{}{}

		// FIXME: The column names are unique, the columns are independent, concurrency is thus possible
		for _, col := range cols {
			if i < col.Len() {
				record[col.Name()] = col.ItemAt(pkIndex)
			} else {
				record[col.Name()] = nil
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
		// FIXME: can be done concurrently
		delete(d.cols, k)
	}

	// clear the index
	for k := range d.index {
		// FIXME: Can be done concurrently
		delete(d.index, k)
	}	
}

// Gets the pointer to a given column, or creates it if it does not exist
func (d *Dataframe) Col(name string) Column {
	col := d.cols[name]

	// if col == nil {
	// 	newCol := ColumnStruct{Name: name, items: map[int]interface{}{}, Dtype: ObjectType}
	// 	d.cols[name] = &newCol 
	// 	return &newCol
	// }

	return col
}

// Access method to return the keys in order
func (d *Dataframe) Keys() []string {
	count := len(d.index)
	orderedKeyMap := make(OrderedStringMapType, count)

	for key, i := range d.index {
		orderedKeyMap[i] = key.(string)
	}

	return orderedKeyMap.ToSlice().([]string)
}

// access method to return all column names
func (d *Dataframe) ColumnNames() []string {
	count := len(d.cols)
	names := make([]string, count)

	i := 0
	for key := range d.cols {
		names[i] = key
		i++
	}

	return names
}

// Pretty prints the record in this dataframe
func (d *Dataframe) PrettyPrintRecords() error {
	// FIXME:
	// Is it possible to print the data as a table instead of row-wise data,
	// so as to give the actual picture of how the data is stored.
	// e.g.
	// ----------------------------------------
	// | Col 1   | Col 2   | Col 3 | Col 4    |
	// ----------------------------------------
	// | foo     | 45      | 90    | hyu      |
	data, err := d.ToArray()
	if err != nil {
		return err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return utils.PrettyPrintJSON(dataJSON)
}

// Returns the indices of the pks that have not been deleted, i.e. that have no nil
func (d *Dataframe) getIndicesInOrder() []int {
	count := len(d.index)
	indices := make([]int, count)

	// FIXME:What if this were converted into a range from zero to count, instead of range d.index
	// wouldn't this allow for concurrently updating the indices slice as no data races would be expected.
	counter := 0
	for _, i := range d.index {
		indices[counter] = i
		counter++
	}

	sort.Slice(indices, func(i, j int) bool {	return indices[i] < indices[j]	})
	return indices
}

// Inserts a single record
func (d *Dataframe) insertRecord(record map[string]interface{}) error {
	key, err := createKey(record, d.pkFields)
	if err != nil {
		return fmt.Errorf("failed to create key for %v using field %v", record, d.pkFields)
	}

	row, ok := d.index[key]
	if !ok {
		row = len(d.index)
		d.index[key] = row
	}		

	for fieldName, value := range record {
		// FIXME:
		// to take advantage of having values in separate columns,
		// these values can be saved concurrently
		col := d.Col(fieldName)			
		col.insert(row, value)
	}

	return nil
}

// Fills up the columns with the given value to reach a given length for all columns
func (d *Dataframe) normalizeCols(defaultValue interface{})  {
	pkIndices := d.getIndicesInOrder()
	finalLength := len(pkIndices)

	for _, col := range d.cols {
		colLength := col.Len()
		
		for i := colLength; i < finalLength; i++ {
			pkIndex := pkIndices[i]
			col.insert(pkIndex, defaultValue)
		}
	}
}

// Converts the primary key field list to a map for easy checking against, to see if field is pkField or not
func (d *Dataframe) getPkFieldMap() map[string]struct{} {
	_map := make(map[string]struct{}, len(d.pkFields))

	for _, field := range d.pkFields {
		// FIXME: This can be done concurrently
		_map[field] = struct{}{}
	}

	return _map
}

// Creates a Key to be used to identify the given record
func createKey(record map[string]interface{}, primaryFields []string) (string, error)  {
	key := ""
	separator := "_"

	// FIXME: Could using strings.Join more expressive of what is actually being done here? Try that.
	for _, pkField := range primaryFields {
		if value, ok := record[pkField]; ok {
			key += fmt.Sprintf("%v_", value)
		} else {
			return "", fmt.Errorf("key error: %s in record %v", pkField, record)
		}
	}
	
	return strings.TrimRight(key, separator), nil
}

// reorders pks and indices and the cols
func (d *Dataframe) defragmentize()  {
	pkIndices := d.getIndicesInOrder()
	keys := d.Keys()

	for _, col := range d.cols {
		// FIXME:
		// These columns are independent. Their defragmentation can be done concurrently
		col.Defragmentize(pkIndices)
	}

	for newRow, key := range keys {
		// FIXME:
		// This could be done concurrently
		// even if two keys were alike, this is supposed to be an index, and thus only one key should be present
		d.index[key] = newRow
	}
}

// Filters this dataframe and returns the filtered **copy** of this dataframe
func (d *Dataframe) getFilteredDf(filter filterType) (*Dataframe, error) {
	newDf, err := d.Copy()
	if err != nil {
		return nil, err
	}

	if filter == nil {
		return newDf, nil
	}

	// toggle the values in the filter, and delete the unwanted items
	for i, v := range filter {
		// FIXME:
		// This can be done concurrently
		filter[i] = !v
	}

	err = newDf.Delete(filter)
	if err != nil {
		return nil, err
	}

	return newDf, nil

}

// Groups this dataframe, basing on the groupbyOption passed, and returns a new grouped Dataframe copy
func (d *Dataframe) getGroupedDf(gopt *groupByOption) (*Dataframe, error) {
	aggs := mergeAggregations(gopt.aggs)
	groupedData := map[string][]map[string]interface{}{}
	mergedRecords := []map[string]interface{}{}
	index := []string{}

	// FIXME:
	// Is it possible to group these items without first going back to row-wise structure?
	// Afterall, we claimed it was more efficient to group when the data is columnar, or was that wishful thinking?
	records, err := d.ToArray()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		key, err := createKey(record, gopt.fields)
		if err != nil {
			return nil, err
		}

		if data, ok := groupedData[key]; ok {
			groupedData[key] = append(data, record)
		} else {
			// index maintains order of the groups
			index = append(index, key)
			
			mergedRecord := map[string]interface{}{}
			for _, field := range gopt.fields {
				mergedRecord[field] = record[field]
			}
			mergedRecords = append(mergedRecords, mergedRecord)

			groupedData[key] = []map[string]interface{}{record}
		}
	}


	for i, key := range index {
		mergedRecord := mergedRecords[i]
		df, err := FromArray(groupedData[key], d.pkFields)
		if err != nil {
			return nil, err
		}

		for field, aggFunc := range aggs {
			mergedRecord[field] = aggFunc(df.Col(field).Items())
		}

		mergedRecords[i] = mergedRecord
	}	
	
	return FromArray(mergedRecords, gopt.fields)
}

// Orders the items in the columns of this dataframe basing on the sort options passed.
func (d *Dataframe) getSortedDf(options... sortOption) (*Dataframe, error) {
	// FIXME: 
	// Is it possible to sort this dataframe without going back to row-wise structure?
	// Afterall we claimed it was more efficient to do so when the data was columnar, or was that wishful thinking?
	records, err := d.ToArray()
	if err != nil {
		return nil, err
	}

	if options == nil {
		return d.Copy()
	}
		
	sort.SliceStable(records, func(i, j int) bool {
		prev := records[i]
		next := records[j]

		for _, opt := range options {
			for field, order := range opt {	
				nextValue := next[field]
				prevValue := prev[field]

				if nextValue == prevValue {
					continue
				}

				if prevValue == nil {
					// nils will be pushed up by default
					return order == ASC
				}

				if nextValue == nil {
					// nils will be pushed up by default
					return order == DESC
				}

				switch p := prevValue.(type) {
				case string:
					if order == ASC {
						return p < nextValue.(string)
					} else {
						return p > nextValue.(string)
					}
				default:
					prevAsFloat := convertToFloat64(prevValue)
					nextAsFloat := convertToFloat64(nextValue)

					if order == ASC {
						return prevAsFloat < nextAsFloat
					} else {
						return prevAsFloat > nextAsFloat
					}
				}	
				
			}
		}

		return true
	})
	

	return FromArray(records, d.pkFields)
}

// Applys the given rowWiseFunc functions on the dataframe
func (d *Dataframe) apply(rowWiseFuncMap map[string][]rowWiseFunc) error {
	for field, txs := range rowWiseFuncMap {
		// FIXME: This second for loop works on each field/column at a time,
		// and so this should be done concurrently
		for _, tx := range txs {
			if col, ok := d.cols[field]; ok {
				// FIXME: This third for loop works on individual items,
				// and so this too can be done concurrently
				switch col.GetDatatype() {
				case IntType:
					for i, v := range col.Items().([]int) {
						col.insert(i, tx(v))
					}

				case FloatType:
					for i, v := range col.Items().([]int) {
						col.insert(i, tx(v))
					}

				case StringType:
					for i, v := range col.Items().([]int) {
						col.insert(i, tx(v))
					}

				case BoolType:
					for i, v := range col.Items().([]int) {
						col.insert(i, tx(v))
					}
				}
			}	
		}			
	}
	return nil
}