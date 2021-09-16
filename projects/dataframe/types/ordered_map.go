package types

type OrderedMap map[int]interface{}

// Converts an ordered map to a slice
func (o OrderedMap) ToSlice() []interface{} {
	count := len(o)
	slice := make([]interface{}, count)

	for i := 0; i < count; i++ {
		slice[i] = o[i]
	}

	return slice
}