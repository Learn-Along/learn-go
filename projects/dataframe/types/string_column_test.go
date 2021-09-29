package types

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for StringColumns should fill any gaps in keys and Items with "", nil respectively
func TestStringColumn_insert(t *testing.T)  {
	col := StringColumn{name: "hi", items: OrderedStringMapType{0: "6", 1: "70"}}
	col.insert(4, "60")
	expectedItems := []string{"6", "70", "", "", "60"}

	if !utils.AreStringSliceEqual(expectedItems, col.Items().([]string)) {
		t.Fatalf("items expected: %v, got %v", expectedItems, col.Items())
	}
}

func BenchmarkStringColumn_insert(b *testing.B)  {
	col := StringColumn{name: "hi", items: OrderedStringMapType{0: "6", 1: "70"}}

	for i := 0; i < b.N; i++ {
		col.insert(4, "60")
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 21.22 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestStringColumn_GreaterThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "4", 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, false, true, false, true, false},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, true, false, false, false, false},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.GreaterThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkStringColumn_GreaterThan(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 855,400,310 ns/op	 | 97,572,326 B/op	     | 775,572 allocs/op     | x  	   |
	// | Add goroutine in for loop		| 4,449,787,656 ns/op| 363,255,202 B/op	     | 3,102,174 allocs/op   |		   |
	// | With wrapper around goroutine	| 4,437,230,299 ns/op| 363251869 B/op	     | 3102156 allocs/op 	 |         |
	// | With wait groups 				| 4,067,743,934 ns/op| 714,285,405 B/op	     | 8,164,777 allocs/op   |         |
}

// GreaterOrEquals should return a slice of booleans where true is for values greater or equal to the value,
// false is for otherwise
func TestStringColumn_GreaterOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{false, false, true, true, true, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, true, true, false, false, false},
		},
		{
			operand: 4, 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.GreaterOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkStringColumn_GreaterOrEquals(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 871,616,373 ns/op	 | 97,562,900 B/op	     | 775,526 allocs/op     | x  	   |
}

// LessThan should return a slice of booleans where true is for values less than the value,
// false is for otherwise
func TestStringColumn_LessThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{true, true, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{false, true, false, false, false, false},
		},
		{
			operand: 4, 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.LessThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkStringColumn_LessThan(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.LessThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 1,142,363,432 ns/op	|	97,571,606 B/op	     | 775,569 allocs/op     | x  	   |
}

// LessOrEquals should return a slice of booleans where true is for values less or equal to the value,
// false is for otherwise
func TestStringColumn_LessOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{true, true, true, false, true, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.LessOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkStringColumn_LessOrEquals(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.LessOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 5,437,853,397 ns/op	| 540,413,572 B/op	    | 4653405 allocs/op     | x  	  |
}

// Equals should return a slice of booleans where true is for values equal to the value,
// false is for otherwise
func TestStringColumn_Equals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{false, false, true, false, true, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: 0, 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.Equals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkStringColumn_Equals(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.Equals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				  | memory 				  | allocations			  | Choice  |
	// |--------------------------------|---------------------|-----------------------|-----------------------|---------|
	// | None				    		| 1,199,388,041 ns/op | 127101200 B/op	      | 1034124 allocs/op     | x  	    |
}


// IsLike should return a slice of booleans where true is for values that match the regexp pattern passed,
// false is for otherwise
func TestStringColumn_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items OrderedStringMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, true, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^\d`), 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, true, true, true, true},
		},
		{
			operand: regexp.MustCompile("^Duhaga"), 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
	}

	for index, tr := range testData {
		col := StringColumn{name: "hi", items: tr.items}
		output := col.IsLike(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkStringColumn_IsLike(b *testing.B)  {
	items := OrderedStringMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := StringColumn{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.IsLike(regexp.MustCompile("^10"))
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 12,941,657,683 ns/op	| 287,952,280 B/op		| 27,307,263 allocs/op  | x  	  |
}





