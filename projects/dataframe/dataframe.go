package dataframe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/learn-along/learn-go/projects/dataframe/internal"
	"github.com/learn-along/learn-go/projects/dataframe/types"
	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/newqf"
)

type View interface {
	ItemAt(i int) bool
	Len() int
	Slice() []bool
}

//// external
type FieldConfig struct {
	Name string
	Type Dtype
}


const (
	IntType Dtype = iota
	Float64Type
	StringType
	BooleanType
)

type Dtype int

////////

/*
* The Dataframe that exposes methods for selection, aggregation, sorting, grouping, filtering etc.
 */
type Dataframe struct {
	q qframe.QFrame
	pkFields []string
	index map[string]int
	fieldOrder []FieldConfig
	columnOrder []string
} 

/*
* Returns a QFrame pointer when given an array/slice of maps of string, value
*
* supported data types of the values include:
* - string
* - int
* - float64
* - bool
 */
func FromArray(records []map[string]interface{}, primaryFields []string, fieldOrder []FieldConfig) (*Dataframe, error) {
	df := Dataframe{pkFields: primaryFields, index: map[string]int{}, fieldOrder: fieldOrder}
	noOfFields := len(fieldOrder)
	numberOfRecords := len(records)
	data := make(map[string]types.DataSlice, noOfFields)
	df.columnOrder = make([]string, 0, noOfFields)

	for i, record := range records {
		key, err := createKey(record, primaryFields)
		if err != nil {
			return nil, fmt.Errorf("error creating key: %s", err)
		}

		if v, ok := df.index[key]; ok {
			return nil, fmt.Errorf("primary fields '%v' are the same for \n%v \nand \n%v", primaryFields, v, record)
		} else {
			df.index[key] = i
		}

		err = addRecordToDataList(data, fieldOrder, numberOfRecords, record)
		if err != nil {
			return nil, err
		}
	}

	for _, fconfig := range fieldOrder {
		df.columnOrder = append(df.columnOrder, fconfig.Name)
	}

	df.q = qframe.New(data, newqf.ColumnOrder(df.columnOrder...))
	return &df, nil
}



/*
* Returns a QFrame pointer when given an map of maps of string, value
*
* supported data types of the values include:
* - string
* - int
* - float64
* - bool
 */
 func FromMap(records map[interface{}]map[string]interface{}, primaryFields []string, fieldOrder []FieldConfig) (*Dataframe, error) {
	df := Dataframe{pkFields: primaryFields, index: map[string]int{}, fieldOrder: fieldOrder}
	noOfFields := len(fieldOrder)
	numberOfRecords := len(records)
	data := make(map[string]types.DataSlice, noOfFields)
	df.columnOrder = make([]string, 0, noOfFields)

	count := 0
	for _, record := range records {
		key, err := createKey(record, primaryFields)
		if err != nil {
			return nil, fmt.Errorf("error creating key: %s", err)
		}

		if v, ok := df.index[key]; ok {
			return nil, fmt.Errorf("primary fields '%v' are the same for \n%v \nand \n%v", primaryFields, v, record)
		} else {
			df.index[key] = count
		}

		err = addRecordToDataList(data, fieldOrder, numberOfRecords, record)
		if err != nil {
			return nil, err
		}

		count++
	}

	for _, fconfig := range fieldOrder {
		df.columnOrder = append(df.columnOrder, fconfig.Name)
	}

	df.q = qframe.New(data, newqf.ColumnOrder(df.columnOrder...))
	return &df, nil
}

