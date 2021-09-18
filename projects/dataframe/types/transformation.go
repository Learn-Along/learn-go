package types

// map of column name and the array function to apply to its values
type transformation map[string]rowWiseFunc

// function that transforms each element into another element
type rowWiseFunc func(interface{}) interface{}