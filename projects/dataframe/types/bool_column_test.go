package types

import (
	"regexp"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for BoolColumns should fill any gaps in keys and Items with "", nil respectively
func TestBoolColumn_insert(t *testing.T)  {
	col := BoolColumn{name: "hi", items: OrderedBoolMapType{0: false, 1: true}}
	col.insert(4, true)
	expectedItems := []interface{}{false, true, false, false, true}

	if !utils.AreSliceEqual(expectedItems, col.Items().([]interface{})) {
		t.Fatalf("items expected: %v, got %v", expectedItems, col.Items())
	}
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestBoolColumn_GreaterThan(t *testing.T)  {
	operand := 2
	col := BoolColumn{name: "hi", items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true}}
	expected := filterType{true, false, true, true, false, true}
	output := col.GreaterThan(operand)

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkBoolColumn_GreaterThan(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{name: "hi", items: items}

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
func TestBoolColumn_GreaterOrEquals(t *testing.T)  {
	operand := 2
	col := BoolColumn{name: "hi", items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true}}
	expected := filterType{true, false, true, true, false, true}
	output := col.GreaterOrEquals(operand)

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkBoolColumn_GreaterOrEquals(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{name: "hi", items: items}

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
func TestBoolColumn_LessThan(t *testing.T)  {
	operand := 2
	col := BoolColumn{name: "hi", items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true}}
	expected := filterType{true, false, true, true, false, true}
	output := col.LessThan(operand)

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkBoolColumn_LessThan(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = true
	}

	col := BoolColumn{name: "hi", items: items}

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
func TestBoolColumn_LessOrEquals(t *testing.T)  {
	operand := 2
	col := BoolColumn{name: "hi", items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true}}
	expected := filterType{true, false, true, true, false, true}
	output := col.LessOrEquals(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkBoolColumn_LessOrEquals(b *testing.B)  {
	items := OrderedBoolMapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = false
	}

	col := BoolColumn{name: "hi", items: items}

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
func TestBoolColumn_Equals(t *testing.T)  {
	type testRecord struct {
		operand LiteralOrColumn;
		items OrderedBoolMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: false, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{false, true, false, false, true, false},
		},
		{
			operand: BoolColumn{name: "hoo", items: OrderedBoolMapType{0: false, 1: true, 2: true, 3: false, 4: false, 5: false}}, 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{name: "hi", items: tr.items}
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

	col := BoolColumn{name: "hi", items: items}

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
func TestBoolColumn_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items OrderedBoolMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^true$`), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{true, false, true, true, false, true},
		},
		{
			operand: regexp.MustCompile("^fa"), 
			items: OrderedBoolMapType{0: true, 1: false, 2: true, 3: true, 4: false, 5: true},
			expected: filterType{false, true, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := BoolColumn{name: "hi", items: tr.items}
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

	col := BoolColumn{name: "hi", items: items}

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