// Inserts items passed as a list of maps into the Dataframe,
// It will overwrite any record whose primary field values match with the new records
func (d *Dataframe) Insert(records []map[string]interface{}) error {
	noOfFields := len(d.fieldOrder)
	data := make(map[string]types.DataSlice, noOfFields)
	allRecords, err := d.ToArray()
	if err != nil {
		return err
	}

	for _, record := range records {
		key, err := createKey(record, d.pkFields)
		if err != nil {
			return err
		}

		if pos, ok := d.index[key]; ok {
			// overwrite
			allRecords[pos] = record
		} else {
			allRecords = append(allRecords, record)
		}		
	}

	numberOfRecords := len(allRecords)

	for i, record := range allRecords {
		key, err := createKey(record, d.pkFields)
		if err != nil {
			return fmt.Errorf("error creating key: %s", err)
		}

		d.index[key] = i

		err = addRecordToDataList(data, d.fieldOrder, numberOfRecords, record)
		if err != nil {
			return err
		}		
	}

	d.q = qframe.New(data, newqf.ColumnOrder(d.columnOrder...))
	return nil
}

// Deletes the items that fulfill the filters
func (d *Dataframe) Delete(filter internal.FilterType) error {
	return nil
}

// Updates the items that fulfill the given filters with the new value
func (d *Dataframe) Update(filter internal.FilterType, value map[string]interface{}) error  {
	return nil
}

// Selects a given number of fields, and returns a query instance of the same
func (d *Dataframe) Select(fields ...string) qframe.QFrame {
	return qframe.QFrame{}
}


// Merges the dataframes dfs to d
func (d *Dataframe) Merge(dfs ...*Dataframe) error {
	return nil
}

// Returns the number of actual active items
func (d *Dataframe) Count() int {
	return 0
}

// Copies the dataframe and returns the new copy
func (d *Dataframe) Copy() (*Dataframe, error) {
	return nil, nil
}

// Converts that dataframe into a slice of records (maps). If selectedFields is a non-empty slice 
// the fields are limited only to the passed fields
func (d *Dataframe) ToArray(selectedFields ...string) ([]map[string]interface{}, error) {	
	q := d.q 

	if len(selectedFields) > 0 {
		q = d.q.Select(selectedFields...)
	}
	
	w := new(bytes.Buffer)
	err := q.ToJSON(w)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal(w.Bytes(), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Clears all the data held by the dataframe except the primary key fields
func (d *Dataframe) Clear() {
}

// Gets the pointer to a given column, or creates it if it does not exist
func (d *Dataframe) Col(name string) View {
	return qframe.BoolView{}
}

// Access method to return the keys in order
func (d *Dataframe) Keys() []string {
	keys := make([]string, len(d.index))

	for k, v := range d.index {
		keys[v] = k
	}

	return keys
}

// access method to return all column names
func (d *Dataframe) ColumnNames() []string {
	return d.q.ColumnNames()
}

// Pretty prints the record in this dataframe
func (d *Dataframe) String() string {
	return ""
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


// Adds a record to a data map of DataSlice
func addRecordToDataList(data map[string]types.DataSlice, fieldOrder []FieldConfig, finalColumnLength int, record map[string]interface{}) error {
	for _, fconfig := range fieldOrder {
		switch fconfig.Type {
		case IntType:
			if data[fconfig.Name] == nil {
				data[fconfig.Name] = make([]int, 0, finalColumnLength)
			}

			data[fconfig.Name] = append(data[fconfig.Name].([]int), record[fconfig.Name].(int))
		case Float64Type:
			if data[fconfig.Name] == nil {
				data[fconfig.Name] = make([]float64, 0, finalColumnLength)
			}

			data[fconfig.Name] = append(data[fconfig.Name].([]float64), record[fconfig.Name].(float64))
		case StringType:
			if data[fconfig.Name] == nil {
				data[fconfig.Name] = make([]string, 0, finalColumnLength)	
			}

			data[fconfig.Name] = append(data[fconfig.Name].([]string), record[fconfig.Name].(string))
		case BooleanType:
			if data[fconfig.Name] == nil {
				data[fconfig.Name] = make([]bool, 0, finalColumnLength)	
			}

			data[fconfig.Name] = append(data[fconfig.Name].([]bool), record[fconfig.Name].(bool))
		default:
			return fmt.Errorf("'%v' dtype is unknown", fconfig.Type)	
		}
	}

	return nil
}
