package types

import (
	"fmt"
	"strconv"
)

var (
	MAX aggregateFunc = getMax
	MIN aggregateFunc = getMin
	SUM aggregateFunc = getSum
	MEAN aggregateFunc = getMean
	// MODE
	// RANGE
	// PERCENTILE(int)
)

// map of column name and the aggregateFunc function to apply to its values
type aggregation map[string]aggregateFunc

// aggregation function to convert array of values into single value especially during grouping
type aggregateFunc func([]interface{}) interface{}


// Aggregation function to get the maximum value in the list of values
func getMax(values []interface{}) interface{} {
	var a interface{} = nil

	defer func() {
		if r := recover(); r != nil {
			a = nil
		}
	}()

	for _, v := range values {
		if v == nil { continue }
		if a == nil { 
			isStr := false
			if a, isStr = v.(string); !isStr {
				a = convertToFloat64(v)
			}
		}
		
		switch v := v.(type) {
		case int:	
			if a.(float64) < float64(v) { a = float64(v) }
		case int8:
			if a.(float64) < float64(v) { a = float64(v) }
		case int16:
			if a.(float64) < float64(v) { a = float64(v) }
		case int32:
			if a.(float64) < float64(v) { a = float64(v) }
		case int64:
			if a.(float64) < float64(v) { a = float64(v) }
		case float32:
			if a.(float64) < float64(v) { a = float64(v) }
		case float64:
			if a.(float64) < float64(v) { a = float64(v) }
		case string:
			if a.(string) < v { a = v }
		}			
	}

	return a
}

// Aggregation function to get the minimum value in the list of values
func getMin(values []interface{}) interface{} {
	var a interface{} = nil

	defer func() {
		if r := recover(); r != nil {
			a = nil
		}
	}()

	for _, v := range values {
		if v == nil { continue }
		if a == nil { 
			isStr := false
			if a, isStr = v.(string); !isStr {
				a = convertToFloat64(v)
			}
		}

		switch v := v.(type) {
		case int:	
			if a.(float64) > float64(v) { a = float64(v) }
		case int8:
			if a.(float64) > float64(v) { a = float64(v) }
		case int16:
			if a.(float64) > float64(v) { a = float64(v) }
		case int32:
			if a.(float64) > float64(v) { a = float64(v) }
		case int64:
			if a.(float64) > float64(v) { a = float64(v) }
		case float32:
			if a.(float64) > float64(v) { a = float64(v) }
		case float64:
			if a.(float64) > float64(v) { a = float64(v) }
		case string:
			if a.(string) > v { a = v }
		}			
	}

	return a
}

// Aggregation function to get the sum of the values
func getSum(values []interface{}) interface{} {
	var a interface{} = nil

	defer func() {
		if r := recover(); r != nil {
			a = nil
		}
	}()

	for _, v := range values {
		if v == nil { continue }
		if a == nil { 
			a = convertToFloat64(v)
			continue
		}

		switch v := v.(type) {
		case int:
			a = a.(float64) + float64(v)
		case int8:
			a = a.(float64) + float64(v)
		case int16:
			a = a.(float64) + float64(v)
		case int32:
			a = a.(float64) + float64(v)
		case int64:
			a = a.(float64) + float64(v)
		case float32:
			a = a.(float64) + float64(v)
		case float64:
			a = a.(float64) + float64(v)
		default:
			return nil
		}			
	}

	return a
}

// Aggregation function to get the mean value in the list of values 
// It returns nil if there are some nil values
func getMean(values []interface{}) interface{} {
	a := getSum(values)

	defer func() {
		if r := recover(); r != nil {
			a = nil
		}
	}()

	if a != nil {
		count := float64(len(values))

		switch a := a.(type) {
		case int:
			return float64(a) / count
		case int8:
			return float64(a) / count
		case int16:
			return float64(a) / count
		case int32:
			return float64(a) / count
		case int64:
			return float64(a) / count
		case float32:
			return float64(a) / count
		case float64:
			return float64(a) / count
		}
	}

	return a
}

// Converts a given value of unknown type to float64
func convertToFloat64(value interface{}) float64 {
	v := fmt.Sprintf("%v", value)
	valueAsFloat, _ := strconv.ParseFloat(v, 64)
	return valueAsFloat
}