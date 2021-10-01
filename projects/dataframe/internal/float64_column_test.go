package internal

import (
	"regexp"
	"testing"
)

// insert for Float64Columns should fill any gaps in keys and Items with "", nil respectively
func TestFloat64Column_insert(t *testing.T)  {
	col := Float64Column{name: "hi", items: OrderedFloat64MapType{0: 6, 1: 70}}
	col.insert(4, 60)
	expectedItems := []float64{6.0, 70.0, 0.0, 0.0, 60.0}
	gotItems := col.Items().([]float64)
	
	for i := range expectedItems {
		got := gotItems[i]
		expected := expectedItems[i]
		if got != expected {
			t.Fatalf("Index %d had %v; expected %v", i, got, expected)
		}
	}
}

func BenchmarkFloat64Column_insert(b *testing.B)  {
	col := Float64Column{name: "hi", items: OrderedFloat64MapType{0: 6, 1: 70}}

	for i := 0; i < b.N; i++ {
		col.insert(4, 60)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 20.00 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestFloat64Column_GreaterThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: filterType{true, true, false, false, false, true},
		},
		{
			operand: Float64Column{name: "foo", items: OrderedFloat64MapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 690, 4: -2, 5: 67},
			expected: filterType{false, false, false, true, false, false},
		},
		{
			operand: IntColumn{name: "foo", items: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 690, 4: -2, 5: 67},
			expected: filterType{false, false, false, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.GreaterThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkFloat64Column_GreaterThan(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 919,153,136 ns/op	 | 37,289,891 B/op	     | 12,302 allocs/op      | x  	   |
}

// GreaterOrEquals should return a slice of booleans where true is for values greater or equal to the value,
// false is for otherwise
func TestFloat64Column_GreaterOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: filterType{true, true, true, false, true, true},
		},
		{
			operand: Float64Column{name: "foo", items: OrderedFloat64MapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{true, true, false, true, false, true},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.GreaterOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkFloat64Column_GreaterOrEquals(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.GreaterOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				    		| 990,669,059 ns/op	 | 41142735 B/op	     | 13943 allocs/op       | x  	   |
}

// LessThan should return a slice of booleans where true is for values less than the value,
// false is for otherwise
func TestFloat64Column_LessThan(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: filterType{false, false, false, true, false, false},
		},
		{
			operand: Float64Column{name: "foo", items: OrderedFloat64MapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{false, true, false, false, false, false},
		},
		{
			operand: IntColumn{name: "foo", items: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{false, true, false, false, false, false},
		},
		{
			operand: 4, 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.LessThan(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkFloat64Column_LessThan(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.LessThan(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				 | allocations			 | Choice  |
	// |--------------------------------|-----------------------|------------------------|-----------------------|---------|
	// | None				    		| 941,402,221 ns/op		| 38462338 B/op	   		 | 12767 allocs/op       | x  	   |
}

// LessOrEquals should return a slice of booleans where true is for values less or equal to the value,
// false is for otherwise
func TestFloat64Column_LessOrEquals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: -69, 4: -2, 5: 67},
			expected: filterType{false, false, true, true, true, false},
		},
		{
			operand: Float64Column{name: "foo", items: OrderedFloat64MapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 690, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: IntColumn{name: "foo", items: OrderedIntMapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 690, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: 4, 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, true, false, true, false},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.LessOrEquals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkFloat64Column_LessOrEquals(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.LessOrEquals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 939,445,430 ns/op		| 38465104 B/op	   		| 12792 allocs/op       | x  	  |
}

// Equals should return a slice of booleans where true is for values equal to the value,
// false is for otherwise
func TestFloat64Column_Equals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: -2, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{false, false, true, false, true, false},
		},
		{
			operand: Float64Column{name: "foo", items: OrderedFloat64MapType{0: 23, 1: 60, 2: -2, 3: 69}}, 
			items: OrderedFloat64MapType{0: 23, 1: 6, 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{true, false, true, true, false, false},
		},
		{
			operand: 0, 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.Equals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}
}

func BenchmarkFloat64Column_Equals(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.Equals(1000)
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				  | memory 				  | allocations			  | Choice  |
	// |--------------------------------|---------------------|-----------------------|-----------------------|---------|
	// | None				    		| 334,640,964 ns/op	  | 20051509 B/op	      | 4797 allocs/op        | x  	    |
}


// IsLike should return a slice of booleans where true is for values that match the regexp pattern passed,
// false is for otherwise
func TestFloat64Column_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items OrderedFloat64MapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^\d`), 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{true, true, true, true, true, true},
		},
		{
			operand: regexp.MustCompile("^Duhaga"), 
			items: OrderedFloat64MapType{0: 23, 1: 500, 2: 2, 3: 69, 4: 0, 5: 67},
			expected: filterType{false, false, false, false, false, false},
		},
	}

	for index, tr := range testData {
		col := Float64Column{name: "hi", items: tr.items}
		output := col.IsLike(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkFloat64Column_IsLike(b *testing.B)  {
	items := OrderedFloat64MapType{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = float64(i)
	}

	col := Float64Column{name: "hi", items: items}

	for i := 0; i < b.N; i++ {
		col.IsLike(regexp.MustCompile("^10"))
	}

	// Results:
	// ========
	// benchtime=10s
	// 
	// | Change 						| time				 	| memory 				| allocations			| Choice  |
	// |--------------------------------|-----------------------|-----------------------|-----------------------|---------|
	// | None				    		| 6,285,800,090 ns/op	| 393,864,724 B/op		| 18,077,117 allocs/op  | x  	  |
}





