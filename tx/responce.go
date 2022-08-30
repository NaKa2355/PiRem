package tx

import "github.com/NaKa2355/ir"

type ResponceType uint8

type RespCases struct {
	Result          func(ResultResp)
	ResultIRRawData func(ResultIRRawDataResp)
}

type ResultResp struct {
	Err error
}

type ResultIRRawDataResp struct {
	Value ir.RawData
	Err   error
}
