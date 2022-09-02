package irdevice

import (
	"encoding/json"
	"pirem/irdata"
	"pirem/irdevice/tx"
	"time"

	"github.com/NaKa2355/irdevctrl"
)

type EventDispatcher struct {
	dev irdevctrl.Controller
}

func (eventQueue *EventDispatcher) Init(pluginPath string, jsonDevConf json.RawMessage) error {
	return nil
}

func (eventQueue *EventDispatcher) InitMock(dev irdevctrl.Controller) {
	eventQueue.dev = dev
}

func (eventQueue EventDispatcher) handleReceiveIRReq(req tx.ReceiveIRReq) {
	rawData, err := eventQueue.dev.ReceiveIRData()

	irData := irdata.Data{Type: irdata.Raw, IRData: rawData}

	resp := tx.ResultIRDataResp{Value: irData, Err: err}
	req.RespChan <- resp

	close(req.RespChan)
}

func (eventQueue EventDispatcher) handleSendIRReq(req tx.SendIRReq) {
	var err error

	rawData, err := req.Param.IRData.ConvertToRawData()
	if err == nil {
		err = eventQueue.dev.SendIRData(rawData)
	}

	time.Sleep(130 * time.Millisecond)

	resp := tx.ResultResp{Err: err}
	req.RespChan <- resp

	close(req.RespChan)
}

func (eventQueue *EventDispatcher) handleRemoveDevReq(req tx.RemoveDevReq) {
	err := eventQueue.dev.Drop()

	resp := tx.ResultResp{Err: err}

	req.RespChan <- resp

	close(req.RespChan)
}

func (eventQueue EventDispatcher) handleReq(req tx.Request) {
	req.Match(tx.ReqCases{
		ReceiveIR: func(value tx.ReceiveIRReq) {
			eventQueue.handleReceiveIRReq(value)
		},

		SendIR: func(value tx.SendIRReq) {
			eventQueue.handleSendIRReq(value)
		},

		RemoveDev: func(value tx.RemoveDevReq) {
			eventQueue.handleRemoveDevReq(value)
		},
	})
}

func (eventQueue EventDispatcher) GetBufferSize() uint16 {
	return eventQueue.dev.GetBufferSize()
}

func (eventQueue EventDispatcher) GetFeatures() irdevctrl.Features {
	return eventQueue.dev.GetSupportingFeatures()
}

func (eventQueue EventDispatcher) Start(reqChan <-chan tx.Request) {
	for {
		req, ok := <-reqChan
		if !ok {
			break
		}

		eventQueue.handleReq(req)
	}
}
