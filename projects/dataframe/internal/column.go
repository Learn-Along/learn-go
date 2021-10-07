package internal

import (
	"regexp"
)

const (
	IntType Datatype = iota
	FloatType
	StringType
	BoolType
)


type Datatype int

/*
For now let us support the following types
	string 
	bool 
	int 
	float64
*/
type Item interface {}

/*
For now let us support the following types
	[]string 
	[]bool 
	[]int 
	[]float64
*/
type ItemSlice interface {}

/*
This can be a literal value of type:
	string 
	int
	float64
Or it can be a Column
*/
type LiteralOrColumn interface{}

/*
The base column interface for all column types
*/
type Column interface{
	Items() ItemSlice
	Insert(index int, value Item)
	DeleteMany(indices []int)
	GreaterThan(operand LiteralOrColumn) FilterType
	GreaterOrEquals(operand LiteralOrColumn) FilterType
	LessThan(operand LiteralOrColumn) FilterType
	LessOrEquals(operand LiteralOrColumn) FilterType
	Equals(operand LiteralOrColumn) FilterType
	IsLike(pattern *regexp.Regexp) FilterType
	Tx(op RowWiseFunc) Transformation
	Agg(aggFunc AggregateFunc) Aggregation
	Order(option SortOrder) SortOption
	Len() int 
	Name() string 
	ItemAt(index int) Item
	GetDatatype() Datatype
	Defragmentize(newOrder []int)
}
