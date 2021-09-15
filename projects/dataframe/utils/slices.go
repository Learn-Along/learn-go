package utils

import "sort"

type SortOrder int 

const (
	ASC SortOrder = iota
	DESC
)

// Extracts a given field from n list of maps, returning a list f the values for that field,
// where the field does not exist, a nil is added
func ExtractFieldFromMapList(data []map[string]interface{}, field string) []interface{} {
	extracted := []interface{}{}

	for _, record := range data {
		if value, ok := record[field]; ok {
			extracted = append(extracted, value)
		} else {
			extracted = append(extracted, nil)
		}
	}

	return extracted
}

// Checks to see if the string slices are equal to each other, returning false or true
func AreStringSliceEqual(first []string, second []string) bool {
	length := len(first)

	if length != len(second) {
		return false
	}

	for i := 0; i < length; i++ {
		if(first[i] != second[i]){
			return false
		}
	}

	return true
}

// Checks to see if slices are equal to each other, returning false or true
func AreSliceEqual(first []interface{}, second []interface{}) bool {
	length := len(first)

	if length != len(second) {
		return false
	}

	for i := 0; i < length; i++ {
		if(first[i] != second[i]){
			return false
		}
	}

	return true
}

// Sorts a given string slice in a given order
func SortStringSlice(slice []string, order SortOrder) []string {
	copyOfSlice := append([]string{}, slice...)

	sortFunc := func(i, j int) bool {
		return copyOfSlice[i] < copyOfSlice[j]
	}

	if order == DESC {
		sortFunc = func(i, j int) bool {
			return copyOfSlice[i] > copyOfSlice[j]
		}	
	}

	sort.Slice(copyOfSlice, sortFunc)

	return copyOfSlice
}