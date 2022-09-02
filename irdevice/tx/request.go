package tx

import "pirem/irdata"

type SendIRReq struct {
	value irdata.Data
}

func NewSendIRReq(value irdata.Data) SendIRReq {
	req := SendIRReq{}
	req.value = value
	return req
}

func (req SendIRReq) GetValue() irdata.Data {
	return req.value
}

type RecvIRReq struct{}
