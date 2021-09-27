package types

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/learn-along/learn-go/projects/dataframe/utils"
)

// insert for columns should fill any gaps in keys and Items with "", nil respectively
func TestColumn_insert(t *testing.T)  {
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{0: "hi", 1: "wow"}}
	col.insert(4, "yeah")
	expectedItems := []interface{}{"hi", "wow", nil, nil, "yeah"}

	if !utils.AreSliceEqual(expectedItems, col.Items()) {
		t.Fatalf("items expected: %v, got %v", expectedItems, col.Items())
	}
}

// GreaterThan should return a slice of booleans where true is for values greater than the value,
// false is for otherwise
func TestColumn_GreaterThan(t *testing.T)  {
	operand := 2
	col := Column{Name: "hi", Dtype: ObjectType, items: map[int]interface{}{
		0: 23.4, 1: 10, 2: -2, 3: 69, 4: 0.23, 5: 67}}
	expected := filterType{true, true, false, true, false, true}
	output := col.GreaterThan(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkColumn_GreaterThan(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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
func TestColumn_GreaterOrEquals(t *testing.T)  {
	operand := 2
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{
		0: 23.4, 1: 10, 2: 2, 3: 69, 4: 0.23, 5: 67}}
	expected := filterType{true, true, true, true, false, true}
	output := col.GreaterOrEquals(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkColumn_GreaterOrEquals(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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


// // XGreaterOrEquals should return a slice of booleans where true is for values greater or equal to the value,
// // false is for otherwise
// func TestColumn_XGreaterOrEquals(t *testing.T)  {
// 	operand := 2
// 	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{
// 		0: 23.4, 1: 10, 2: 2, 3: 69, 4: 0.23, 5: 67}}
// 	expected := filterType{true, true, true, true, false, true}
// 	output := col.XGreaterOrEquals(float64(operand))

// 	for i := 0; i < 6; i++ {
// 		if output[i] != expected[i] {
// 			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
// 		}
// 	}
// }

// func BenchmarkColumn_XGreaterOrEquals(b *testing.B)  {
// 	items := map[int]interface{}{}
// 	numberOfItems := 9000000

// 	for i := 0; i < numberOfItems; i++ {
// 		items[i] = i
// 	}

// 	col := Column{Name: "hi", Dtype: ObjectType, items: items}

// 	for i := 0; i < b.N; i++ {
// 		col.XGreaterOrEquals(1000)
// 	}

// 	// Results:
// 	// ========
// 	// benchtime=10s
// 	// 
// 	// | Change 						| time				 | memory 				 | allocations			 | Choice  |
// 	// |--------------------------------|--------------------|-----------------------|-----------------------|---------|
// 	// | NonX				    		| 871,616,373 ns/op	 | 97,562,900 B/op	     | 775,526 allocs/op     |  	   |
// 	// | X								| 6,315,185,648 ns/op| 540,413,944 B/op	 	 | 4,653,417 allocs/op   | x       |
// 	// portion size = 4, 9,536,656,396 ns/op	363,305,816 B/op	 3,102,418 allocs/op
// 	// portion size = 20, 6,379,516,892 ns/op	221,575,688 B/op	 1,861,412 allocs/op
// 	// portion size = 36, 7,710,833,247 ns/op	363,282,429 B/op	 3,102,308 allocs/op
// 	// portion size = 40, 4,511,972,400 ns/op	221,576,900 B/op	 1,861,417 allocs/op
// 	// portion size = 42, 5359138934 ns/op	274704190 B/op	 2326693 allocs/op
// 	// portion size = 45, 5666726788 ns/op	274710688 B/op	 2326723 allocs/op
// 	// portion size = 50, 5,698,110,606 ns/op	274,658,768 B/op	 2,326,475 allocs/op
// 	// portion size = 400, 5,610,789,320 ns/op	274,702,634 B/op	 2,326,686 allocs/op
// 	// portion size = 4000, 7,120,493,393 ns/op	363,237,666 B/op	 3,102,092 allocs/op
// }


// LessThan should return a slice of booleans where true is for values less than the value,
// false is for otherwise
func TestColumn_LessThan(t *testing.T)  {
	operand := 2
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{
		0: 23.4, 1: 10, 2: -2, 3: 69, 4: 0.23, 5: 67}}
	expected := filterType{false, false, true, false, true, false}
	output := col.LessThan(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkColumn_LessThan(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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
func TestColumn_LessOrEquals(t *testing.T)  {
	operand := 2
	col := Column{Name: "hi", Dtype: StringType, items: map[int]interface{}{
		0: 23.4, 1: 10, 2: 2, 3: 69, 4: 0.23, 5: 67}}
	expected := filterType{false, false, true, false, true, false}
	output := col.LessOrEquals(float64(operand))

	for i := 0; i < 6; i++ {
		if output[i] != expected[i] {
			t.Fatalf("index %d: expected: %v, got %v", i, expected[i], output[i])
		}
	}
}

func BenchmarkColumn_LessOrEquals(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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
func TestColumn_Equals(t *testing.T)  {
	type testRecord struct {
		operand interface{};
		items orderedMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: "hi", 
			items: orderedMapType{0: 23.4, 1: "hi", 2: 2, 3: 69, 4: 0.23, 5: 67},
			expected: filterType{false, true, false, false, false, false},
		},
		{
			operand: -2, 
			items: orderedMapType{0: 23.4, 1: "hi", 2: -2, 3: 69, 4: -2, 5: 67},
			expected: filterType{false, false, true, false, true, false},
		},
		{
			operand: 0.23, 
			items: orderedMapType{0: 23.4, 1: "hi", 2: 2, 3: 69, 4: 0.23, 5: 67},
			expected: filterType{false, false, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := Column{Name: "hi", Dtype: ObjectType, items: tr.items}
		output := col.Equals(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkColumn_Equals(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = i
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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
func TestColumn_IsLike(t *testing.T)  {
	type testRecord struct {
		operand *regexp.Regexp;
		items orderedMapType;
		expected filterType
	}

	testData := []testRecord{
		{
			operand: regexp.MustCompile("(?i)^L"), 
			items: orderedMapType{0: "London", 1: "Loe", 2: "livingstone", 3: "69", 4: "Duhaga", 5: "Yoo"},
			expected: filterType{true, true, true, false, false, false},
		},
		{
			operand: regexp.MustCompile(`^\d`), 
			items: orderedMapType{0: "London", 1: "Loe", 2: "Livingstone", 3: "69", 4: "Duhaga", 5: "2Yoo"},
			expected: filterType{false, false, false, true, false, true},
		},
		{
			operand: regexp.MustCompile("^Duhaga"), 
			items: orderedMapType{0: "London", 1: "Loe", 2: "Livingstone", 3: "69", 4: "Duhaga", 5: "Yoo"},
			expected: filterType{false, false, false, false, true, false},
		},
	}

	for index, tr := range testData {
		col := Column{Name: "hi", Dtype: ObjectType, items: tr.items}
		output := col.IsLike(tr.operand)
	
		for i := 0; i < 6; i++ {
			if output[i] != tr.expected[i] {
				t.Fatalf("test record #%d, index %d: expected: %v, got %v", index, i, tr.expected[i], output[i])
			}
		}
	}

}

func BenchmarkColumn_IsLike(b *testing.B)  {
	items := map[int]interface{}{}
	numberOfItems := 9000000

	for i := 0; i < numberOfItems; i++ {
		items[i] = fmt.Sprintf("%v", i)
	}

	col := Column{Name: "hi", Dtype: ObjectType, items: items}

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





