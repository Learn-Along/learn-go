package utils

import (
	"encoding/json"
	"log"
)

// The PrettyPrintJSON function prints out JSON in a pretty format
func ExamplePrettyPrintJSON()  {
	type dataType struct {
		Foo string `json:"foo"`
		Jesus string `json:"Jesus"`
		Woo string `json:"Woo"`
	}

	data := dataType{
		Foo: "bar",
		Jesus: "the LORD",
		Woo: "hoo",
	}
	
	dataAsJson, err := json.Marshal(data)
	if err != nil {
		log.Fatal("error marshalling the data: ", err)
	}

	PrettyPrintJSON(dataAsJson)
	// Output:
	// { 
	//	"foo": "bar",
	//	"Jesus": "the LORD",
	//	"Woo": "hoo"
	// }
}