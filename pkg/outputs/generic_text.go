package outputs

import (
	"fmt"
	"log"
)

//PrintGenericTextOutput iterate all k,v of a struct and print them
//
func PrintGenericTextOutput(t interface{}, region string, message string) {
	// fmt.Println(t)
	// v := reflect.ValueOf(t)
	// vt := reflect.TypeOf(t)

	log.Printf("Running on region: %v", region)

	if message != "" {
		log.Println(message)
	}

	if rec, ok := t.(map[string]interface{}); ok {
		for key, val := range rec {
			log.Printf("%s: %s", key, val)
		}
	} else {
		fmt.Println(t)
	}

	// if v.Kind() == reflect.Slice {

	// } else {
	// 	for i := 0; i < v.NumField(); i++ {
	// 		fmt.Printf("%v : %v", vt.Field(i).Name, v.Field(i))
	// 	}
	// }
}
