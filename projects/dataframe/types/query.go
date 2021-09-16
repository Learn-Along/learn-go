package types

type Query struct{
	ops []colTransform
}

type SortOrder int

type sortOption struct {
	Col *Column
	Order SortOrder
}

const (
	ASC SortOrder = iota
	DESC
)

// Actually executes the query
func (q *Query) Execute() ([]map[string]interface{}, error) {
	newDf := Dataframe{}
	for _, op := range q.ops {
		newCol := op()
		newDf.cols[newCol.Name] = &newCol
	}

	return newDf.ToArray()
}

// Given a list of boolean corresponding to indices of the items,
// true meaning the item should be included, false meaning that item should be excluded
// the method then returns a query instance
func (q *Query) Where(filter []bool) *Query {
	return q
}

// Sorts the data by the col provided in the sort option, and int he order given
func (q *Query) SortBy(options ...sortOption) *Query {
	return q
}

// Groups the data into gorups that have same values for the given columns
func (q *Query) GroupBy(cols ...*Column) *Query {
	return q
}

// Applies the col transforms to the query
func (q *Query) Apply(ops ...colTransform) *Query {
	return nil
}


// Logic combinations

// Combines a list of filters to produce a combined AND logical filter
func AND(filters ...[]bool) []bool {
	combinedFilters := []bool{}
	filterLengths := []int{}

	maxLength := 0
	for _, filter := range filters {
		filterLength := len(filter)
		filterLengths = append(filterLengths, filterLength)

		if maxLength < filterLength {
			maxLength = filterLength
		}
	}

	for row := 0; row < maxLength; row++ {
		value := true

		for i, filter := range filters {
			if row < filterLengths[i] {
				value = value && filter[row]				
			} else {
				value = false
				break
			}
		}

		combinedFilters = append(combinedFilters, value)
	}

	return combinedFilters
}

// Combines a list of filters to produce a combined OR logical filter
func OR(filters ...[]bool) []bool {
	combinedFilters := []bool{}
	filterLengths := []int{}

	maxLength := 0
	for _, filter := range filters {
		filterLength := len(filter)
		filterLengths = append(filterLengths, filterLength)

		if maxLength < filterLength {
			maxLength = filterLength
		}
	}

	for row := 0; row < maxLength; row++ {
		value := false

		for i, filter := range filters {
			if row < filterLengths[i] {
				value = value || filter[row]				
			}
		}

		combinedFilters = append(combinedFilters, value)
	}

	return combinedFilters
}

// Inverts a given filter to produce a NOT logical filter
func NOT(filter []bool) []bool {
	combinedFilters := []bool{}
	for _, value := range filter {
		combinedFilters = append(combinedFilters, !value)
	}

	return combinedFilters
}
