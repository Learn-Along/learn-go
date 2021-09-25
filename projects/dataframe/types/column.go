package types

import (
	"regexp"
	"runtime"
)

const (
	IntType Datatype = iota
	FloatType
	StringType
	ObjectType
	BooleanType
	ArrayType
)

var MAX_PROCS = runtime.GOMAXPROCS(0)

type Datatype int

// type boolTuple struct {i int; v bool}

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

// Returns an array of booleans corresponding in position to each item,
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

// Returns an array of booleans corresponding in position to each item,
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

// // Same as @GreaterOrEquals but optimized to take advantage of multicore systems
// func (c *Column) XGreaterOrEquals(operand float64) filterType {
// 	count := len(c.items)
// 	flags := make(filterType, count)
// 	type pRange struct{start int; stop int}
// 	queue := make(chan pRange)
// 	results := make(chan boolTuple)
// 	// For cache size of 64 KB, cache line of 64 B, and knowing that size of interface{} is 16 B
// 	portionSize := 40

// 	// create MAX_PROCS goroutines only 
// 	for i := 0; i < MAX_PROCS; i++ {
// 		go func (input chan pRange, output chan boolTuple)  {
// 			for p := range input {
// 				for index := p.start; index < p.stop; index++ {
// 					switch v := c.items[index].(type) {
// 					case int:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case int8:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case int16:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case int32:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case int64:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case float32:
// 						output <- boolTuple{i: index, v: float64(v) >= operand}
// 					case float64:
// 						output <- boolTuple{i: index, v: v >= operand}
// 					default:
// 						output <- boolTuple{i: index, v: false}
// 					}
// 				}				
// 			}			
// 		}(queue, results)
// 	}

// 	// Sort of the scheduler the portions to pass to each goroutine
// 	go func (output chan pRange)  {
// 		for i := 0; i < count; i += portionSize {
// 			stop := i + portionSize
// 			if stop > count {
// 				stop = count
// 			}

// 			output <- pRange{start: i, stop: stop}		
// 		}
// 		close(output)		
// 	}(queue)

// 	for i := 0; i < count; i++ {
// 		r := <- results
// 		flags[r.i] = r.v
// 	}

// 	return flags
// }

// Returns an array of booleans corresponding in position to each item,
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

// Returns an array of booleans corresponding in position to each item,
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

// Returns an array of booleans corresponding in position to each item,
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

// Returns an array of booleans corresponding in position to each item,
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
	return transformation{k: c.Name, v: op}
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
