package internal

// struct with k as column name and v as the array function to apply to its values
type Transformation struct{k string; v RowWiseFunc}

// function that transforms each element into another element
type RowWiseFunc func(interface{}) interface{}

// Merges a slice of transformations into a map of lists of RowWiseFunc functions
func MergeTransformations(txns []Transformation) map[string][]RowWiseFunc {
	res := map[string][]RowWiseFunc{}	

	for _, txn 	:= range txns {
		prev, ok := res[txn.k]
		if !ok {
			prev = []RowWiseFunc{}
		} 

		res[txn.k] = append(prev, txn.v)
	}

	return res
}