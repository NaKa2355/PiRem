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
	"sync"
	"time"

	"github.com/NaKa2355/irdevctrl"
)

type Devices map[string]*Device

type Device struct {
	mu           sync.RWMutex
	pluginPath   string
	buffSize     uint16
	timeout      time.Duration
	name         string
	featurs      irdevctrl.Features
	reqChan      chan<- message.Message
	deviceConfig json.RawMessage
}

func (dev *Device) InitAndSetupMock(plugin_path string, timeout time.Duration, mock irdevctrl.Controller) {
	eventDispatcher := EventDispatcher{}
	dev.pluginPath = plugin_path
	eventDispatcher.InitMock(mock)
	dev.buffSize = eventDispatcher.GetBufferSize()
	dev.featurs = eventDispatcher.GetFeatures()
	dev.timeout = timeout
	reqChan := make(chan message.Message)
	dev.reqChan = reqChan

	go eventDispatcher.Start(reqChan)
}

func (dev *Device) Setup() error {
	eventDispatcher := EventDispatcher{}
	if err := eventDispatcher.Init(dev.pluginPath, dev.deviceConfig); err != nil {
		return err
	}

	dev.mu.Lock()
	dev.buffSize = eventDispatcher.GetBufferSize()
	dev.featurs = eventDispatcher.GetFeatures()

	reqChan := make(chan message.Message)
	dev.reqChan = reqChan
	dev.mu.Unlock()

	go eventDispatcher.Start(reqChan)
	return nil
}

func (dev *Device) GetPluginPath() string {
	var pluginPath string
	dev.mu.RLock()
	pluginPath = dev.pluginPath
	dev.mu.RUnlock()
	return pluginPath
}

func (dev *Device) GetBuffSize() uint16 {
	var buffSize uint16
	dev.mu.RLock()
	buffSize = dev.buffSize
	dev.mu.RUnlock()
	return buffSize
}

func (dev *Device) GetFeatures() irdevctrl.Features {
	var features irdevctrl.Features
	dev.mu.RLock()
	features = dev.featurs
	dev.mu.RUnlock()
	return features
}

func (dev *Device) SendIR(irdata irdata.Data) error {
	m := message.NewRoundTrip(tx.NewSendIRReq(irdata))

	dev.mu.RLock()
	dev.reqChan <- m

	resp, err := m.Receive(dev.timeout)
	if err != nil {
		dev.mu.RLock()
		return err
	}

	dev.mu.RUnlock()
	return resp.GetValue().(tx.SendIRResp).GetValue()
}

func (dev *Device) ReceiveIR() (irdata.Data, error) {
	m := message.NewRoundTrip(tx.RecvIRReq{})

	dev.mu.RLock()
	dev.reqChan <- m

	resp, err := m.Receive(dev.timeout)
	if err != nil {
		dev.mu.RUnlock()
		return irdata.Data{}, err
	}

	dev.mu.RUnlock()
	return resp.GetValue().(tx.RecvIRResp).GetValue()
}

func (dev *Device) Drop() {
	dev.mu.Lock()
	close(dev.reqChan)
	dev.mu.Unlock()
}

func (dev *Device) UnmarshalJSON(data []byte) error {
	devicePrim := struct {
		PluginPath string          `json:"plugin_path"`
		Timeout    int             `json:"timeout"`
		DeviceConf json.RawMessage `json:"device_config"`
		Name       string          `json:"name"`
	}{}

	if err := json.Unmarshal(data, &devicePrim); err != nil {
		return err
	}

	dev.mu.Lock()
	dev.pluginPath = devicePrim.PluginPath
	dev.timeout = time.Duration(devicePrim.Timeout) * time.Second
	dev.deviceConfig = devicePrim.DeviceConf
	dev.name = devicePrim.Name
	dev.mu.Unlock()
	return nil
}

func (dev *Device) MarshalJSON() ([]byte, error) {
	devicePrim := struct {
		Name       string             `json:"name"`
		PluginPath string             `json:"plugin_path"`
		BuffSize   uint16             `json:"buffsize"`
		Features   irdevctrl.Features `json:"features"`
	}{}

	dev.mu.RLock()
	devicePrim.Name = dev.name
	devicePrim.BuffSize = dev.buffSize
	devicePrim.Features = dev.featurs
	devicePrim.PluginPath = dev.pluginPath
	dev.mu.RUnlock()

	return json.Marshal(devicePrim)
}
