package internal

import (
	"fmt"
	"regexp"
)

type StringColumn struct {
	Title string
	Values OrderedStringMapType
}

// returns the Title of the column
func (c *StringColumn) Name() string {
	return c.Title
}

// Number of Values in int column
func (c *StringColumn) Len() int {
	return c.Values.Len()
}

// Number of Values in int column
func (c *StringColumn) ItemAt(index int) Item {
	return c.Values[index]
}

// Returns a list of Items
func (c *StringColumn) Items() ItemSlice {
	return c.Values.ToSlice()
}

// Returns the data type of the given column
func (c *StringColumn) GetDatatype() Datatype {
	return StringType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *StringColumn) Defragmentize(newOrder []int) {
	c.Values.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// it ignores the Insert if the value is not a string
func (c *StringColumn) Insert(index int, value Item) {
	nextIndex := c.Values.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.Values[i] = ""		
		}
	}

	switch v := value.(type) {
	case string:
		c.Values[index] = v
	}
}

// Deletes many indices at once
func (c *StringColumn) DeleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.Values, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) GreaterThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.Values.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v > op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = v > operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) GreaterOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.Values.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v >= op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = v >= operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) LessThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.Values.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v < op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = v < operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) LessOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsString string 
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.Values.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v <= op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = v <= operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) Equals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operands []string

	switch v := operand.(type) {
	case string:
		for i, v := range c.Values {
			flags[i] = v == operand
		}	
		return flags

	case StringColumn:
		operands = v.Values.ToSlice().([]string)
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v == op
			}
		}
		return flags

	default:
		return flags
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *StringColumn) IsLike(pattern *regexp.Regexp) FilterType  {
	count := len(c.Values)
	flags := make(FilterType, count)

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = pattern.MatchString(fmt.Sprintf("%v", v))
		}
	}

	return flags
}

// Returns transformer method specific to this column to transform its values from one thing to another
// It is passed a function expecting a value any type
func (c *StringColumn) Tx(op RowWiseFunc) Transformation {
	return Transformation{k: c.Title, v: op}
}

// Returns an Aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *StringColumn) Agg(aggFunc AggregateFunc) Aggregation {
	return Aggregation{c.Title: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *StringColumn) Order(option SortOrder) SortOption {
	return SortOption{c.Title: option}
}
