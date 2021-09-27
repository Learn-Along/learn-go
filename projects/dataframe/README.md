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

## Quick Start

- Install [Go +1.17.1](https://golang.org/dl/). Installation instructions are [here](https://golang.org/doc/install)

- Start a go project of your choice e.g. 'df_example'

  ```sh
  # create the df_example directory and enter it
  mkdir df_example && cd df_example
  # initialize the go project using go modules
  go mod init github.com/<your username>/df_example
  # open the folder in Visual studio code
  code .
  ```

- Create your first go file in the project, say 'main.go'

  ```sh
  code main.go
  ```

- Install this package

  ```sh
  go get github.com/learn-along/learn-go/projects/dataframe
  ```

- Add the import to your 'main.go' file and add some fancy code e.g.

  ```go
  package main

  import (
    "encoding/json"
    "fmt"
    "log"
    "regexp"

    "github.com/learn-along/learn-go/projects/dataframe/types"
    "github.com/learn-along/learn-go/projects/dataframe/utils"
  )

  var dataArray = []map[string]interface{}{
      {"first name": "John", "last name": "Doe", "age": 30, "location": "Kampala" },
      {"first name": "Jane", "last name": "Doe", "age": 50, "location": "Lusaka" },
      {"first name": "Paul", "last name": "Doe", "age": 19, "location": "Kampala" },
      {"first name": "Richard", "last name": "Roe", "age": 34, "location": "Nairobi" },
      {"first name": "Reyna", "last name": "Roe", "age": 45, "location": "Nairobi" },
      {"first name": "Ruth", "last name": "Roe", "age": 60, "location": "Kampala" },
  }

  func main() {
    pkFields := []string{"first name", "last name"}
      df, err := types.FromArray(dataArray, pkFields)
    if err != nil {
      log.Fatal("initialiation: ", err)
    }

    // Prints the records in the dataframe if they were converted into JSON
    fmt.Printf("\ndf:\n")
    err = df.PrettyPrintRecords()
    if err != nil {
      log.Fatal("pretty print (Df): ", err)
    }

      txData, err := df.Select().Apply(
      df.Col("first name").Tx(
        func (v interface{}) interface{} {return fmt.Sprintf("Edit-%v", v)},
      ),
    ).Execute()
    if err != nil {
      log.Fatal("select: ", err)
    }

    txJSON, err := json.Marshal(txData)
    if err != nil {
      log.Fatal("JSON: ", err)
    }

    // Pretty prints txData
    fmt.Printf("\ntxData:\n")
    err = utils.PrettyPrintJSON(txJSON)
    if err != nil {
      log.Fatal("pretty print (txData): ", err)
    }

      df.Insert(txData)

      // prints the original dataArray plus txData appended to the bottom
    fmt.Printf("\nmerged Df:\n")
    err = df.PrettyPrintRecords()
    if err != nil {
      log.Fatal("pretty print (merged Df): ", err)
    }

      // delete all the records you added as txData
      df.Delete(df.Col("first name").IsLike(regexp.MustCompile("^Edit")))

      // prints dataArray without txData
    fmt.Printf("\nstripped Df:\n")
    err = df.PrettyPrintRecords()
    if err != nil {
      log.Fatal("pretty print (stripped Df): ", err)
    }

      // ... etc....there is Update, Merge, Copy, FromMap etc.
  }
  ```

- Run your go app

  ```sh
  go run main.go
  ```

_Side note: Or...you can just clone this repo and run the [example.go](./example.go) file as `go run example.go`_

## How to Test

- Install [Go +1.17.1](https://golang.org/dl/). Installation instructions are [here](https://golang.org/doc/install)

- Clone the repo and enter the dataframe folder

  ```sh
  git clone git@github.com:Learn-Along/learn-go.git && cd learn-go/projects/dataframe
  ```

- Test the code

  ```sh
  go test ./... -race -timeout=30s
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
data, err = df1.Select("age", "name", "date").GroupBy("name", "date").Agg(
                df1.Col("vote").Agg(MAX),
                df1.Col("age").Agg(MIN),
                // Or supply a custom aggregregate func that returns a single value given an array of values
                df1.Col("address").Agg(func(arr []interface{}) {return arr[0]}),
            ).Execute()

// filter
data, err = df1.Select("age", "name", "date").Where(
                        AND(
                          OR(
                            df1.Col("age").GreaterThan(20),
                            df1.Col("name").IsLike(regexp.MustCompile("^(?i)john$")),
                            df1.Col("source").Equals("https://sopherapps.com"),
                          ),
                        ),
                      ).Execute()

// pipe the operations one after another
// the ... should be replace with appropriate arguments of course.
data, err = df1.Select(...).Where(...).GroupBy(...).SortBy(...).Apply(...).Execute()
```

## Opportunities

This library has a lot of opportunity to improve. Some include:

- [ ] There is a lot of reallocation due to saving and returning interface{}, thus the GC slows down the entire app trying to follow up on each allocation
- [ ] Actually make use of the columnar structure in filtering, selecting, grouping etc. [Right now it seems to make things just more complicated than rowise structure as most if not all these operations are done in row-wise manner]
- [ ] Optimize memory usage. There is a lot of copying and heavy use of suboptimal type methods
- [ ] Optimize speed. There are way too many loops.
- [ ] Take advantage of concurrent design. No goroutines were used. These might help improve speed or readability or both.
- [ ] Clean up the code. Some of the functions arre really dirty. The tests also are quite dirty.
- [ ] Add benchmark tests
- [ ] Add the features for accessing any set of records using their indices (x and y) like pandas does it or just add `Limit(int)` and `Skip(int)` methods. I actually feel `Limit` and `Skip` would align better with the current API which mimics SQL.
- [ ] More stuff you will find. Just create an issue. If I am taking too long to respond (as I have been known to in the past), shoot me an email at [sales@sopherapps.com](mailto:sales@sopherapps.com)
- [ ] Saving the data as interface{} caused a lot of reallocations and poor performance. There was need to convert the data straight into given data tyes to make the operations faster

## License

Copyright (c) 2021 [Martin Ahindura](https://github.com/Tinitto) Licensed under the [MIT License](./LICENSE)
