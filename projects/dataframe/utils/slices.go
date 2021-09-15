package utils

import "fmt"

// Merges two slices into one
func MergeSlices(first []interface{}, second []interface{}) []interface{} {
	return nil
}

// Checks to see if slices are equal to each other, return error if they are not
func SliceEquals(first []interface{}, second []interface{}) error {
	length := len(first)

	if length != len(second) {
		return fmt.Errorf("length mismatch")
	}

	for i := 0; i < length; i++ {
		if(first[i] != second[i]){
			return fmt.Errorf("index %d values don't match", i)
		}
	}

	return nil
}