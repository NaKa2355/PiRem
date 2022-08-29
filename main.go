package main

import (
	"fmt"
	"pirem/reqjson"
	"pirem/respjson"
)

func main() {
	jsonStr := []byte(`{"type": "raw","ir_data": {"data": ["10ms", "20us", "10ms"]}}`)
	rawData, err := reqjson.JsonToIRRawData(jsonStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rawData)

	result, err := respjson.IRRawDataToJson(rawData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))

	rawData, err = reqjson.JsonToIRRawData(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rawData)
}
