package irdevice

import (
	"encoding/json"
	"fmt"
	"pirem/irdata"
	"pirem/irdevice/tx"
	"pirem/message"
	"plugin"
	"time"

	"github.com/NaKa2355/irdevctrl"
)

type EventDispatcher struct {
	dev irdevctrl.Controller
}

func (eventDispatcher *EventDispatcher) Init(pluginPath string, jsonDevConf json.RawMessage) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("faild to open plugin: %s", err)
	}

	InitFromJson, err := p.Lookup("InitFromJson")
	if err != nil {
		return fmt.Errorf("loaded plugin does not have excepted function: %s", err)
	}

	dev, err := InitFromJson.(func(json.RawMessage) (irdevctrl.Controller, error))(jsonDevConf)
	if err != nil {
		return fmt.Errorf("loaded plugin does not have excepted function: %s", err)
	}

	eventDispatcher.dev = dev
	return nil
}

func (eventDispatcher *EventDispatcher) InitMock(dev irdevctrl.Controller) {
	eventDispatcher.dev = dev
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

func (eventDispatcher EventDispatcher) GetBufferSize() uint16 {
	return eventDispatcher.dev.GetBufferSize()
}

func (eventDispatcher EventDispatcher) GetFeatures() irdevctrl.Features {
	return eventDispatcher.dev.GetSupportingFeatures()
}

func (eventDispatcher EventDispatcher) Start(reqChan <-chan message.Message) {
	for {
		req, ok := <-reqChan
		if !ok {
			break
		}
		eventDispatcher.handleReq(req)
	}
}
