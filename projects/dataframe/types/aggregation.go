package types

var (
	MAX aggregateFunc = max
	MIN aggregateFunc = min
	SUM aggregateFunc = sum
	MEAN aggregateFunc = mean
	// MODE
	// RANGE
	// PERCENTILE
)

// map of column name and the aggregateFunc function to apply to its values
type aggregation map[string]aggregateFunc

// aggregation function to convert array of values into single value especially during grouping
type aggregateFunc func([]interface{}) interface{}


// Aggregation function to get the maximum value in the list of values
func max(values []interface{}) interface{} {
	var a interface{} = nil
	for _, v := range values {
		if v == nil { continue }
		if a == nil { a = v }
		
		switch v := v.(type) {
		case int:
			if a.(int) < v { a = v }
		case int8:
			if a.(int8) < v { a = v }
		case int16:
			if a.(int16) < v { a = v }
		case int32:
			if a.(int32) < v { a = v }
		case int64:
			if a.(int64) < v { a = v }
		case float32:
			if a.(float32) < v { a = v }
		case float64:
			if a.(float64) < v { a = v }
		case string:
			if a.(string) < v { a = v }
		}			
	}

	return a
}

// Aggregation function to get the minimum value in the list of values
func min(values []interface{}) interface{} {
	var a interface{} = nil
	for _, v := range values {
		if v == nil { continue }
		if a == nil { a = v }

		switch v := v.(type) {
		case int:
			if a.(int) > v { a = v }
		case int8:
			if a.(int8) > v { a = v }
		case int16:
			if a.(int16) > v { a = v }
		case int32:
			if a.(int32) > v { a = v }
		case int64:
			if a.(int64) > v { a = v }
		case float32:
			if a.(float32) > v { a = v }
		case float64:
			if a.(float64) > v { a = v }
		case string:
			if a.(string) > v { a = v }
		}			
	}

	return a
}

// Aggregation function to get the sum of the values. 
// It returns nil if there are some nil values
func sum(values []interface{}) interface{} {
	var a interface{} = nil
	for _, v := range values {
		if v == nil { continue }
		if a == nil { 
			a = v 
			continue
		}

		switch v := v.(type) {
		case int:
			a = a.(int) + v
		case int8:
			a = a.(int8) + v
		case int16:
			a = a.(int16) + v
		case int32:
			a = a.(int32) + v
		case int64:
			a = a.(int64) + v
		case float32:
			a = a.(float32) + v
		case float64:
			a = a.(float64) + v
		default:
			a = nil
		}			
	}

	return a
}

// Aggregation function to get the mean value in the list of values 
// It returns nil if there are some nil values
func mean(values []interface{}) interface{} {
	var a interface{} = nil
	for _, v := range values {
		if v == nil { continue }
		if a == nil { 
			a = v 
			continue
		}

		switch v := v.(type) {
		case int:
			a = a.(int) + v
		case int8:
			a = a.(int8) + v
		case int16:
			a = a.(int16) + v
		case int32:
			a = a.(int32) + v
		case int64:
			a = a.(int64) + v
		case float32:
			a = a.(float32) + v
		case float64:
			a = a.(float64) + v
		default:
			a = nil
		}			
	}

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
