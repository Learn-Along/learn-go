/*
the dataframe package offers a data structure that makes it easy to manipulate a collection of records.

The operations that can be done on Dataframe include:

- Select
- Filter (using Where)
- Sortby
- Groupby

TODO:

- Limit
- Offset

*/
package dataframe

import "github.com/learn-along/learn-go/projects/dataframe/internal"

/*
The Dataframe data structure that acts like a database, allowing selection, filtering, grouping and sorting.

- Planned: Add pagination
*/
type Dataframe internal.Dataframe

const (
	// Sortby 

	/*
	Enum Value for sorting in Ascending order
	*/
	ASC = internal.ASC
	/*
	Enum Value for sorting in Descending order
	*/
	DESC = internal.DESC
)

var (
	// Initialization

	/*
	Constructs a Dataframe from an array of maps and returns a pointer to it
	*/
	FromArray = internal.FromArray
	/*
	Constructs a Dataframe from an array of maps and returns a pointer to it
	*/
	FromMap = internal.FromMap

	// Filter 

	/*
	Combines a list of filters to produce a combined AND logical filter
	*/
	AND = internal.AND
	/*
	Combines a list of filters to produce a combined OR logical filter
	*/
	OR = internal.OR 
	/*
	Combines a list of filters to produce a combined NOT logical filter
	*/
	NOT = internal.NOT

	// GroupBy

	/*
	Aggregation function to get the sum of a given field in each group after a Groupby is done
	*/
	SUM = internal.SUM
	/*
	Aggregation function to get the maximum value of a given field in each group after a Groupby is done
	*/
	MAX = internal.MAX 
	/*
	Aggregation function to get the minimum value of a given field in each group after a Groupby is done
	*/
	MIN = internal.MIN 
	/*
	Aggregation function to get the mean value of a given field in each group after a Groupby is done
	*/
	MEAN = internal.MEAN
	/*
	Aggregation function to get the difference between maximum and minimum values 
	of a given field in each group after a Groupby is done
	*/
	RANGE = internal.RANGE
	/*
	Aggregation function to get the number of records in each group after a Groupby is done
	*/
	COUNT = internal.COUNT
)