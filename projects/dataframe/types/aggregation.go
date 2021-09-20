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
	COUNT aggregateFunc = getCount
	RANGE aggregateFunc = getRange
	// PERCENTILE(int) etc.
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
			if val := float64(v); a.(float64) < val { a = val }
		case int8:
			if val := float64(v); a.(float64) < val { a = val }
		case int16:
			if val := float64(v); a.(float64) < val { a = val }
		case int32:
			if val := float64(v); a.(float64) < val { a = val }
		case int64:
			if val := float64(v); a.(float64) < val { a = val }
		case float32:
			if val := float64(v); a.(float64) < val { a = val }
		case float64:
			if val := float64(v); a.(float64) < val { a = val }
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
			if val := float64(v); a.(float64) > val { a = val }	
		case int8:
			if val := float64(v); a.(float64) > val { a = val }	
		case int16:
			if val := float64(v); a.(float64) > val { a = val }	
		case int32:
			if val := float64(v); a.(float64) > val { a = val }	
		case int64:
			if val := float64(v); a.(float64) > val { a = val }	
		case float32:
			if val := float64(v); a.(float64) > val { a = val }	
		case float64:
			if val := float64(v); a.(float64) > val { a = val }	
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

// Returns the number of items in the values array
func getCount(values []interface{}) interface{} {
	return len(values)
}

// Returns the difference between the biggest and the smallest value in the values array,
// if all values are numbers (or nil which are ignored), else it returns nil
func getRange(values []interface{}) interface{} {
	var max interface{} = nil
	var min interface{} = nil

	defer func() {
		if r := recover(); r != nil {
			min = nil
			max = nil
		}
	}()

	for _, v := range values {
		if v == nil { continue }
		if max == nil { 
			isStr := false
			if max, isStr = v.(string); !isStr {
				max = convertToFloat64(v)
			}
		}

		if min == nil { 
			isStr := false
			if min, isStr = v.(string); !isStr {
				min = convertToFloat64(v)
			}
		}
		
		switch v := v.(type) {
		case int:	
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case int8:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case int16:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case int32:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case int64:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case float32:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case float64:
			val := float64(v)
			if max.(float64) < val { max = val }
			if min.(float64) > val { min = val }
		case string:
			return nil
		}			
	}

	if min != nil && max != nil {
		return max.(float64) - min.(float64)
	}

	return nil
}

/*
*	Helpers
*/

// Converts a given value of unknown type to float64
func convertToFloat64(value interface{}) float64 {
	v := fmt.Sprintf("%v", value)
	valueAsFloat, _ := strconv.ParseFloat(v, 64)
	return valueAsFloat
}

// Merges a slice of aggregations into one aggregation.
// Inorder to have only one aggregation per column, only the last aggregateFunc passed for that column
// is kept
func mergeAggregations(aggs []aggregation) aggregation {
	res := aggregation{}

	for _, agg 	:= range aggs {
		for key, v := range agg {
			res[key] = v
		}
	}

	return res
}