package tx

import "github.com/NaKa2355/ir"

type Request interface {
	Match(ReqCases)
	RecvResp() Responce
}

type ReqCases struct {
	SendIR      func(SendIRReq)
	ReceiveIR   func(ReceiveIRReq)
	GetBuffSize func(GetBuffSizeReq)
	RemoveDev   func(RemoveDevReq)
}

type SendIRReq struct {
	RespChan chan ResultResp
	Param    ir.RawData
}

func (req SendIRReq) Match(cases ReqCases) {
	if cases.SendIR == nil {
		return
	}
	cases.SendIR(req)
}

func (req SendIRReq) RecvResp() Responce {
	return <-req.RespChan
}

type ReceiveIRReq struct {
	RespChan chan ResultIRRawDataResp
}

func (req ReceiveIRReq) Match(cases ReqCases) {
	if cases.ReceiveIR == nil {
		return
	}
	cases.ReceiveIR(req)
}

func (req ReceiveIRReq) RecvResp() Responce {
	return <-req.RespChan
}

type GetBuffSizeReq struct {
	RespChan chan ResultUInt32Resp
}

func (req GetBuffSizeReq) Match(cases ReqCases) {
	if cases.GetBuffSize == nil {
		return
	}
	cases.GetBuffSize(req)
}

func (req GetBuffSizeReq) RecvResp() Responce {
	return <-req.RespChan
}

type RemoveDevReq struct {
	RespChan chan ResultResp
}

func (req RemoveDevReq) Match(cases ReqCases) {
	if cases.RemoveDev == nil {
		return
	}
	cases.RemoveDev(req)
}

func (req RemoveDevReq) RecvResp() Responce {
	return <-req.RespChan
}
