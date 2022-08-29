package reqjson

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

func JsonToIRRawData(jsonData []byte) (ir.RawData, error) {
	irData := IRData{}
	if err := json.Unmarshal(jsonData, &irData); err != nil {
		return nil, err
	}
	if irData.Type == "raw" {
		return unmarshalIRRawData(irData.IRData)
	} else {
		return nil, fmt.Errorf("type: \"%s\" is not supported: %s", irData.Type, ErrInvaildInput)
	}
}

func unmarshalIRRawData(jsonData json.RawMessage) (ir.RawData, error) {
	rawData := IRRawData{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, err
	}

	return strRawDataToRawData(rawData.Data)
}

func strRawDataToRawData(strRawData []string) (ir.RawData, error) {
	if len(strRawData)%2 == 0 {
		return nil, fmt.Errorf("number of pluses is must be odd: %s", ErrInvaildInput)
	}
	pulses := make([]ir.Pulse, len(strRawData))
	var pulseWidth int16 = 0
	var err error
	for i, strPulse := range strRawData {
		if _, err = fmt.Sscanf(strPulse, "%dms", &pulseWidth); err == nil {
			pulses[i].Prefix = ir.Milli
		} else if _, err = fmt.Sscanf(strPulse, "%dus", &pulseWidth); err == nil {
			pulses[i].Prefix = ir.Micro
		} else {
			return pulses, fmt.Errorf("use \"ms\" or \"us\" as prefix like \"10ms\": %s", ErrInvaildInput)
		}
		if pulseWidth < 0 {
			return pulses, fmt.Errorf("a pulse width must be positive interger: %s", ErrInvaildInput)
		}
		pulses[i].Width = pulseWidth
	}
	return pulses, nil
}
