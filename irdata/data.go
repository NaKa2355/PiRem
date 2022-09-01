package irdata

import (
	"encoding/json"
	"errors"

	"github.com/NaKa2355/irdevctrl"
)

type Data struct {
	Type   DataType
	IRData irdevctrl.DataConverter
}

func (irdata *Data) UnmarshalJSON(data []byte) error {
	jsonIrData := struct {
		Type   DataType        `json:"type"`
		IRData json.RawMessage `json:"ir_data"`
	}{}

	if err := json.Unmarshal(data, &jsonIrData); err != nil {
		return err
	}
	switch jsonIrData.Type {
	case Raw:
		rawData := irdevctrl.RawData{}
		if err := json.Unmarshal(jsonIrData.IRData, &rawData); err != nil {
			return err
		}
		irdata.IRData = rawData
	default:
		return errors.New("unsupported type")
	}
	return nil
}

func (irdata Data) MarshalJSON() ([]byte, error) {
	var err error
	irDataPrim := struct {
		Type   DataType          `json:"type"`
		IRData irdevctrl.RawData `json:"ir_data"`
	}{}
	irDataPrim.Type = irdata.Type
	irDataPrim.IRData, err = irdata.IRData.ConvertToRawData()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(irDataPrim)
}
