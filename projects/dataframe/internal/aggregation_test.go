package internal

import "testing"

// MAX should return the maximum value (as a float64 value) in a given list of items
func TestMAX(t *testing.T)  {
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
		got := MAX(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkMAX(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		MAX(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 57.17 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// MIN should return the minimum value (as a float64 value) in a given list of items
func TestMIN(t *testing.T)  {
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
		got := MIN(tr.input)
		if got != tr.expected {
			t.Fatalf("loop %d: expected %v; got %v", loop, tr.expected, got)
		}
	}
}

func BenchmarkMIN(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		MIN(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 66.58 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// SUM should return the sum (as a float64 value) of the given list of items, but will return nil if 
// the values are not all numbers (nil values are treated as zero)
func TestSUM(t *testing.T)  {
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
		got := SUM(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkSUM(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		SUM(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 63.92 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// MEAN should return the mean (as a float64 value) of the given list of items.
// It returns nil if the values are not numbers (nil values are treated as zero)
func TestMEAN(t *testing.T)  {
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
		got := MEAN(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkMEAN(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		MEAN(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 85.74 ns/op	     | 40 B/op	       		 | 3 allocs/op           |  x  	   |
}

// COUNT should return the number of items in the list of values, including nils
func TestCOUNT(t *testing.T)  {
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
		got := COUNT(tr.input)
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkCOUNT(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		COUNT(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 42.17 ns/op	     | 24 B/op	       		 | 1 allocs/op           |  x  	   |
}

// RANGE should return the range (as a float64 value) of the given list of items
// i.e. maximum minus minimum.
// It returns nil if the values are not numbers (nil values are ignored)
func TestRANGE(t *testing.T)  {
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
		got := RANGE(tr.input)
		if got != tr.expected {
			t.Fatalf("loop %d: expected %v; got %v", loop, tr.expected, got)
		}
	}
}

func BenchmarkRANGE(b *testing.B)  {
	input := []float64{1.9, 2.3, 4.8, 0.8}

	for i := 0; i < b.N; i++ {
		RANGE(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None  						| 70.68 ns/op	     | 32 B/op	       		 | 2 allocs/op           |  x  	   |
}

// mergeAggregations should merge an aggregation list into a single aggregation 
// ensuring that the last aggregateFunc to be attached to a given column is the one kept,
// the previous ones are overwritten, to avoid ambiguity
func TestMergeAggregations(t *testing.T)  {
	type testRecord struct {
		input []aggregation;
		expected aggregation
	}

	testData := []testRecord{
		{
			input: []aggregation{{"hi": MAX}, {"hi": MIN, "yoo": RANGE}, {"hi": SUM, "an": RANGE}, {"an": MIN}},
			expected: aggregation{
				"hi": SUM,
				"yoo": RANGE,
				"an": MIN,
			},
		},		
	}

	sampleArray := []interface{}{2, 1, 45, 6}

	for _, tr := range testData {
		res := mergeAggregations(tr.input)

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
	input := []aggregation{
		{"hi": MAX}, 
		{"hi": MIN, "yoo": RANGE},
		{"hi": SUM, "an": RANGE},
		{"an": MIN},
	}

	for i := 0; i < b.N; i++ {
		mergeAggregations(input)
	}
	
	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | aggregation as map  		| 560.1 ns/op	     | 256 B/op	             | 2 allocs/op           |  x  	   |
}