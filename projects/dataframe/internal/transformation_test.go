package internal

import "testing"

var (
	_addTen = func (i interface{}) interface{}  {return i.(int) + 10 }
	_minusTen = func (i interface{}) interface{}  {return i.(int) - 10 }
	_multiplyByTen = func (i interface{}) interface{}  {return i.(int) * 10 }
	_divideByTen = func (i interface{}) interface{}  {return i.(int) / 10 }
)

// MergeTransformations should merge an Transformation list into a map of slices of RowWiseFunc functions
func TestMergeTransformations(t *testing.T)  {
	type testRecord struct {
		input []Transformation;
		expected map[string][]RowWiseFunc
	}

	sampleValue := 20

	testData := []testRecord{
		{
			input: []Transformation{
				{k: "hi", v: _addTen}, 
				{k: "hi", v: _minusTen},
				{k: "yoo", v: _multiplyByTen}, 
				{k: "hi", v: _divideByTen},
				{k: "an", v: _multiplyByTen}, 
				{k: "an", v: _minusTen},
			},
			expected: map[string][]RowWiseFunc{
				"hi": {_addTen, _minusTen, _divideByTen},
				"yoo": {_multiplyByTen},
				"an": {_multiplyByTen, _minusTen},
			},
		},		
	}


	for _, tr := range testData {
		res := MergeTransformations(tr.input)

		for key, v := range tr.expected {
			for i, agg := range v {
				got := res[key][i](sampleValue)
				expected := agg(sampleValue)

				if got != expected {
					t.Fatalf("for key '%s', expected %v; got %v",key,  agg, res[key][i])
				}
			}
		}
	}
}

/*
* Benchmark tests
*/

func Benchmark_mergeTransformations(b *testing.B) {
	input := []Transformation{
		{k: "hi", v: _addTen}, 
		{k: "hi", v: _minusTen},
		{k: "yoo", v: _multiplyByTen}, 
		{k: "hi", v: _divideByTen},
		{k: "an", v: _multiplyByTen}, 
		{k: "an", v: _minusTen},
	}

	for i := 0; i < b.N; i++ {
		MergeTransformations(input)
	}

	// Results:
	// ========
	// 
	// | Change 					| time				 | memory 				 | allocations			 | Choice  |
	// |----------------------------|--------------------|-----------------------|-----------------------|---------|
	// | Transformation as map  	| 1073 ns/op         | 488 B/op              | 8 allocs/op           |    	   |
	// | Transformation as struct  	| 766.5 ns/op        | 488 B/op              | 8 allocs/op           | x	   |
	// | Add goroutine in for loop	| 3689 ns/op	     | 632 B/op	      		 | 10 allocs/op   		 |		   |
}