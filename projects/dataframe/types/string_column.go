package types

import (
	"fmt"
	"regexp"
)

type StringColumn struct {
	name string
	items OrderedStringMapType
}

// returns the name of the column
func (c *StringColumn) Name() string {
	return c.name
}

// Number of items in int column
func (c *StringColumn) Len() int {
	return c.items.Len()
}

// Number of items in int column
func (c *StringColumn) ItemAt(index int) Item {
	return c.items[index]
}

// Returns a list of Items
func (c *StringColumn) Items() ItemSlice {
	return c.items.ToSlice()
}

// Returns the data type of the given column
func (c *StringColumn) GetDatatype() Datatype {
	return StringType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *StringColumn) Defragmentize(newOrder []int) {
	c.items.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
func (c *StringColumn) insert(index int, value Item) {
	nextIndex := c.items.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.items[i] = ""		
		}
	}

	c.items[index] = value.(string)
}

// Deletes many indices at once
func (c *StringColumn) deleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.items, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) GreaterThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.items.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.items[i]; ok {
				flags[i] = v > op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.items[i]; ok {
			flags[i] = v > operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) GreaterOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.items.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.items[i]; ok {
				flags[i] = v >= op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.items[i]; ok {
			flags[i] = v >= operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) LessThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsString string
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.items.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.items[i]; ok {
				flags[i] = v < op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.items[i]; ok {
			flags[i] = v < operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) LessOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsString string 
	var operands []string

	switch v := operand.(type) {
	case string:
		operandAsString = v
	case StringColumn:
		operands = v.items.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.items[i]; ok {
				flags[i] = v <= op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.items[i]; ok {
			flags[i] = v <= operandAsString
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *StringColumn) Equals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operands []string

	switch v := operand.(type) {
	case StringColumn:
		operands = v.items.ToSlice().([]string)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.items[i]; ok {
				flags[i] = v == op
			}
		}

		return flags
	}

	for i, v := range c.items {
		flags[i] = v == operand
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *StringColumn) IsLike(pattern *regexp.Regexp) filterType  {
	count := len(c.items)
	flags := make(filterType, count)

	for i := 0; i < count; i++ {
		if v, ok := c.items[i]; ok {
			flags[i] = pattern.MatchString(fmt.Sprintf("%v", v))
		}
	}

	return flags
}

// Returns transformer method specific to this column to transform its values from one thing to another
// It is passed a function expecting a value any type
func (c *StringColumn) Tx(op rowWiseFunc) transformation {
	return transformation{k: c.name, v: op}
}

// Returns an aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *StringColumn) Agg(aggFunc aggregateFunc) aggregation {
	return aggregation{c.name: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *StringColumn) Order(option sortOrder) sortOption {
	return sortOption{c.name: option}
}
