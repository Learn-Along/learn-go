package internal

import (
	"fmt"
	"regexp"
)

type BoolColumn struct {
	Title string
	Values OrderedBoolMapType
}

// returns the Title of the column
func (c *BoolColumn) Name() string {
	return c.Title
}

// Number of Values in int column
func (c *BoolColumn) Len() int {
	return c.Values.Len()
}

// Number of Values in int column
func (c *BoolColumn) ItemAt(index int) Item {
	return c.Values[index]
}

// Returns a list of Items
func (c *BoolColumn) Items() ItemSlice {
	return c.Values.ToSlice()
}

// Returns the data type of the given column
func (c *BoolColumn) GetDatatype() Datatype {
	return BoolType
}

// Reorders the OrderedMapType ensuring that any gaps in the data are removed
// So as to go back to a sequantial key list
func (c *BoolColumn) Defragmentize(newOrder []int) {
	c.Values.Defragmentize(newOrder)
}

// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
// It ignores value is not a boolean
func (c *BoolColumn) Insert(index int, value Item) {
	nextIndex := c.Values.Len()

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.Values[i] = false		
		}
	}

	switch v := value.(type) {
	case bool:
		c.Values[index] = v
	}	
}

// Deletes many indices at once
func (c *BoolColumn) DeleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.Values, i)
	}	
}

// Returns an array of booleans corresponding in position to each item,
// true if item is true and operand is false or else false
// The operand can reference a constant, or a Col
func (c *BoolColumn) GreaterThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsBool bool
	var operands []bool

	switch v := operand.(type) {
	case bool:
		operandAsBool = v
	case BoolColumn:
		operands = v.Values.ToSlice().([]bool)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = !op && v
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = !operandAsBool && v
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if operand is false, or if operand is true and item is true or else false
// The operand can reference a constant, or a Col
func (c *BoolColumn) GreaterOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsBool bool
	var operands []bool

	switch v := operand.(type) {
	case bool:
		operandAsBool = v
	case BoolColumn:
		operands = v.Values.ToSlice().([]bool)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = !op || (v && op)
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = !operandAsBool || (v && operandAsBool)
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is false, and operand is true or else false
// The operand can reference a constant, or a Col
func (c *BoolColumn) LessThan(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsBool bool
	var operands []bool

	switch v := operand.(type) {
	case bool:
		operandAsBool = v
	case BoolColumn:
		operands = v.Values.ToSlice().([]bool)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = !v && op
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = !v && operandAsBool
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if operand is true or if both operand and item are false or else false
// The operand can reference a constant, or a Col
func (c *BoolColumn) LessOrEquals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operandAsBool bool
	var operands []bool

	switch v := operand.(type) {
	case bool:
		operandAsBool = v
	case BoolColumn:
		operands = v.Values.ToSlice().([]bool)
	default:
		return flags
	}

	if operands != nil {
		for i, op := range operands {
			if v, ok := c.Values[i]; ok {
				flags[i] = op || (!v && !op)
			}
		}

		return flags
	}

	for i := 0; i < count; i++ {
		if v, ok := c.Values[i]; ok {
			flags[i] = operandAsBool || (!v && !operandAsBool)
		}
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *BoolColumn) Equals(operand LiteralOrColumn) FilterType {
	count := len(c.Values)
	flags := make(FilterType, count)
	var operands []bool
	var operandAsBool bool

	switch v := operand.(type) {
	case bool:
		operandAsBool = v
	case BoolColumn:
		operands = v.Values.ToSlice().([]bool)
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
		flags[i] = v == operandAsBool
	}

	return flags
}

// Returns an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *BoolColumn) IsLike(pattern *regexp.Regexp) FilterType  {
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
func (c *BoolColumn) Tx(op RowWiseFunc) Transformation {
	return Transformation{k: c.Title, v: op}
}

// Returns an Aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *BoolColumn) Agg(aggFunc AggregateFunc) Aggregation {
	return Aggregation{c.Title: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *BoolColumn) Order(option SortOrder) SortOption {
	return SortOption{c.Title: option}
}
