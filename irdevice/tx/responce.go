package tx

import (
	"pirem/irdata"
)

type ResponceType uint8

type ResultResp struct {
	Err error
}

type ResultIRDataResp struct {
	Value irdata.Data
	Err   error
}
