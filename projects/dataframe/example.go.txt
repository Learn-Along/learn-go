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