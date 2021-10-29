package dataframe

import "github.com/learn-along/learn-go/projects/dataframe/internal"

const (
	ASC internal.SortOrder = iota
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

var (
	MAX internal.AggregateFunc = internal.GetMax
	MIN internal.AggregateFunc = internal.GetMin
	SUM internal.AggregateFunc = internal.GetSum
	MEAN internal.AggregateFunc = internal.GetMean
	COUNT internal.AggregateFunc = internal.GetCount
	RANGE internal.AggregateFunc = internal.GetRange
	// PERCENTILE(int) etc.
)

type action struct {
	_type actionType
	payload interface{}
}

type actionType int

/*
* GroupBy Options 
*/
type groupByOption struct {
	fields []string
	aggs []internal.Aggregation
	q *query
}

// aggregates the different groups
func (g *groupByOption) Agg(aggs...internal.Aggregation) *query {
	g.aggs = append(g.aggs, aggs...)
	g.q.ops = append(g.q.ops, action{_type: GROUPBY_ACTION, payload: g})
	return g.q
}

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
	var gopt *groupByOption
	filters := []internal.FilterType{}
	sortOptions := []internal.SortOption{}
	txList := []internal.Transformation{}
	selectedFields := []string{}

	// combine similar actions together
	for _, act := range q.ops {
		switch act._type {
		case FILTER_ACTION:
			filters = append(filters, act.payload.(internal.FilterType))
		case GROUPBY_ACTION:
			gopt = act.payload.(*groupByOption)
		case SORT_ACTION:
			sortOptions = append(sortOptions, act.payload.([]internal.SortOption)...)
		case APPLY_ACTION:
			txList = append(txList, act.payload.([]internal.Transformation)...)
		case SELECT_ACTION:
			selectedFields = append(selectedFields, act.payload.([]string)...)
		}
	}

	df, err := q.df.getFilteredDf(AND(filters...))
	if err != nil {
		return nil, err
	}

	if gopt != nil {
		df, err = df.getGroupedDf(gopt)
		if err != nil {
			return nil, err
		}
	}

	if len(sortOptions) > 0 {
		df, err = df.getSortedDf(sortOptions...)
		if err != nil {
			return nil, err
		}
	}

	if len(txList) > 0 {
		mergedTxs := internal.MergeTransformations(txList)	
		err = df.apply(mergedTxs)
		if err != nil {
			return nil, err
		}		
	}

	// FIXME: Could we return a dataframe instead of converting this to an array first
	// This could mean that less columns are even generated and passed out
	return df.ToArray(selectedFields...)	
}

// Given a list of boolean corresponding to indices of the items,
// true meaning the item should be included, false meaning that item should be excluded
// the method then returns a query instance
func (q *query) Where(filter internal.FilterType) *query {
	q.ops = append(q.ops, action{_type: FILTER_ACTION, payload: filter})
	return q
}

// Sorts the data by the col provided in the sort option, and int he order given
func (q *query) SortBy(options ...internal.SortOption) *query {
	q.ops = append(q.ops, action{_type: SORT_ACTION, payload: options})
	return q
}

// Groups the data into groups that have same values for the given columns/fields
func (q *query) GroupBy(fields ...string) *groupByOption {
	return &groupByOption{q: q, fields: fields, aggs: []internal.Aggregation{}}
}

// Applies the col transforms to the query
func (q *query) Apply(ops ...internal.Transformation) *query {
	q.ops = append(q.ops, action{_type: APPLY_ACTION, payload: ops})
	return q
}


// Logic combinations

// Combines a list of filters to produce a combined AND logical filter
func AND(filters ...internal.FilterType) internal.FilterType{
	// FIXME: What if I got the maximum length of the filters,
	// used make to create a FilterType of that length (values default to false),
	// and then got rid of append at the bottom
	combinedFilter := internal.FilterType{}

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
func OR(filters ...internal.FilterType) internal.FilterType {
	// FIXME: What if I got the maximum length of the filters, (use a utility from slices.go)
	// used make to create a FilterType of that length (values default to false),
	// and then got rid of append at the bottom
	combinedFilter := internal.FilterType{}

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
func NOT(filter internal.FilterType) internal.FilterType {
	count := len(filter)
	combinedFilter := make(internal.FilterType, count)

	for i, value := range filter {
		// FIXME: concurrency possible
		combinedFilter[i] = !value
	}

	return combinedFilter
}
