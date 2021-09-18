package types

const (
	ASC sortOrder = iota
	DESC
)

const (
	// This order is important. 
	// filter first, 
	// then group, 
	// then sort each group,
	// then apply whatever,
	// then select the field
	FILTER_ACTION actionType = iota
	GROUPBY_ACTION
	SORT_ACTION
	APPLY_ACTION
	SELECT_ACTION
)

type action struct {
	_type actionType
	payload interface{}
}

type actionType int

type sortOrder int

type sortOption map[string]sortOrder

type filterType []bool

/**
* query
*/
type query struct{
	ops []action
	df *Dataframe
}

// Actually executes the query
func (q *query) Execute() ([]map[string]interface{}, error) {
	// may need to add a recover defer
	// filters := []filterType{}
	// aggs := []aggregation{}
	// sortOptions := []sortOption{}
	// txList := []transformation{}

	// // combine similar actions together
	// for _, act := range q.ops {
	// 	switch act._type {
	// 	case FILTER_ACTION:
	// 		filters = append(filters, act.payload.(filterType))
	// 	case GROUPBY_ACTION:
	// 		aggs = append(aggs, act.payload.([]aggregation)...)
	// 	case SORT_ACTION:
	// 		sortOptions = append(sortOptions, act.payload.([]sortOption)...)
	// 	case APPLY_ACTION:
	// 		txList = append(txList, act.payload.([]transformation)...)
	// 	}
	// }

	// filteredDf, err := q.df.filter(AND(filters...))
	// if err != nil {
	// 	return nil, err
	// }

	// groupedDfs, err := filteredDf.groupby(mergeAggregations(aggs))
	// mergedTxs := mergeTransformations(txList)

	// for _, df := range groupedDfs {
	// 	// sort the df
	// 	// apply the app functions on df
	// }

	// merge the groupedDfs

	return nil, nil
}

// Given a list of boolean corresponding to indices of the items,
// true meaning the item should be included, false meaning that item should be excluded
// the method then returns a query instance
func (q *query) Where(filter filterType) *query {
	q.ops = append(q.ops, action{_type: FILTER_ACTION, payload: filter})
	return q
}

// Sorts the data by the col provided in the sort option, and int he order given
func (q *query) SortBy(options ...sortOption) *query {
	q.ops = append(q.ops, action{_type: SORT_ACTION, payload: options})
	return q
}

// Groups the data into groups that have same values for the given columns
func (q *query) GroupBy(aggs ...aggregation) *query {
	q.ops = append(q.ops, action{_type: GROUPBY_ACTION, payload: aggs})
	return q
}

// Applies the col transforms to the query
func (q *query) Apply(ops ...transformation) *query {
	q.ops = append(q.ops, action{_type: APPLY_ACTION, payload: ops})
	return nil
}


// Logic combinations

// Combines a list of maps of filters to produce a combined AND logical filter
func AND(filters ...filterType) filterType{
	combinedFilter := filterType{}

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
func OR(filters ...filterType) filterType {
	combinedFilter := filterType{}

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
func NOT(filter filterType) filterType {
	count := len(filter)
	combinedFilter := make(filterType, count)

	for i, value := range filter {
		combinedFilter[i] = !value
	}

	return combinedFilter
}
