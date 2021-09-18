package types

import "testing"

// mergeTransformations should merge an transformation list into a map of slices of rowWiseFunc functions
func TestMergeTransformations(t *testing.T)  {
	type testRecord struct {
		input []transformation;
		expected map[string][]rowWiseFunc
	}

	sampleValue := 20

	_addTen := func (i interface{}) interface{}  {return i.(int) + 10 }
	_minusTen := func (i interface{}) interface{}  {return i.(int) - 10 }
	_multiplyByTen := func (i interface{}) interface{}  {return i.(int) * 10 }
	_divideByTen := func (i interface{}) interface{}  {return i.(int) / 10 }

	testData := []testRecord{
		{
			input: []transformation{{"hi": _addTen}, {"hi": _minusTen, "yoo": _multiplyByTen}, {"hi": _divideByTen, "an": _multiplyByTen}, {"an": _minusTen}},
			expected: map[string][]rowWiseFunc{
				"hi": {_addTen, _minusTen, _divideByTen},
				"yoo": {_multiplyByTen},
				"an": {_multiplyByTen, _minusTen},
			},
		},		
	}


	for _, tr := range testData {
		res := mergeTransformations(tr.input)

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