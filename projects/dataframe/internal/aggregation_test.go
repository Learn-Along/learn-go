package internal

import "testing"

// GetMax should return the maximum value (as a float64 value) in a given list of items
func TestGetMax(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "yoohoo", "salut"},
			expected: "yoohoo",
		},
		{
			input: []float64{1.9, 2.3, 4.8, 0.8},
			expected: 4.8,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 98,
		},		
	}

	for _, tr := range testData {
		got := GetMax(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkGetMax(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetMax(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 57.17 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// GetMin should return the minimum value (as a float64 value) in a given list of items
func TestGetMin(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "yoohoo", "salut"},
			expected: "hello",
		},
		{
			input: []float64{1, 2.3, 4, 0.8},
			expected: 0.8,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 4,
		},		{
			input: []float64{80.7, 78.6, 98.5, 98.509},
			expected: 78.6,
		},		
	}

	for loop, tr := range testData {
		got := GetMin(tr.input)
		if got != tr.expected {
			t.Fatalf("loop %d: expected %v; got %v", loop, tr.expected, got)
		}
	}
}

func BenchmarkGetMin(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetMin(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 66.58 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// GetSum should return the sum (as a float64 value) of the given list of items, but will return nil if 
// the values are not all numbers (nil values are treated as zero)
func TestGetSum(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item;
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "yoohoo", "salut"},
			expected: nil,
		},
		{
			input: []float64{1, 2.3, 4, 0.8},
			expected: 8.1,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 260,
		},
		{
			input: []float64{80.7, 78.6, 98.5, 98.509},
			expected: 356.309,
		},
		
	}

	for _, tr := range testData {
		got := GetSum(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkGetSum(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetSum(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 63.92 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// GetMean should return the mean (as a float64 value) of the given list of items.
// It returns nil if the values are not numbers (nil values are treated as zero)
func TestGetMean(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "yoohoo", "salut"},
			expected: nil,
		},
		{
			input: []float64{1, 2.3, 4, 0.8},
			expected: 2.025,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 65.0,
		},
		{
			input: []float64{80.7, 78.6, 98.5, 98.509},
			expected: 89.07725,
		},		
	}

	for _, tr := range testData {
		got := GetMean(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkGetMean(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetMean(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 85.74 ns/op	     | 40 B/op	       		 | 3 allocs/op           |  x  	   |
}

// GetCount should return the number of items in the list of values, including nils
func TestGetCount(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item;
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "hello", "hello", "yoohoo", "salut", "hello", "yoohoo", "salut", "yoohoo"},
			expected: 10,
		},
		{
			input: []float64{1, 2.3, 4, 0.8, 2.3, 4},
			expected: 6,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 4,
		},		{
			input: []float64{80.7, 78.6, 98.5, 98.509},
			expected: 4,
		},		
	}

	for _, tr := range testData {
		got := GetCount(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkGetCount(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetCount(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 42.17 ns/op	     | 24 B/op	       		 | 1 allocs/op           |  x  	   |
}

// GetRange should return the range (as a float64 value) of the given list of items
// i.e. maximum minus minimum.
// It returns nil if the values are not numbers (nil values are ignored)
func TestGetRange(t *testing.T)  {
	type testRecord struct {
		input ItemSlice;
		expected Item
	}

	testData := []testRecord{
		{
			input: []string{"hi", "hello", "yoohoo", "salut"},
			expected: nil,
		},
		{
			input: []float64{1, 2.3, 4, 0.8},
			expected: 3.2,
		},
		{
			input: []int{80, 78, 98, 4},
			expected: 94,
		},
		{
			input: []float64{80.7, 78.6, 98.5, 98.509},
			expected: float64(98.509) - float64(78.6),
		},
		
	}

	for loop, tr := range testData {
		got := GetRange(tr.input)
		if got != tr.expected {
			t.Fatalf("loop %d: expected %v; got %v", loop, tr.expected, got)
		}
	}
}

func BenchmarkGetRange(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		GetRange(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 70.68 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// MergeAggregations should merge an Aggregation list into a single Aggregation 
// ensuring that the last AggregateFunc to be attached to a given column is the one kept,
// the previous ones are overwritten, to avoid ambiguity
func TestMergeAggregations(t *testing.T)  {
	type testRecord struct {
		input []Aggregation;
		expected Aggregation
	}

	testData := []testRecord{
		{
			input: []Aggregation{{"hi": GetMax}, {"hi": GetMin, "yoo": GetRange}, {"hi": GetSum, "an": GetRange}, {"an": GetMin}},
			expected: Aggregation{
				"hi": GetSum,
				"yoo": GetRange,
				"an": GetMin,
			},
		},		
	}

	sampleArray := []interface{}{2, 1, 45, 6}

	for _, tr := range testData {
		res := MergeAggregations(tr.input)

		for key, agg := range tr.expected {
			got := res[key](sampleArray)
			expected := agg(sampleArray)

			if got != expected {
				t.Fatalf("for key '%s', expected %v; got %v",key,  agg, res[key])
			}
		}
	}
}

func Benchmark_mergeAggregations(b *testing.B)  {
	input := []Aggregation{
		{"hi": GetMax}, 
		{"hi": GetMin, "yoo": GetRange},
		{"hi": GetSum, "an": GetRange},
		{"an": GetMin},
	}

	for i := 0; i < b.N; i++ {
		MergeAggregations(input)
	}
	
	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | Aggregation as map  		| 560.1 ns/op	     | 256 B/op	             | 2 allocs/op           |  x  	   |
}