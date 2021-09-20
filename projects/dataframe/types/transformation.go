package types

// map of column name and the array function to apply to its values
type transformation map[string]rowWiseFunc

// function that transforms each element into another element
type rowWiseFunc func(interface{}) interface{}

// Merges a slice of transformations into a map of lists of rowWiseFunc functions
func mergeTransformations(aggs []transformation) map[string][]rowWiseFunc {
	res := map[string][]rowWiseFunc{}

	for _, agg 	:= range aggs {
		for key, v := range agg {
			prev, ok := res[key]
			if !ok {
				prev = []rowWiseFunc{}
			} 

			res[key] = append(prev, v)
		}
	}

	return res
}