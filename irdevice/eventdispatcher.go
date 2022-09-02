package irdevice

import (
	"encoding/json"
	"pirem/irdata"
	"pirem/irdevice/tx"
	"pirem/message"
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

func (eventDispatcher EventDispatcher) handleReceiveIRReq(m message.Message) {
	rawData, err := eventDispatcher.dev.ReceiveIRData()

	irData := irdata.Data{Type: irdata.Raw, IRData: rawData}

	resp := tx.NewRecvIRResp(irData, err)

	m.SendBack(message.NewOneWay(resp))
}

func (eventDispatcher EventDispatcher) handleSendIRReq(m message.Message) {
	var err error

	rawData, err := m.GetValue().(tx.SendIRReq).GetValue().IRData.ConvertToRawData()

	if err == nil {
		err = eventDispatcher.dev.SendIRData(rawData)
	}

	time.Sleep(130 * time.Millisecond)

	resp := tx.NewSendIRResp(err)
	m.SendBack(message.NewOneWay(resp))
}

func (eventDispatcher EventDispatcher) handleReq(m message.Message) {
	switch m.GetValue().(type) {
	case tx.SendIRReq:
		eventDispatcher.handleSendIRReq(m)
	case tx.RecvIRReq:
		eventDispatcher.handleReceiveIRReq(m)
	default:
		return
	}
}

func (eventQueue EventDispatcher) GetBufferSize() uint16 {
	return eventQueue.dev.GetBufferSize()
}

func (eventQueue EventDispatcher) GetFeatures() irdevctrl.Features {
	return eventQueue.dev.GetSupportingFeatures()
}

func (eventQueue EventDispatcher) Start(reqChan <-chan message.Message) {
	for {
		req, ok := <-reqChan
		if !ok {
			break
		}
		eventQueue.handleReq(req)
	}
}
