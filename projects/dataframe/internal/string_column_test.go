package internal

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/internal/utils"
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
			expected: filterType{false, false, true, false, true, false},
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
		col.GreaterThan("1000")
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | Typed				    		| 1,640,126,487 ns/op| 90,068,110 B/op	 	 | 1307661 allocs/op   	 | x  	   |
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
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: "4", 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, true, false, true, false, true},
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
		col.GreaterOrEquals("1000")
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | Typed							| 1,472,384,774 ns/op	| 90,057,602 B/op	 	 | 1,307,610 allocs/op	 | x	   |
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
			expected: filterType{true, true, true, true, true, true},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{false, true, false, false, false, false},
		},
		{
			operand: "4", 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, false, true, false, true, false},
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
		col.LessThan("1000")
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 1,859,208,357 ns/op	| 103569454 B/op	 	 | 1525559 allocs/op     | x  	   |
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
			expected: filterType{true, true, true, true, true, true},
		},
		{
			operand: -2, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "-69", 4: "-2", 5: "67"},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "690", 4: "-2", 5: "67"},
			expected: filterType{true, true, true, false, false, false},
		},
		{
			operand: "4", 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, false, true, false, true, false},
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
		col.LessOrEquals("1000")
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 2,056,327,446 ns/op	| 103,567,231 B/op	 	| 1,525,548 allocs/op   | x  	  |
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
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: StringColumn{name: "foo", items: OrderedStringMapType{0: "23", 1: "60", 2: "-2", 3: "69"}}, 
			items: OrderedStringMapType{0: "23", 1: "6", 2: "-2", 3: "69", 4: "-2", 5: "67"},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: "0", 
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
		col.Equals("1000")
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				  | memory 				  | allocations			  | Choice  |
	// |--------------------------------|---------------------|-----------------------|-----------------------|---------|
	// | None				    		| 1,199,388,041 ns/op | 127,101,200 B/op	  | 103,4124 allocs/op    | x  	    |
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
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^\d`), 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{true, true, true, true, true, true},
		},
		{
			operand: regexp.MustCompile("^69"), 
			items: OrderedStringMapType{0: "23", 1: "500", 2: "2", 3: "69", 4: "0", 5: "67"},
			expected: filterType{false, false, false, true, false, false},
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
	// | None 							| 5,173,150,662 ns/op	| 508,459,324 B/op		| 22,576,725 allocs/op  | x		  |
}





