package internal

import (
	"fmt"
	"regexp"
)

type IntColumn struct {
	Title string
	Values OrderedIntMapType
}

// returns the Title of the column
func (c *IntColumn) Name() string {
	return c.Title
}

// Number of Values in int column
func (c *IntColumn) Len() int {
	return c.Values.Len()
}

// Number of Values in int column
func (c *IntColumn) ItemAt(index int) Item {
	return c.Values[index]
}

// Returns a list of Items
func (c *IntColumn) Items() ItemSlice {
	return c.Values.ToSlice()
}

// Returns the data type of the given column
func (c *IntColumn) GetDatatype() Datatype {
	return IntType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *IntColumn) Defragmentize(newOrder []int) {
	c.Values.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// it ignores the Insert if the value is not a number of int or float64 types
func (c *IntColumn) Insert(index int, value Item) {
	nextIndex := c.Values.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.Values[i] = 0		
		}
	}

	switch v := value.(type) {
	case int:
		c.Values[index] = v
	case float64:
		c.Values[index] = int(v)
	}
}

// Deletes many indices at once
func (c *IntColumn) DeleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.Values, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) GreaterThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.Values.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.Values.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v > operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) GreaterOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.Values.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.Values.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v >= operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) LessThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsInt int
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.Values.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.Values.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v < operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) LessOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsInt int 
	var operands []int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.Values.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.Values.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
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
			flags[i] = v <= operandAsInt
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *IntColumn) Equals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operands []int
	var operandAsInt int

	switch v := operand.(type) {
	case int:
		operandAsInt = v
	case float64:
		operandAsInt = int(v)
	case IntColumn:
		operands = v.Values.ToSlice().([]int)
	case Float64Column:
		operands = make([]int, 0, count)
		for _, v := range v.Values.ToSlice().([]float64) {
			operands = append(operands, int(v))
		} 
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = v == op
			}
		}

		return flags
	}

	for i, v := range c.Values {
		flags[i] = v == operandAsInt
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *IntColumn) IsLike(pattern *regexp.Regexp) FilterType  {
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
func (c *IntColumn) Tx(op RowWiseFunc) Transformation {
	return Transformation{k: c.Title, v: op}
}

// Returns an Aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *IntColumn) Agg(aggFunc AggregateFunc) Aggregation {
	return Aggregation{c.Title: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *IntColumn) Order(option SortOrder) SortOption {
	return SortOption{c.Title: option}
}
