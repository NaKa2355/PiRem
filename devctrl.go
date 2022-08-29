package main

import (
	"pirem/tx"
	"time"

	"github.com/NaKa2355/ir"
)

type DevController struct {
	dev ir.Device
}

func (devctrl *DevController) Init(dev ir.Device) {
	devctrl.dev = dev
}

func (devctrl DevController) handleReceiveIRReq(req tx.ReceiveIRReq) {
	rawData, err := devctrl.dev.ReceiveIRData()
	resp := tx.ResultIRRawDataResp{Value: rawData, Err: err}
	req.RespChan <- resp
	close(req.RespChan)
}

func (devctrl DevController) handleSendIRReq(req tx.SendIRReq) {
	err := devctrl.dev.SendIRData(req.Param)
	resp := tx.ResultResp{Err: err}
	time.Sleep(130 * time.Millisecond)
	req.RespChan <- resp
	close(req.RespChan)
}

func (devctrl *DevController) handleRemoveDevReq(req tx.RemoveDevReq) {
	err := devctrl.dev.Drop()
	devctrl.dev = nil
	resp := tx.ResultResp{Err: err}
	req.RespChan <- resp
	close(req.RespChan)
}

func (devctrl DevController) handleReq(req tx.Request) {
	req.Match(tx.ReqCases{
		ReceiveIR: func(value tx.ReceiveIRReq) {
			devctrl.handleReceiveIRReq(value)
		},
		SendIR: func(value tx.SendIRReq) {
			devctrl.handleSendIRReq(value)
		},
		RemoveDev: func(value tx.RemoveDevReq) {
			devctrl.handleRemoveDevReq(value)
		},
	})
}

func (devctrl DevController) Start(reqChan <-chan tx.Request) {
	for {
		if devctrl.dev == nil {
			break
		}
		req := <-reqChan
		devctrl.handleReq(req)
	}
}
