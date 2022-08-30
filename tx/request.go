package tx

import "github.com/NaKa2355/ir"

type Request interface {
	Match(ReqCases)
}

type ReqCases struct {
	SendIR    func(SendIRReq)
	ReceiveIR func(ReceiveIRReq)
	RemoveDev func(RemoveDevReq)
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

type ReceiveIRReq struct {
	RespChan chan ResultIRRawDataResp
}

func (req ReceiveIRReq) Match(cases ReqCases) {
	if cases.ReceiveIR == nil {
		return
	}
	cases.ReceiveIR(req)
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
