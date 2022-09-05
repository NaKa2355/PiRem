package irdevice

/*
デーモンでデバイスを管理するための構造体を定義
SendReq関数でデバイスに紐づいたイベントループへリクエストを送る
モックアップはInitMockから作成、実際のデバイスはjson.Unmarshal()で作成する
*/

import (
	"encoding/json"
	"pirem/irdata"
	"pirem/irdevice/tx"
	"pirem/message"
	"time"

	"github.com/NaKa2355/irdevctrl"
)

type Devices map[string]Device

type Device struct {
	pluginPath      string
	buffSize        uint16
	timeout         time.Duration
	featurs         irdevctrl.Features
	reqChan         chan<- message.Message
	eventDispatcher EventDispatcher
}

func (dev *Device) InitMock(plugin_path string, timeout time.Duration, mock irdevctrl.Controller) {
	eventDispatcher := EventDispatcher{}
	dev.pluginPath = plugin_path
	eventDispatcher.InitMock(mock)
	dev.buffSize = eventDispatcher.GetBufferSize()
	dev.featurs = eventDispatcher.GetFeatures()
	dev.timeout = timeout
	dev.eventDispatcher = eventDispatcher
}

func (dev *Device) StartDispatcher() {
	reqChan := make(chan message.Message)
	dev.reqChan = reqChan

	go dev.eventDispatcher.Start(reqChan)
}

func (dev Device) GetPluginPath() string {
	return dev.pluginPath
}

func (dev Device) GetBuffSize() uint16 {
	return dev.buffSize
}

func (dev Device) GetFeatures() irdevctrl.Features {
	return dev.featurs
}

func (dev Device) SendIR(irdata irdata.Data) error {
	m := message.NewRoundTrip(tx.NewSendIRReq(irdata))
	dev.reqChan <- m

	resp, err := m.Receive()
	if err != nil {
		return err
	}

	return resp.GetValue().(tx.SendIRResp).GetValue()
}

func (dev Device) ReceiveIR() (irdata.Data, error) {
	m := message.NewRoundTrip(tx.RecvIRReq{})
	dev.reqChan <- m

	resp, err := m.Receive()
	if err != nil {
		return irdata.Data{}, err
	}

	return resp.GetValue().(tx.RecvIRResp).GetValue()
}

func (dev Device) Drop() {
	close(dev.reqChan)
}

func (dev *Device) UnmarshalJSON(data []byte) error {
	devicePrim := struct {
		PluginPath string          `json:"plugin_path"`
		Timeout    int             `json:"timeout"`
		DeviceConf json.RawMessage `json:"device_config"`
	}{}

	if err := json.Unmarshal(data, &devicePrim); err != nil {
		return err
	}

	eventDispatcher := EventDispatcher{}
	if err := eventDispatcher.Init(devicePrim.PluginPath, devicePrim.DeviceConf); err != nil {
		return err
	}

	dev.pluginPath = devicePrim.PluginPath
	dev.buffSize = eventDispatcher.GetBufferSize()
	dev.featurs = eventDispatcher.GetFeatures()
	dev.timeout = time.Duration(devicePrim.Timeout) * time.Second
	dev.eventDispatcher = eventDispatcher
	return nil
}

func (dev Device) MarshalJSON() ([]byte, error) {
	devicePrim := struct {
		PluginPath string             `json:"plugin_path"`
		BuffSize   uint16             `json:"buffsize"`
		Features   irdevctrl.Features `json:"features"`
	}{}

	devicePrim.BuffSize = dev.buffSize
	devicePrim.Features = dev.featurs
	devicePrim.PluginPath = dev.pluginPath

	return json.Marshal(devicePrim)
}
