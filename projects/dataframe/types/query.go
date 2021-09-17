package types

const (
	ASC SortOrder = iota
	DESC
)

type SortOrder int

type sortOption struct {
	Col *Column
	Order SortOrder
}

/**
* Query
*/
type Query struct{
	ops []colTransform
}

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
func (q *Query) Where(filter Filter) *Query {
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

// Combines a list of maps of filters to produce a combined AND logical filter
func AND(filters ...Filter) Filter{
	combinedFilter := Filter{}

	for _, filter := range filters {
		currentLength := len(combinedFilter)
		newArrayLength := len(filter)
		
		if currentLength == 0 {
			combinedFilter = filter
			continue
		}

		for row, value := range combinedFilter {	
			if row < newArrayLength {
				combinedFilter[row] = value && filter[row]
			} else {
				combinedFilter[row] = false
			}
		}

		// fill up any new rows that didn't exist originally, with false
		for row := currentLength; row < newArrayLength; row++ {
			combinedFilter = append(combinedFilter, false)
		}
	}

	return combinedFilter
}

// Combines a list of filters to produce a combined OR logical filter
func OR(filters ...Filter) Filter {
	combinedFilter := Filter{}

	for _, filter := range filters {
		currentLength := len(combinedFilter)
		newArrayLength := len(filter)
			
		if currentLength == 0 {
			combinedFilter = filter
			continue
		}			

		for row, value := range combinedFilter {	
			if row < newArrayLength {
				combinedFilter[row] = value || filter[row]
			}
		}

		// fill up any new rows that didn't exist originally, with the new value
		for row := currentLength; row < newArrayLength; row++ {
			combinedFilter = append(combinedFilter, filter[row])
		}

	}

	return combinedFilter
}

// Inverts a given filter to produce a NOT logical filter
func NOT(filter Filter) Filter {
	count := len(filter)
	combinedFilter := make(Filter, count)

	for i, value := range filter {
		combinedFilter[i] = !value
	}

	return combinedFilter
}
