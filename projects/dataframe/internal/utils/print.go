package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// Pretty prints JSON data
func PrettyPrintJSON(data []byte) error {
	var prettyJSON bytes.Buffer
    err := json.Indent(&prettyJSON, data, "", "\t")
    if err != nil {
        log.Println("JSON parse error: ", err)
		return err
    }
	fmt.Println(prettyJSON.String())
	return nil
}