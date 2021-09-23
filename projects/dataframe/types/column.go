package types

import (
	"math"
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

type comparator func(value interface{}) bool


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
	// type boolTuple struct {i int; v bool}
	// count := len(c.items)
	// flags := make(filterType, count)	
	// queue := make(chan boolTuple)
	// n := runtime.GOMAXPROCS(0)
	// portionSize := int(math.Ceil(float64(count)/float64(n)))
	// // portionSize := 4

	// // one goroutine for each process, to utilize the cache to the max
	// // https://appliedgo.net/concurrencyslower/
	// for i := 0; i < n; i++ {
	// 	// goroutines should not update global variables so as to avoid heap memory,
	// 	// but you can read them
	// 	start := i * portionSize
	// 	stop := start + portionSize
		
	// 	go func(start int, stop int, maxCount int, chn chan boolTuple) {
	// 		for i := start; i < stop && i < maxCount; i++ {
	// 			switch v := c.items[i].(type) {
	// 			case int:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case int8:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case int16:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case int32:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case int64:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case float32:
	// 				chn <- boolTuple{i: i, v: float64(v) > operand}
	// 			case float64:
	// 				chn <- boolTuple{i: i, v: v > operand}
	// 			default:
	// 				chn <- boolTuple{i: i, v: false}
	// 			}
	// 		}			
	// 	}(start, stop, count, queue)		
	// }

	// for i := 0; i < count; i++ {
	// 	t := <- queue
	// 	flags[t.i] = t.v		
	// }

	// return flags
	// count := len(c.items)
	// flags := make(filterType, count)

	// for i, v := range c.items {
	// 	switch v := v.(type) {
	// 	case int:
	// 		flags[i] = float64(v) > operand
	// 	case int8:
	// 		flags[i] = float64(v) > operand
	// 	case int16:
	// 		flags[i] = float64(v) > operand
	// 	case int32:
	// 		flags[i] = float64(v) > operand
	// 	case int64:
	// 		flags[i] = float64(v) > operand
	// 	case float32:
	// 		flags[i] = float64(v) > operand
	// 	case float64:
	// 		flags[i] = v > operand
	// 	default:
	// 		flags[i] = false
	// 	}
	// }

	// return flags
	return c.compareWithItems(func(v interface{}) bool {
		switch v := v.(type) {
		case int:
			return float64(v) > operand
		case int8:
			return float64(v) > operand
		case int16:
			return float64(v) > operand
		case int32:
			return float64(v) > operand
		case int64:
			return float64(v) > operand
		case float32:
			return float64(v) > operand
		case float64:
			return v > operand
		default:
			return false
		}
	})
}

// Utility to run comparison methods faster on the given items
func (c *Column) compareWithItems(compare comparator) filterType {
	type boolTuple struct {i int; v bool}
	count := len(c.items)
	flags := make(filterType, count)	
	n := 3 // runtime.GOMAXPROCS(0)
	// if n < 1 {
	// 	n = 1
	// }
	// each interface{} is about 16 bytes
	assumedCacheLineSize := 64
	portionSize := int(assumedCacheLineSize / 16)
	queue := make(chan boolTuple, count)

	// one goroutine for each process, to utilize the cache to the max
	// https://appliedgo.net/concurrencyslower/
	for i := 0; i < n; i++ {
		// goroutines should not update global variables so as to avoid heap memory,
		// but you can read them
		// Split the array to the given portion size
		start := i * portionSize
		stop := int(math.Min(float64(start + portionSize), float64(count)))

		if start >= count {
			break
		}

		// create n or less goroutines
		go func(start int, stop int, maxCount int, chn chan boolTuple) {
			// each goroutine pushes to the channel a number of times equal to portionSize
			for i := start; i < stop; i++ {
				chn <- boolTuple{i: i, v: compare(c.items[i])}
			}			
		}(start, stop, count, queue)		
	}

	for i := 0; i < count; i++ {
		t := <- queue
		flags[t.i] = t.v		
	}

	return flags	
}

// // Returns an array of booleans corresponding in position to each item,
// // true if item is greater than operand or else false
// // The operand can reference a constant, or a Col
// func (c *Column) GreaterThan(operand float64) filterType {
// 	count := len(c.items)
// 	flags := make(filterType, count)	
// 	ch := make(chan boolTuple)
// 	n := runtime.GOMAXPROCS(0)
// 	portionSize := int(math.Ceil(float64(count)/float64(n)))

// 	// one goroutine for each process, to utilize the cache to the max
// 	// https://appliedgo.net/concurrencyslower/
// 	for i := 0; i < n; i++ {
// 		// goroutines should not update global variables so as to avoid heap memory,
// 		// but can read them
// 		startIndex := i * portionSize
// 		stopIndex := startIndex + portionSize

// 		go func(start int, stop int, maxCount int, chn chan boolTuple) {

// 			for i := start; i < stop && i < maxCount; i++ {
// 				switch v := c.items[i].(type) {
// 				case int:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case int8:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case int16:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case int32:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case int64:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case float32:
// 					chn <- boolTuple{i: i, v: float64(v) > operand}
// 				case float64:
// 					chn <- boolTuple{i: i, v: v > operand}
// 				default:
// 					chn <- boolTuple{i: i, v: false}
// 				}
// 			}
			
// 		}(startIndex, stopIndex, count, ch)
		
// 	}

// 	for i := 0; i < count; i++ {
// 		t := <- ch
// 		flags[t.i] = t.v		
// 	}	

// 	return flags
// }

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
