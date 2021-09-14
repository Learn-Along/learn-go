package types

import "regexp"

type Datatype int

type arrayFunc func([]interface{}) []interface{}

type Column struct {
	Name string
	items []*item
	Dtype Datatype
	keys []string
}

type colTransform func() Column

type item struct {
	pk string
	value interface{}
}

const (
	IntType Datatype = iota
	FloatType
	StringType
	ObjectType
	BooleanType
	ArrayType
)

// Returns a filter function that gets only values greater than the operand
// The operand can reference a constant, or a Col
func (c *Column) GreaterThan(operand interface{}) Filter {
	return func() []string {return nil}
}

// Returns a filter function that gets only values greater than or equal to the operand
// The operand can reference a constant, or a Col
func (c *Column) GreaterOrEquals(operand interface{}) Filter {
	return func() []string {return nil}
}

// Returns a filter function that gets only values less than the operand
// The operand can reference a constant, or a Col
func (c *Column) LessThan(operand interface{}) Filter {
	return func() []string {return nil}
}

// Returns a filter function that gets only values less than or equal to the operand
// The operand can reference a constant, or a Col
func (c *Column) LessOrEquals(operand interface{}) Filter {
	return func() []string {return nil}
}

// Returns a filter function that gets only values equal to the operand
// The operand can reference a constant, or a Col
func (c *Column) Equals(operand interface{}) Filter {
	return func() []string {return nil}
}

// Returns a filter function that gets only values that are like the regex expression
func (c *Column) IsLike(pattern *regexp.Regexp) Filter  {
	return func() []string {return nil}
}

// Returns a transformer method to transform the column from one form to another
// It is passed a function expecting an array of values of any type
func (c *Column) Tx(op arrayFunc) colTransform {
	return func() Column {return Column{}}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *Column) Order(option SortOrder) sortOption {
	return sortOption{Col: c, Order: option}
}

// Returns a list a function that extracts the list of values of the given items
func (c *Column) Items() []interface{} {
	return nil
}
