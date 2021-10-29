package internal

import (
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/internal/utils"
)

// Len returns the length of the map
func TestOrderedStringMap_Len(t *testing.T)  {
	type testRecord struct {
		input OrderedStringMapType;
		expected int
	}

	testData := []testRecord{
		{
			input: OrderedStringMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			expected: 4,
		},
		{
			input: OrderedStringMapType{0: "4"},
			expected: 1,
		},
		{
			input: OrderedStringMapType{0: "4", 2: "yoohoo", 3: "salut"},
			expected: 3,
		},
	}

	for _, tr := range testData {
		got := tr.input.Len()
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkStringMap_Len(b *testing.B)  {
	input := OrderedStringMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"}

	for i := 0; i < b.N; i++ {
		input.Len()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 0.4084 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// ToSlice converts an ordered map into a slice, ignoring gaps in indices automatically
func TestOrderedStringMap_ToSlice(t *testing.T)  {
	type testRecord struct {
		input OrderedStringMapType;
		expected []string
	}

	testData := []testRecord{
		{
			input: OrderedStringMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			expected: []string{"hi", "hello", "yoohoo", "salut"},
		},
		{
			input: OrderedStringMapType{0: "4", 1: "2", 2: "yoohoo", 3: "salut"},
			expected: []string{"4", "2", "yoohoo", "salut"},
		},
		{
			input: OrderedStringMapType{0: "4", 2: "yoohoo", 3: "salut"},
			expected: []string{"4", "yoohoo", "salut"},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice()
		if !utils.AreStringSliceEqual(got.([]string), tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkStringMap_ToSlice(b *testing.B)  {
	input := OrderedStringMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"}

	for i := 0; i < b.N; i++ {
		input.ToSlice()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 546.5 ns/op	     | 176 B/op	       		 | 5 allocs/op           |  x  	   |
}

// Defragmentize should reorder the ordered dict basing
// on the new indices passed to it
func TestOrderedStringMap_Defragmentize(t *testing.T)  {
	type testRecord struct {
		_map OrderedStringMapType;
		newOrder []int;
		expected []string
	}

	testData := []testRecord{
		{
			_map: OrderedStringMapType{0: "hi", 1: "hello", 2: "yoohoo", 3: "salut"},
			newOrder: []int{3, 1, 0},
			expected: []string{"salut", "hello", "hi"},
		},
		{
			_map: OrderedStringMapType{0: "4", 1: "2", 2: "yoohoo", 3: "salut"},
			newOrder: []int{1, 0, 3, 2},
			expected: []string{"2", "4", "salut", "yoohoo"},
		},
	}

	for _, tr := range testData {
		tr._map.Defragmentize(tr.newOrder)
		got := tr._map.ToSlice()

		if !utils.AreStringSliceEqual(got.([]string), tr.expected) {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkStringMap_Defragmentize(b *testing.B)  {
	input := OrderedStringMapType{0: "hi", 1: "hello", 3: "salut"}

	for i := 0; i < b.N; i++ {
		input.Defragmentize([]int{0, 3, 1})
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 292.3 ns/op	     | 256 B/op	       		 | 2 allocs/op           |  x  	   |
}

// Len returns the length of the map
func TestOrderedIntMap_Len(t *testing.T)  {
	type testRecord struct {
		input OrderedIntMapType;
		expected int
	}

	testData := []testRecord{
		{
			input: OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8},
			expected: 4,
		},
		{
			input: OrderedIntMapType{0: 4},
			expected: 1,
		},
		{
			input: OrderedIntMapType{0: 4, 2: 67, 3: 76},
			expected: 3,
		},
	}

	for _, tr := range testData {
		got := tr.input.Len()
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkIntMap_Len(b *testing.B)  {
	input := OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.Len()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 0.8400 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// ToSlice converts an ordered map into a slice, ignoring gaps in indices automatically
func TestOrderedIntMap_ToSlice(t *testing.T)  {
	type testRecord struct {
		input OrderedIntMapType;
		expected []int
	}

	testData := []testRecord{
		{
			input: OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8},
			expected: []int{0: 4, 1: 67, 2: 90, 3: 8},
		},
		{
			input: OrderedIntMapType{0: 4, 1: 67, 3: 8},
			expected: []int{4, 67, 8},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice().([]int)
		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkIntMap_ToSlice(b *testing.B)  {
	input := OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.ToSlice()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 427.4 ns/op	     | 144 B/op	       		 | 5 allocs/op           |  x  	   |
}

// Defragmentize should reorder the ordered dict basing
// on the new indices passed to it
func TestOrderedIntMap_Defragmentize(t *testing.T)  {
	type testRecord struct {
		_map OrderedIntMapType;
		newOrder []int;
		expected []int;
	}

	testData := []testRecord{
		{
			_map: OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8},
			newOrder: []int{3, 1, 0},
			expected: []int{8, 67, 4},
		},
		{
			_map: OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8},
			newOrder: []int{1, 0, 3, 2},
			expected: []int{67, 4, 8, 90},
		},
	}

	for _, tr := range testData {
		tr._map.Defragmentize(tr.newOrder)
		got := tr._map.ToSlice().([]int)

		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkIntMap_Defragmentize(b *testing.B)  {
	input := OrderedIntMapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.Defragmentize([]int{0, 1, 3})
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 210.1 ns/op	     | 192 B/op	       		 | 2 allocs/op           |  x  	   |
}

// Len returns the length of the map
func TestOrderedFloat64Map_Len(t *testing.T)  {
	type testRecord struct {
		input OrderedFloat64MapType;
		expected int
	}

	testData := []testRecord{
		{
			input: OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8},
			expected: 4,
		},
		{
			input: OrderedFloat64MapType{0: 4},
			expected: 1,
		},
		{
			input: OrderedFloat64MapType{0: 4, 2: 67, 3: 76},
			expected: 3,
		},
	}

	for _, tr := range testData {
		got := tr.input.Len()
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkFloat64Map_Len(b *testing.B)  {
	input := OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.Len()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 0.3927 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// ToSlice converts an ordered map into a slice, ignoring gaps in indices automatically
func TestOrderedFloat64Map_ToSlice(t *testing.T)  {
	type testRecord struct {
		input OrderedFloat64MapType;
		expected []float64
	}

	testData := []testRecord{
		{
			input: OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8},
			expected: []float64{0: 4, 1: 67, 2: 90, 3: 8},
		},
		{
			input: OrderedFloat64MapType{0: 4, 1: 67, 3: 8},
			expected: []float64{4, 67, 8},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice().([]float64)
		
		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkFloat64Map_ToSlice(b *testing.B)  {
	input := OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.ToSlice()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 631.9 ns/op	     | 144 B/op	       		 | 5 allocs/op           |  x  	   |
}

// Defragmentize should reorder the ordered dict basing
// on the new indices passed to it
func TestOrderedFloat64Map_Defragmentize(t *testing.T)  {
	type testRecord struct {
		_map OrderedFloat64MapType;
		newOrder []int;
		expected []float64
	}

	testData := []testRecord{
		{
			_map: OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8},
			newOrder: []int{3, 1, 0},
			expected: []float64{8, 67, 4},
		},
		{
			_map: OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8},
			newOrder: []int{1, 0, 3, 2},
			expected: []float64{67, 4, 8, 90},
		},
	}

	for _, tr := range testData {
		tr._map.Defragmentize(tr.newOrder)
		got := tr._map.ToSlice().([]float64)

		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkFloat64Map_Defragmentize(b *testing.B)  {
	input := OrderedFloat64MapType{0: 4, 1: 67, 2: 90, 3: 8}

	for i := 0; i < b.N; i++ {
		input.Defragmentize([]int{0, 3})
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 214.9 ns/op	     | 192 B/op	       		 | 2 allocs/op           |  x  	   |
}

// Len returns the length of the map
func TestOrderedBoolMap_Len(t *testing.T)  {
	type testRecord struct {
		input OrderedBoolMapType;
		expected int
	}

	testData := []testRecord{
		{
			input: OrderedBoolMapType{0: true, 1: true, 2: false, 3: true},
			expected: 4,
		},
		{
			input: OrderedBoolMapType{0: false},
			expected: 1,
		},
		{
			input: OrderedBoolMapType{0: true, 2: false, 3: true},
			expected: 3,
		},
	}

	for _, tr := range testData {
		got := tr.input.Len()
		if got != tr.expected {
			t.Fatalf("expected %v; got %v", tr.expected, got)
		}
	}
}

func BenchmarkBoolMap_Len(b *testing.B)  {
	input := OrderedBoolMapType{0: false, 1: true, 2: false, 3: true}

	for i := 0; i < b.N; i++ {
		input.Len()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 0.8741 ns/op	     | 0 B/op	       		 | 0 allocs/op           |  x  	   |
}

// ToSlice converts an ordered map into a slice, ignoring gaps in indices automatically
func TestOrderedBoolMap_ToSlice(t *testing.T)  {
	type testRecord struct {
		input OrderedBoolMapType;
		expected []bool
	}

	testData := []testRecord{
		{
			input: OrderedBoolMapType{0: false, 1: true, 2: true, 3: false},
			expected: []bool{false, true, true, false},
		},
		{
			input: OrderedBoolMapType{0: false, 1: true, 3: false},
			expected: []bool{false, true, false},
		},
	}

	for _, tr := range testData {
		got := tr.input.ToSlice().([]bool)

		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkBoolMap_ToSlice(b *testing.B)  {
	input := OrderedBoolMapType{0: false, 1: true, 2: false, 3: true}

	for i := 0; i < b.N; i++ {
		input.ToSlice()
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 418.6 ns/op	     | 116 B/op	       		 | 5 allocs/op           |  x  	   |
}

// Defragmentize should reorder the ordered dict basing
// on the new indices passed to it
func TestOrderedBoolMap_Defragmentize(t *testing.T)  {
	type testRecord struct {
		_map OrderedBoolMapType;
		newOrder []int;
		expected []bool
	}

	testData := []testRecord{
		{
			_map: OrderedBoolMapType{0: false, 1: true, 2: false, 3: false},
			newOrder: []int{3, 1, 0},
			expected: []bool{false, true, false},
		},
		{
			_map: OrderedBoolMapType{0: false, 1: true, 2: false, 3: false},
			newOrder: []int{1, 0, 3, 2},
			expected: []bool{true, false, false, false},
		},
	}

	for _, tr := range testData {
		tr._map.Defragmentize(tr.newOrder)
		got := tr._map.ToSlice().([]bool)

		if len(got) != len(tr.expected) {
			t.Fatalf("expected length: %d; got %d", len(got), len(tr.expected))
		}

		for i, v := range tr.expected {
			gotValue := got[i]

			if gotValue != v {
				t.Fatalf("expected %v; got %v", v, gotValue)
			}
		}
	}
}

func BenchmarkBoolMap_Defragmentize(b *testing.B)  {
	input := OrderedBoolMapType{0: false, 1: true, 2: false, 3: true}

	for i := 0; i < b.N; i++ {
		input.Defragmentize([]int{1, 0})
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | None				  		| 162.3 ns/op	     | 144 B/op	       		 | 2 allocs/op           |  x  	   |
}