package types

import (
	"fmt"
	"regexp"
)

type IntColumn struct {
	name string
	items OrderedIntMapType
}

// returns the name of the column
func (c *IntColumn) Name() string {
	return c.name
}

// Number of items in int column
func (c *IntColumn) Len() int {
	return c.items.Len()
}

// Number of items in int column
func (c *IntColumn) ItemAt(index int) Item {
	return c.items[index]
}

// Returns a list of Items
func (c *IntColumn) Items() ItemSlice {
	return c.items.ToSlice()
}

// Returns the data type of the given column
func (c *IntColumn) GetDatatype() Datatype {
	return IntType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *IntColumn) Defragmentize(newOrder []int) {
	c.items.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// it ignores the insert if the value is not a number of int or float64 types
func (c *IntColumn) insert(index int, value Item) {
	nextIndex := c.items.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.items[i] = 0		
		}
	}

	switch v := value.(type) {
	case int:
		c.items[index] = v
	case float64:
		c.items[index] = int(v)
	}
}

// Deletes many indices at once
func (c *IntColumn) deleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.items, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) GreaterThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.items.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.items.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v > operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) GreaterOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.items.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.items.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v >= operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) LessThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.items.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.items.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v < operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) LessOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsInt int 
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.items.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.items.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v <= operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) Equals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operands []int
	var operandAsInt int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.items.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.items.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
		flags[i] = v == operandAsInt
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *IntColumn) IsLike(pattern *regexp.Regexp) filterType  {
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
func (c *IntColumn) Tx(op rowWiseFunc) transformation {
	return transformation{k: c.name, v: op}
}

// Returns an aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *IntColumn) Agg(aggFunc aggregateFunc) aggregation {
	return aggregation{c.name: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *IntColumn) Order(option sortOrder) sortOption {
	return sortOption{c.name: option}
}
