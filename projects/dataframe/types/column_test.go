package types

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for columns should fill any gaps in keys and Items with "", nil respectively
func TestColumn_insert(t *testing.T)  {
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{0: "hi", 1: "wow"}}
	col.insert(4, "yeah")
	expectedItems := []interface{}{"hi", "wow", nil, nil, "yeah"}

	if !utils.AreSliceEqual(expectedItems, col.Items()) {
		t.Fatalf("items expected: %v, got %v", expectedItems, col.Items())
	}
}

// Greater that should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestColumn_GreaterThan(t *testing.T)  {
	operand := 2
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{
		0: 23.4, 1: 10, 2: -2, 3: 69, 4: 0.23, 5: 67}}
	expected := filterType{true, true, false, true, false, true}
	output := col.GreaterThan(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

/*
* Benchmark tests
*/
func BenchmarkColumn_GreaterThan(b *testing.B)  {
	items := map[int]interface{}{}
	// numberOfItems := 9000000
	// FIXME: deadlock seems to happen at > 16
	// Note that portionSize is constant at 4
	// I am not sure if it has something to do with the number of cores (procs) = 4
	// so probably with 16 items, all items are consumed in the 4 goroutines 
	// but when the number is beyond 16, 4 goroutines are created but the chan sends are more than 
	// 1 for one of the goroutines. The issue is why then do buffered channels also suffer. 
	// I expect they should be able to wait with the buffer.
	// Note: when the n is set to 3 (hard coded), and number of items is above 12, deadlock!
	// It looks like the max items allowable without deadlock = number of procs * portionSize
	// number of procs corresponds to number of goroutines
	// portionSize corresponds to number of calls to push to channel per routine
	// It looks like for safety, 
	// each goroutine can only push to the channel a number of times equal to portionSize before deadlock occurs
	// I need to think about that.
	numberOfItems := 13

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: StringType, items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 855,400,310 ns/op	 | 97,572,326 B/op	     | 775,572 allocs/op     | x  	   |
	// | Add goroutine in for loop		| 4,449,787,656 ns/op| 363,255,202 B/op	     | 3,102,174 allocs/op   |		   |
	// | With wrapper around goroutine	| 4,437,230,299 ns/op| 363251869 B/op	     | 3102156 allocs/op 	 |         |
}
