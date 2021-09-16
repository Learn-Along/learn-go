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
	combinedFilters := Filter{}

	for _, filter := range filters {

		for field, newArray := range filter {

			oldArray, ok := combinedFilters[field]
			if !ok {
				combinedFilters[field] = newArray
				continue
			}
			
			oldArrayLength := len(oldArray)
			newArrayLength := len(newArray)

			for row, value := range oldArray {	
				if row < newArrayLength {
					combinedFilters[field][row] = value && newArray[row]
				} else {
					combinedFilters[field][row] = false
				}
			}

			// fill up any new rows that didn't exist originally, with false
			for row := oldArrayLength; row < newArrayLength; row++ {
				combinedFilters[field] = append(combinedFilters[field], false)
			}
		}
	}

	return combinedFilters
}

// Combines a list of filters to produce a combined OR logical filter
func OR(filters ...Filter) Filter {
	combinedFilters := Filter{}

	for _, filter := range filters {

		for field, newArray := range filter {

			oldArray, ok := combinedFilters[field]
			if !ok {
				combinedFilters[field] = newArray
				continue
			}
			
			oldArrayLength := len(oldArray)
			newArrayLength := len(newArray)

			for row, value := range oldArray {	
				if row < newArrayLength {
					combinedFilters[field][row] = value || newArray[row]
				}
			}

			// fill up any new rows that didn't exist originally, with the new value
			for row := oldArrayLength; row < newArrayLength; row++ {
				combinedFilters[field] = append(combinedFilters[field], newArray[row])
			}
		}
	}

	return combinedFilters
}

// Inverts a given filter to produce a NOT logical filter
func NOT(filter Filter) Filter {
	combinedFilters := Filter{}

	for field, data := range filter {
		combinedFilters[field] = []bool{}

		for _, value := range data {
			combinedFilters[field] = append(combinedFilters[field], !value)
		}
	}

	return combinedFilters
}
