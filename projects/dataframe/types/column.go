package types

import "regexp"

type Datatype int

type arrayFunc func([]interface{}) []interface{}

type Column struct {
	Name string
	Items []interface{}
	Dtype Datatype
}

type Filter map[string][]bool

type colTransform func() Column

const (
	IntType Datatype = iota
	FloatType
	StringType
	ObjectType
	BooleanType
	ArrayType
)

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
func (c *Column) insert(index int, value interface{}) {
	nextIndex := len(c.Items)

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.Items = append(c.Items, nil)		
		}
	}

	c.Items[index] = value
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *Column) GreaterThan(operand interface{}) Filter {
	return nil
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Column) GreaterOrEquals(operand interface{}) Filter {
	return nil
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *Column) LessThan(operand interface{}) Filter {
	return nil
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Column) LessOrEquals(operand interface{}) Filter {
	return nil
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *Column) Equals(operand interface{}) Filter {
	return nil
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *Column) IsLike(pattern *regexp.Regexp) Filter  {
	return nil
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
