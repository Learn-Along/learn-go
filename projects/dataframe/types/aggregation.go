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
type aggregateFunc func(ItemSlice) Item


// Aggregation function to get the maximum value in the list of values
func getMax(values ItemSlice) Item {
	switch records := values.(type) {
	case []string:
		var maxV string
		for i, v := range records {
			if maxV < v || i == 0 {
				maxV = v
			}
		}

		return maxV

	case []int:
		var maxV int
		for i, v := range records {
			if maxV < v || i == 0 {
				maxV = v
			}
		}

		return maxV

	case []float64:
		var maxV float64
		for i, v := range records {
			if maxV < v || i == 0 {
				maxV = v
			}
		}

		return maxV

	case []bool:
		var maxV bool = false
		for _, v := range records {
			if v {
				return true
			}
		}

		return maxV

	default:
		return nil		
	}
}

// Aggregation function to get the minimum value in the list of values
func getMin(values ItemSlice) Item {
	switch records := values.(type) {
	case []string:
		var minV string
		for i, v := range records {
			if minV > v || i == 0 {
				minV = v
			}
		}

		return minV

	case []int:
		var minV int
		for i, v := range records {
			if minV > v || i == 0 {
				minV = v
			}
		}

		return minV

	case []float64:
		var minV float64
		for i, v := range records {
			if minV > v || i == 0 {
				minV = v
			}
		}

		return minV

	case []bool:
		var minV bool
		for i, v := range records {
			if i == 0 {
				minV = v
			}

			if !v {
				return false
			}
		}

		return minV

	default:
		return nil		
	}

}

// Aggregation function to get the sum of the values
func getSum(values ItemSlice) Item {
	switch records := values.(type) {
	case []int:
		var a int = 0
		for _, v := range records {
			a += v
		}

		return a

	case []float64:
		var a float64 = 0
		for _, v := range records {
			a += v
		}

		return a

	default:
		return nil		
	}
}

// Aggregation function to get the mean value in the list of values 
// It returns nil if there are some nil values
func getMean(values ItemSlice) Item {
	a := getSum(values)

	if a != nil {
		var count float64 = 0

		switch records := values.(type) {
		case []float64:
			count = float64(len(records))
		case []int:
			count = float64(len(records))
		}

		switch a := a.(type) {
		case int:
			return float64(a) / count
		case float64:
			return float64(a) / count
		}
	}

	return a
}

// Returns the number of items in the values array
func getCount(values ItemSlice) Item {
	switch records := values.(type) {
	case []float64:
		return len(records)
	case []int:
		return len(records)
	case []string:
		return len(records)
	case []bool:
		return len(records)
	}

	return 0
}

// Returns the difference between the biggest and the smallest value in the values array,
// if all values are numbers (or nil which are ignored), else it returns nil
func getRange(values ItemSlice) Item {
	switch records := values.(type) {
	case []int:
		var max int
		var min int 

		for i, v := range records {
			if max < v || i == 0 {
				max = v
			}

			if min > v || i == 0 {
				min = v
			}
		}

		return max - min

	case []float64:
		var max float64
		var min float64 

		for i, v := range records {
			if max < v || i == 0 {
				max = v
			}

			if min > v || i == 0 {
				min = v
			}
		}

		return max - min
		
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