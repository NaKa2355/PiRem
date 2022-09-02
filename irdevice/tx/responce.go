package tx

import (
	"pirem/irdata"
)

type SendIRResp struct {
	err error
}

func NewSendIRResp(err error) SendIRResp {
	return SendIRResp{err: err}
}

func (resp SendIRResp) GetValue() error {
	return resp.err
}

type RecvIRResp struct {
	value irdata.Data
	err   error
}

func NewRecvIRResp(value irdata.Data, err error) RecvIRResp {
	return RecvIRResp{value: value, err: err}
}

func (resp RecvIRResp) GetValue() (irdata.Data, error) {
	return resp.value, resp.err
}
