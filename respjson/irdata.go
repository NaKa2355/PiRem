package respjson

import (
	"encoding/json"
	"fmt"

	"github.com/NaKa2355/ir"
)

type IRData struct {
	Type   string          `json:"type"`
	IRData json.RawMessage `json:"ir_data"`
}

type IRRawData struct {
	Data []string `json:"data"`
}

func IRRawDataToJson(data ir.RawData) ([]byte, error) {
	rawMessage, err := rawDataToJsonArray(data)
	if err != nil {
		return nil, err
	}
	irData := IRData{Type: "raw", IRData: rawMessage}
	return json.Marshal(irData)
}

func rawDataToJsonArray(data ir.RawData) (json.RawMessage, error) {
	rawData := IRRawData{make([]string, len(data))}

	for i, pluse := range data {
		if pluse.Prefix == ir.Micro {
			rawData.Data[i] = fmt.Sprintf("%dus", pluse.Width)
		} else if pluse.Prefix == ir.Milli {
			rawData.Data[i] = fmt.Sprintf("%dms", pluse.Width)
		}
	}

	return json.Marshal(rawData)
}
