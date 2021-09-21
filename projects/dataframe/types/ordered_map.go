package types

import "sort"

type orderedMapType map[int]interface{}

// Converts an ordered map to a slice
func (o *orderedMapType) ToSlice() []interface{} {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]interface{}, count)

	// FIXME: if the assumption is there are no gaps in this orderedMap,
	// which assumption is true, as Deframentize needs to be called everytime
	// there is random update to the map, then concurrency here is possible. try a goroutine and a channel
	// and two ranges, one, on *o, the other on the channel
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
		// FIXME: concurrency is possible
		copyOfO[k] = v
		delete(*o, k)
	}	

	for newRow, oldRow := range newOrder {
		// FIXME: concurrency is possible
		(*o)[newRow] = copyOfO[oldRow]
	}
}