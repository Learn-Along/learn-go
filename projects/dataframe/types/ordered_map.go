package types

import "sort"

type orderedMapType map[int]interface{}

// Converts an ordered map to a slice
func (o *orderedMapType) ToSlice() []interface{} {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]interface{}, count)

	counter := 0
	for i := range *o {
		indices[counter] = i
		counter++
	}

	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})

	for i, index := range indices {
		slice[i] = (*o)[index]
	}

	return slice
}

// Reorders the orderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *orderedMapType) Defragmentize(newOrder []int) {
	copyOfO := orderedMapType{}
	for k, v := range *o {
		copyOfO[k] = v
		delete(*o, k)
	}	

	for newRow, oldRow := range newOrder {
		(*o)[newRow] = copyOfO[oldRow]
	}
}