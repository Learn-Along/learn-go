package types

type OrderedMap map[int]interface{}

// Converts an ordered map to a slice
func (o *OrderedMap) ToSlice() []interface{} {
	count := len(*o)
	slice := make([]interface{}, count)

	for i := 0; i < count; i++ {
		slice[i] = (*o)[i]
	}

	return slice
}

// Reorders the OrderedMap ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *OrderedMap) Defragmentize(newOrder []int) {
	copyOfO := OrderedMap{}
	for k, v := range *o {
		copyOfO[k] = v
		delete(*o, k)
	}	

	for newRow, oldRow := range newOrder {
		(*o)[newRow] = copyOfO[oldRow]
	}
}