package internal

import (
	"fmt"
	"regexp"
)

type Float64Column struct {
	Title string
	Values OrderedFloat64MapType
}

// returns the Title of the column
func (c *Float64Column) Name() string {
	return c.Title
}

// Number of Values in int column
func (c *Float64Column) Len() int {
	return c.Values.Len()
}

// Number of Values in int column
func (c *Float64Column) ItemAt(index int) Item {
	return c.Values[index]
}

// Returns a list of Items
func (c *Float64Column) Items() ItemSlice {
	return c.Values.ToSlice()
}

// Returns the data type of the given column
func (c *Float64Column) GetDatatype() Datatype {
	return FloatType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *Float64Column) Defragmentize(newOrder []int) {
	c.Values.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// it ignores the Insert if the value is not an int or float64
func (c *Float64Column) Insert(index int, value Item) {
	nextIndex := c.Values.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.Values[i] = 0		
		}
	}

	switch v := value.(type) {
	case int:
		c.Values[index] = float64(v)
	case float64:
		c.Values[index] = v
	}
}

// Deletes many indices at once
func (c *Float64Column) DeleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.Values, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) GreaterThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.Values.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.Values.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v > operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) GreaterOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.Values.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.Values.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v >= operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) LessThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsFloat float64
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.Values.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.Values.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v < operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) LessOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsFloat float64 
	var operands []float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.Values.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.Values.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
			flags[i] = v <= operandAsFloat
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *Float64Column) Equals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operands []float64
	var operandAsFloat float64

	switch v := operand.(type) {
	case int:
		operandAsFloat = float64(v)
	case float64:
		operandAsFloat = v
	case Float64Column:
		operands = v.Values.ToSlice().([]float64)
	case IntColumn:
		operands = make([]float64, 0, count)
		for _, v := range v.Values.ToSlice().([]int) {
			operands = append(operands, float64(v))
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
		flags[i] = v == operandAsFloat
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *Float64Column) IsLike(pattern *regexp.Regexp) FilterType  {
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
func (c *Float64Column) Tx(op RowWiseFunc) Transformation {
	return Transformation{k: c.Title, v: op}
}

// Returns an Aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *Float64Column) Agg(aggFunc AggregateFunc) Aggregation {
	return Aggregation{c.Title: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *Float64Column) Order(option SortOrder) SortOption {
	return SortOption{c.Title: option}
}
