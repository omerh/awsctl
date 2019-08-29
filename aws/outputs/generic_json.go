package outputs

import (
	"encoding/json"
	"fmt"
)

type jsonOutput struct {
	Payload interface{} `json:"payload"`
	Region  string      `json:"region,omitempty"`
}

//PrintGenericJSONOutput appending region to result and print it to stdout
//
func PrintGenericJSONOutput(v interface{}, region string) {
	output := jsonOutput{Payload: v, Region: region}
	jsonout, _ := json.Marshal(output)
	fmt.Println(string(jsonout))
}
