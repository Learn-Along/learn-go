package types

type Query struct{
	ops []colTransform
}

type Filter func() []string

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

// Narrows down the filter to a given filter and returns a query instance
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

// Combines a list of filters to produce a combined AND logical filter
func AND(filters ...Filter) Filter {
	return func() []string {return nil}
}

// Combines a list of filters to produce a combined OR logical filter
func OR(filters ...Filter) Filter {
	return func() []string {return nil}
}

// Inverts a given filter to produce a NOT logical filter
func NOT(filter Filter) Filter {
	return func() []string {return nil}
}
