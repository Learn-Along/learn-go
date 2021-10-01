package internal

import "sort"

/*
* Supported types are:
* - map[int]string
* - map[int]int
* - map[int]float64
* - map[int]bool
 */
type OrderedMapType interface{
	ToSlice() ItemSlice
	Defragmentize(newOrder []int)
	Len() int
}
// type OrderedMapType map[int]interface{}

type OrderedStringMapType map[int]string

// Gets the length of the OrderedStringMapType
func (o *OrderedStringMapType) Len() int {
	return len(*o)
}

// Converts an ordered map to a slice
func (o *OrderedStringMapType) ToSlice() ItemSlice {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]string, count)

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

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *OrderedStringMapType) Defragmentize(newOrder []int) {
	noOfItems := len(newOrder)
	cleanedMap := make(OrderedStringMapType, noOfItems)

	for newRow, oldRow := range newOrder {
		cleanedMap[newRow] = (*o)[oldRow]
	}

	for key := range (*o) {
		delete(*o, key)
	}

	*o = cleanedMap
}

type OrderedIntMapType map[int]int

// Gets the length of the OrderedStringMapType
func (o *OrderedIntMapType) Len() int {
	return len(*o)
}

// Converts an ordered map to a slice
func (o *OrderedIntMapType) ToSlice() ItemSlice {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]int, count)

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

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *OrderedIntMapType) Defragmentize(newOrder []int) {
	noOfItems := len(newOrder)
	cleanedMap := make(OrderedIntMapType, noOfItems)

	for newRow, oldRow := range newOrder {
		cleanedMap[newRow] = (*o)[oldRow]
	}

	*o = cleanedMap
}

type OrderedFloat64MapType map[int]float64

// Gets the length of the OrderedStringMapType
func (o *OrderedFloat64MapType) Len() int {
	return len(*o)
}

// Converts an ordered map to a slice
func (o *OrderedFloat64MapType) ToSlice() ItemSlice {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]float64, count)

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

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *OrderedFloat64MapType) Defragmentize(newOrder []int) {
	noOfItems := len(newOrder)
	cleanedMap := make(OrderedFloat64MapType, noOfItems)

	for newRow, oldRow := range newOrder {
		cleanedMap[newRow] = (*o)[oldRow]
	}

	*o = cleanedMap
}

type OrderedBoolMapType map[int]bool

// Gets the length of the OrderedStringMapType
func (o *OrderedBoolMapType) Len() int {
	return len(*o)
}

// Converts an ordered map to a slice
func (o *OrderedBoolMapType) ToSlice() ItemSlice {
	count := len(*o)
	indices := make([]int, count)
	slice := make([]bool, count)

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

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (o *OrderedBoolMapType) Defragmentize(newOrder []int) {
	noOfItems := len(newOrder)
	cleanedMap := make(OrderedBoolMapType, noOfItems)

	for newRow, oldRow := range newOrder {
		cleanedMap[newRow] = (*o)[oldRow]
	}

	*o = cleanedMap
}