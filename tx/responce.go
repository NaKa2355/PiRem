package tx

import "github.com/NaKa2355/ir"

type ResponceType uint8

type Responce interface {
	Match(RespCases)
}

type RespCases struct {
	Result          func(ResultResp)
	ResultUInt32    func(ResultUInt32Resp)
	ResultIRRawData func(ResultIRRawDataResp)
}

type ResultResp struct {
	Err error
}

func (resp ResultResp) Match(cases RespCases) {
	cases.Result(resp)
}

type ResultUInt32Resp struct {
	Value uint32
	Err   error
}

func (resp ResultUInt32Resp) Match(cases RespCases) {
	cases.ResultUInt32(resp)
}

type ResultIRRawDataResp struct {
	Value ir.RawData
	Err   error
}

func (resp ResultIRRawDataResp) Match(cases RespCases) {
	cases.ResultIRRawData(resp)
}
