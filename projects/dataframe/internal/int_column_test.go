package internal

import (
	"regexp"
	"testing"
)

// Insert for IntColumns should fill any gaps in keys and Items with "", nil respectively
func TestIntColumn_insert(t *testing.T)  {
	col := IntColumn{Title: "hi", Values: OrderedIntMapType{0: 6, 1: 70}}
	col.Insert(4, 60)
	expectedItems := []int{6, 70, 0, 0, 60}
	gotItems := col.Items().([]int)
	
	for i := range expectedItems {
		got := gotItems[i]
		expected := expectedItems[i]
		if got != expected {
			t.Fatalf("Index %d had %v; expected %v", i, got, expected)
		}
	}
}

func BenchmarkIntColumn_insert(b *testing.B)  {
	col := IntColumn{Title: "hi", Values: OrderedIntMapType{0: 6, 1: 70}}

	for i := 0; i < b.N; i++ {
		col.Insert(4, 60)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 17.84 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestIntColumn_GreaterThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedIntMapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: FilterType{true, true, false, false, false, true},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: 690, 4: -2, 5: 67},
			expected: FilterType{false, false, false, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.GreaterThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkIntColumn_GreaterThan(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 749,171,660 ns/op	 | 31,813,476 B/op	     | 9907 allocs/op        | x  	   |
}

// GreaterOrEquals should return a slice of booleans where true is for values greater or equal to the value,
// false is for otherwise
func TestIntColumn_GreaterOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedIntMapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: FilterType{true, true, true, false, true, true},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: FilterType{true, false, true, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.GreaterOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkIntColumn_GreaterOrEquals(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.GreaterOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 784,941,139 ns/op	 | 36,202,119 B/op	     | 11,826 allocs/op     	 | x  	   |
}

// LessThan should return a slice of booleans where true is for values less than the value,
// false is for otherwise
func TestIntColumn_LessThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedIntMapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: FilterType{false, false, false, true, false, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: FilterType{false, true, false, false, false, false},
		},
		{
			operand: 4, 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.LessThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkIntColumn_LessThan(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.LessThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 795,008,458 ns/op		| 32574746 B/op	   		 | 10243 allocs/op       | x  	   |
}

// LessOrEquals should return a slice of booleans where true is for values less or equal to the value,
// false is for otherwise
func TestIntColumn_LessOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedIntMapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: FilterType{false, false, true, true, true, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedIntMapType{0: 23, 1: 690, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: FilterType{true, false, true, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.LessOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkIntColumn_LessOrEquals(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.LessOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 756,376,738 ns/op	    | 31,814,971 B/op	    	| 9,917 allocs/op       | x  	  |
}

// Equals should return a slice of booleans where true is for values equal to the value,
// false is for otherwise
func TestIntColumn_Equals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedIntMapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: FilterType{false, false, true, false, true, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedIntMapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: FilterType{true, false, true, true, false, false},
		},
		{
			operand: 0, 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.Equals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkIntColumn_Equals(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.Equals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				  | memory 				  | allocations			  | Choice  |
	// |--------------------------------|---------------------|-----------------------|-----------------------|---------|
	// | None				    		| 244,137,151 ns/op	  | 16,773,965 B/op	    	  | 3,378 allocs/op        | x  	    |
}


// IsLike should return a slice of booleans where true is for values that match the regexp pattern passed,
// false is for otherwise
func TestIntColumn_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items OrderedIntMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^\d`), 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{true, true, true, true, true, true},
		},
		{
			operand: regexp.MustCompile("^Duhaga"), 
			items: OrderedIntMapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: FilterType{false, false, false, false, false, false},
		},
	}

	for index, tr := range testData {
		col := IntColumn{Title: "hi", Values: tr.items}
		output := col.IsLike(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkIntColumn_IsLike(b *testing.B)  {
	items := OrderedIntMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := IntColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.IsLike(regexp.MustCompile("^10"))
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 3,343,612,836 ns/op	| 254037489 B/op	    | 18043603 allocs/op    | x  	  |
}





