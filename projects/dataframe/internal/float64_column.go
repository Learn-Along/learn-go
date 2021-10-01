package internal

import (
	"fmt"
	"regexp"
)

type Float64Column struct {
	name string
	items OrderedFloat64MapType
}

// returns the name of the column
func (c *Float64Column) Name() string {
	return c.name
}

// Number of items in int column
func (c *Float64Column) Len() int {
	return c.items.Len()
}

// Number of items in int column
func (c *Float64Column) ItemAt(index int) Item {
	return c.items[index]
}

// Returns a list of Items
func (c *Float64Column) Items() ItemSlice {
	return c.items.ToSlice()
}

// Returns the data type of the given column
func (c *Float64Column) GetDatatype() Datatype {
	return FloatType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *Float64Column) Defragmentize(newOrder []int) {
	c.items.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// it ignores the insert if the value is not an int or float64
func (c *Float64Column) insert(index int, value Item) {
	nextIndex := c.items.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.items[i] = 0		
		}
	}

	switch v := value.(type) {
	case int:
		c.items[index] = float64(v)
	case float64:
		c.items[index] = v
	}
}

// Deletes many indices at once
func (c *Float64Column) deleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.items, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) GreaterThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.items.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.items.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v > operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) GreaterOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.items.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.items.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v >= operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) LessThan(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.items.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.items.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v < operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) LessOrEquals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operandAsFloat float64 
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.items.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.items.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v <= operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) Equals(operand LiteralOrColumn) filterType {
	count := len(c.items)
	flags := make(filterType, count)
	var operands []float64
	var operandAsFloat float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.items.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.items.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
		flags[i] = v == operandAsFloat
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *Float64Column) IsLike(pattern *regexp.Regexp) filterType  {
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
func (c *Float64Column) Tx(op rowWiseFunc) transformation {
	return transformation{k: c.name, v: op}
}

// Returns an aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *Float64Column) Agg(aggFunc aggregateFunc) aggregation {
	return aggregation{c.name: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *Float64Column) Order(option sortOrder) sortOption {
	return sortOption{c.name: option}
}
