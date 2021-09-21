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

/*
* Benchmark tests
*/
func BenchmarkColumn_GreaterThan(b *testing.B)  {
	items := map[int]interface{}{}

	for i := 0; i < 9000000; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: StringType, items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    	| 791834066 ns/op	 | 48366771 B/op	     | 344697 allocs/op      | x  	   |
	// | Add goroutine in for loop	| 13898155496 ns/op	 | 1321728940 B/op	     | 23181212 allocs/op  	 |		   |
}
