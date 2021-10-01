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
	insert(index int, value Item)
	deleteMany(indices []int)
	GreaterThan(operand LiteralOrColumn) filterType
	GreaterOrEquals(operand LiteralOrColumn) filterType
	LessThan(operand LiteralOrColumn) filterType
	LessOrEquals(operand LiteralOrColumn) filterType
	Equals(operand LiteralOrColumn) filterType
	IsLike(pattern *regexp.Regexp) filterType
	Tx(op rowWiseFunc) transformation
	Agg(aggFunc aggregateFunc) aggregation
	Order(option sortOrder) sortOption
	Len() int 
	Name() string 
	ItemAt(index int) Item
	GetDatatype() Datatype
	Defragmentize(newOrder []int)
}
