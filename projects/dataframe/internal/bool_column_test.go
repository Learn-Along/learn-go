package internal

import (
	"regexp"
	"testing"
)

// Insert for BoolColumns should fill any gaps in keys and Items with "", nil respectively
func TestBoolColumn_insert(t *testing.T)  {
	col := BoolColumn{Title: "hi", Values: OrderedBoolMapType{0: false, 1: true}}
	col.Insert(4, true)
	expectedItems := []bool{false, true, false, false, true}
	gotItems := col.Items().([]bool)
	
	for i := range expectedItems {
		got := gotItems[i]
		expected := expectedItems[i]
		if got != expected {
			t.Fatalf("Index %d had %v; expected %v", i, got, expected)
		}
	}
}

func BenchmarkBoolColumn_insert(b *testing.B)  {
	col := BoolColumn{Title: "hi", Values: OrderedBoolMapType{0: false, 1: true}}

	for i := 0; i < b.N; i++ {
		col.Insert(4, true)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 20.15 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestBoolColumn_GreaterThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{true, false, true, true, false, true},
		},
		{
			operand: BoolColumn{Title: "foo", Values: OrderedBoolMapType{0: true, 1: false, 2: false, 3: true}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, true, false, false, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: -1, 1: 7, 2: 6, 3: 5}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: true, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.GreaterThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkBoolColumn_GreaterThan(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 2,216,969 ns/op	 | 9,046,912 B/op	     | 31 allocs/op          | x  	   |
}

// GreaterOrEquals should return a slice of booleans where true is for values greater or equal to the value,
// false is for otherwise
func TestBoolColumn_GreaterOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{true, true, true, true, true, true},
		},
		{
			operand: BoolColumn{Title: "foo", Values: OrderedBoolMapType{0: true, 1: false, 2: false, 3: true}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{true, true, true, true, false, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: -1, 1: 7, 2: 6, 3: 5}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: true, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{true, false, true, true, false, true},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.GreaterOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkBoolColumn_GreaterOrEquals(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.GreaterOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 2,249,874 ns/op	 | 9046041 B/op	      	 | 30 allocs/op          | x  	   |
}

// LessThan should return a slice of booleans where true is for values less than the value,
// false is for otherwise
func TestBoolColumn_LessThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: BoolColumn{Title: "foo", Values: OrderedBoolMapType{0: true, 1: false, 2: false, 3: true}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{ false, false, false, false, false, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: -1, 1: 7, 2: 6, 3: 5}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: true, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{ false, true, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.LessThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkBoolColumn_LessThan(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = true
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.LessThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 2,198,321 ns/op	 	| 9,046,524 B/op	     | 31 allocs/op          | x  	   |
}

// LessOrEquals should return a slice of booleans where true is for values less or equal to the value,
// false is for otherwise
func TestBoolColumn_LessOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, true, false, false, true, false},
		},
		{
			operand: BoolColumn{Title: "foo", Values: OrderedBoolMapType{0: true, 1: false, 2: false, 3: true}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{ true, true, false, true, false, false},
		},
		{
			operand: IntColumn{Title: "foo", Values: OrderedIntMapType{0: -1, 1: 7, 2: 6, 3: 5}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: true, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{ true, true, true, true, true, true},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.LessOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkBoolColumn_LessOrEquals(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.LessOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 2,211,790 ns/op	 	| 9,047,361 B/op	    | 31 allocs/op          | x  	  |
}

// Equals should return a slice of booleans where true is for values equal to the value,
// false is for otherwise
func TestBoolColumn_Equals(t *testing.T)  {
	type testRecord struct {
		operand LiteralOrColumn;
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, true, false, false, true, false},
		},
		{
			operand: BoolColumn{Title: "hoo", Values: OrderedBoolMapType{0: false, 1: true, 2: true, 3: false, 4: false, 5: false}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.Equals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkBoolColumn_Equals(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = true
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.Equals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				  | memory 				  | allocations			  | Choice  |
	// |--------------------------------|---------------------|-----------------------|-----------------------|---------|
	// | None				    		| 2,202,191 ns/op	  | 9,046,113 B/op	      | 30 allocs/op          | x  	    |
}


// IsLike should return a slice of booleans where true is for values that match the regexp pattern passed,
// false is for otherwise
func TestBoolColumn_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items OrderedBoolMapType;
		expected FilterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, false, false, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^true$`), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{true, false, true, true, false, true},
		},
		{
			operand: regexp.MustCompile("^fa"), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: FilterType{false, true, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{Title: "hi", Values: tr.items}
		output := col.IsLike(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkBoolColumn_IsLike(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = true
	}

	col := BoolColumn{Title: "hi", Values: items}

	for i := 0; i < b.N; i++ {
		col.IsLike(regexp.MustCompile("^tru"))
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 3,588,715,970 ns/op	| 118816362 B/op	    | 9051141 allocs/op     | x  	  |
}





