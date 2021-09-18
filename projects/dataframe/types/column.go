package types

import (
	"regexp"
)

const (
	IntType Datatype = iota
	FloatType
	StringType
	ObjectType
	BooleanType
	ArrayType
)


type Datatype int


type Column struct {
	Name string
	items orderedMapType
	Dtype Datatype
}

// Returns a list of Items
func (c *Column) Items() []interface{} {
	return c.items.ToSlice()
}


// Inserts a given value at the given index.
// If the index is beyond the length of keys,
// it fills the gap in both Items and keys with nil and "" respectively
func (c *Column) insert(index int, value interface{}) {
	// FIXME: This seems to be causing gaps
	nextIndex := len(c.items)

	if nextIndex <= index {
		for i := nextIndex; i <= index; i++ {
			c.items[i] = nil		
		}
	}

	c.items[index] = value
}

// Deletes many indices at once
func (c *Column) deleteMany(indices []int)  {
	for _, i := range indices {
		delete(c.items, i)
	}	
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is greater than operand or else false
// The operand can reference a constant, or a Col
func (c *Column) GreaterThan(operand float64) filterType {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		switch v := v.(type) {
		case int:
			flags[i] = float64(v) > operand
		case int8:
			flags[i] = float64(v) > operand
		case int16:
			flags[i] = float64(v) > operand
		case int32:
			flags[i] = float64(v) > operand
		case int64:
			flags[i] = float64(v) > operand
		case float32:
			flags[i] = float64(v) > operand
		case float64:
			flags[i] = v > operand
		default:
			flags[i] = false
		}
	}

	return flags
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is greater than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Column) GreaterOrEquals(operand float64) filterType {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		switch v := v.(type) {
		case int:
			flags[i] = float64(v) >= operand
		case int8:
			flags[i] = float64(v) >= operand
		case int16:
			flags[i] = float64(v) >= operand
		case int32:
			flags[i] = float64(v) >= operand
		case int64:
			flags[i] = float64(v) >= operand
		case float32:
			flags[i] = float64(v) >= operand
		case float64:
			flags[i] = v >= operand
		default:
			flags[i] = false
		}
	}

	return flags
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is less than operand or else false
// The operand can reference a constant, or a Col
func (c *Column) LessThan(operand float64) filterType {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		switch v := v.(type) {
		case int:
			flags[i] = float64(v) < operand
		case int8:
			flags[i] = float64(v) < operand
		case int16:
			flags[i] = float64(v) < operand
		case int32:
			flags[i] = float64(v) < operand
		case int64:
			flags[i] = float64(v) < operand
		case float32:
			flags[i] = float64(v) < operand
		case float64:
			flags[i] = v < operand
		default:
			flags[i] = false
		}
	}

	return flags
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is less than or equal to the operand or else false
// The operand can reference a constant, or a Col
func (c *Column) LessOrEquals(operand float64) filterType {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		switch v := v.(type) {
		case int:
			flags[i] = float64(v) <= operand
		case int8:
			flags[i] = float64(v) <= operand
		case int16:
			flags[i] = float64(v) <= operand
		case int32:
			flags[i] = float64(v) <= operand
		case int64:
			flags[i] = float64(v) <= operand
		case float32:
			flags[i] = float64(v) <= operand
		case float64:
			flags[i] = v <= operand
		default:
			flags[i] = false
		}
	}

	return flags
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is equal to operand or else false
// The operand can reference a constant, or a Col
func (c *Column) Equals(operand interface{}) filterType {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		flags[i] = v == operand
	}

	return flags
}

// Returns a map of an array of booleans corresponding in position to each item,
// true if item is like the regex expression or else false
func (c *Column) IsLike(pattern *regexp.Regexp) filterType  {
	count := len(c.items)
	flags := make(filterType, count)

	for i, v := range c.items {
		switch v := v.(type) {
		case string:
			flags[i] = pattern.MatchString(v)
		case []byte:
			flags[i] = pattern.Match(v)
		default:
			flags[i] = false
		}		
	}

	return flags
}

// Returns transformer method specific to this column to transform its values from one thing to another
// It is passed a function expecting a value any type
func (c *Column) Tx(op rowWiseFunc) transformation {
	return transformation{c.Name: op}
}

// Returns an aggregation function specific to this column to
// merge its values into a single value. It works when GroupBy is used
func (c *Column) Agg(aggFunc aggregateFunc) aggregation {
	return aggregation{c.Name: aggFunc}
}

// Returns a Sort Option that is attached to this column, for the given order
func (c *Column) Order(option sortOrder) sortOption {
	return sortOption{c.Name: option}
}
