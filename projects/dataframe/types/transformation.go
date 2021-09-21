package types

// map of column name and the array function to apply to its values
// type transformation map[string]rowWiseFunc
type transformation struct{k string; v rowWiseFunc}

// function that transforms each element into another element
type rowWiseFunc func(interface{}) interface{}

// Merges a slice of transformations into a map of lists of rowWiseFunc functions
func mergeTransformations(txns []transformation) map[string][]rowWiseFunc {
	res := map[string][]rowWiseFunc{}	

	for _, txn 	:= range txns {
		prev, ok := res[txn.k]
		if !ok {
			prev = []rowWiseFunc{}
		} 

		res[txn.k] = append(prev, txn.v)
	}

	return res
}