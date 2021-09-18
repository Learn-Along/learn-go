# dataframe

dataframe is a data structure that makes analysis of lists of records easy, by allowing selecting fields, sorting, grouping, filtering and applying functions on each record.

## Purpose

Many times when we retrieve data, we wish to manipulate it in a way to give us information we can use.
We wish to, say, filter out what is unnecessary, group it by certain properties, find aggregates on each group e.g. mean, standard deviation etc.

However, most data is received in a form that can't allow such analytics. Essentially it comes as a list of records i.e. in a row-wise fashion e.g.

```JSON
[
  {"title": "Mr.", "author": "John Doe", "age": 38, "vote_result": "YES"},
  {"title": "Ms.", "author": "Jane Doe", "age": 20, "vote_result": "NO"},
  {"title": "Dr.", "author": "Richard Roe", "age": 38, "vote_result": "YES"}
]
```

Dataframe rearranges the data into a columnar data structure, that is easy for selecting fields, filtering by fields etc.

For example, the above data is turned into something like (if it were JSON):

```JSON
[
  {"label": "title", "items": ["Mr.", "Ms.", "Mr."]},
  {"label": "author", "items": ["John Doe", "Jane Doe", "Richard Roe"]},
  {"label": "age", "items": [38, 20, 54]},
  {"label": "vote_result", "items": ["YES", "NO", "YES"]}
]
```

## Dependencies

- [Go +1.17.1](https://golang.org/dl/)

## Development Environment

Operating system: [Ubuntu 18.04](https://releases.ubuntu.com/18.04.5/)
IDE: [Visual Studio Code](https://code.visualstudio.com/)
Go: [Go 1.17.1](https://golang.org/dl/)

## VS Code Extensions

[Go](https://marketplace.visualstudio.com/items?itemName=golang.Go)

## How to Run

- Install [Go +1.17.1](https://golang.org/dl/). Installation instructions are [here](https://golang.org/doc/install)

- Clone the repo and enter the dataframe folder

  ```sh
  git clone git@github.com:Learn-Along/learn-go.git && cd learn-go/projects/dataframe
  ```

- Run the example code

  ```sh
  go run example.go
  ```

## How to Test

- Install [Go +1.17.1](https://golang.org/dl/). Installation instructions are [here](https://golang.org/doc/install)

- Clone the repo and enter the dataframe folder

  ```sh
  git clone git@github.com:Learn-Along/learn-go.git && cd learn-go/projects/dataframe
  ```

- Test the code

  ```sh
  go test -race -timeout=30s
  ```

## Design

This is a description of the data structures:

### Dataframe

This is the data structure that exposes the public API. It is described as below:

```go
type Dataframe struct {}

/*
* Creation functions:
*/

// Create a dataframe from an array of maps, given an array of fields to use to uniquely id records
df1, err := FromArray(records, primaryFields)

// Create a dataframe from a map of maps, given an array of fields to use to uniquely id records
df2, err := FromMap(recordsMap, primaryFields)


/*
* Mutation Methods
*/

// Insert an array of maps as new records.
// It will overwrite any record whose primary fields are the same as those of new records
err = df1.Insert(moreRecords)

// Update any number of items that fulfill a given condition
err = df1.Update(AND(df1.Col("age").GreaterThan(13), df1.Col("name").IsLike(regexp.MustCompile("^john$"))), map[string]interface{}{"age": 20})

// Delete any number of items that fulfill a given condition
err = df1.Delete(AND(df1.Col("age").GreaterThan(3), df1.Col("name").IsLike(regexp.MustCompile("^john"))))

/*
* Selection methods
*/

// select a few fields, with apply to transformation the data accordingly and return the map of records
data, err = df1.Select("age", "name", "date").Apply(
    df1.Col("age").Tx(func(v int) {return v*8  }),
    df1.Col("name").Tx(func(v string) { return fmt.Sprintf("name is %s", v) }),
  ).Execute()

// sort
data, err = df1.Select("age", "name", "date").SortBy(
                df1.Col("age").Order(DESC),
                df1.Col("name").Order(ASC),
            ).Execute()

// groupby
data, err = df1.Select("age", "name", "date").GroupBy(
                df1.Col("age").Agg(MAX),
                df1.Col("date").Agg(MIN),
                // Or supply a custom aggregregate func that returns a single value given an array of values
                df1.Col("name").Agg(func(arr []interface{}) {return arr[0]}),
            ).Execute()

// filter
data, err = df1.Select("age", "name", "date").Where(
              AND(
                OR(
                  df1.Col("age").GreatherThan(20),
                  df1.Col("name").IsLike(regexp.MustCompile("^(?i)john$"))
                  ),
                  df1.Col("source").Equals("https://sopherapps.com")
              )
          ).Execute()

// pipe the operations one after another
// the ... should be replace with appropriate arguments of course.
data, err = df1.Select(...).Where(...).GroupBy(...).SortBy(...).Apply(...).Execute()
```

## License

Copyright (c) 2021 [Martin Ahindura](https://github.com/Tinitto) Licensed under the [MIT License](./LICENSE)
