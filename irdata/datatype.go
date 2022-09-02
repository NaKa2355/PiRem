package irdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"pirem/defs"
)

type DataType int

const (
	Raw DataType = iota
)

func (datatype *DataType) UnmarshalJSON(data []byte) error {
	var strDataType string
	if err := json.Unmarshal(data, &strDataType); err != nil {
		return err
	}
	switch strDataType {
	case "raw":
		*datatype = Raw
	default:
		return fmt.Errorf("type, \"%s\" is not supported: %s", strDataType, defs.ErrInvaildInput)
	}
	return nil
}

func (dataType DataType) MarshalJSON() ([]byte, error) {
	var dataTypePrim string

	switch dataType {
	case Raw:
		dataTypePrim = "raw"
	default:
		return []byte(dataTypePrim), errors.New("unsupported type")
	}

	return json.Marshal(dataTypePrim)
}
