package irdevice

import (
	"encoding/json"
	"pirem/irdevice/tx"

	"github.com/NaKa2355/irdevctrl"
)

type Devices map[string]Device

type Device struct {
	pluginPath string
	buffSize   uint16
	featurs    irdevctrl.Features
	reqChan    chan<- tx.Request
	eventQueue EventQueue
}

func (dev *Device) InitMock(mock irdevctrl.Controller) {
	eventQueue := EventQueue{}
	eventQueue.InitMock(mock)
	dev.buffSize = eventQueue.GetBufferSize()
	dev.featurs = eventQueue.GetFeatures()
	dev.eventQueue = eventQueue
}

func (dev *Device) GenerateEventQueue(queueLength int) {
	reqChan := make(chan tx.Request, queueLength)
	go dev.eventQueue.Start(reqChan)
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

func (dev Device) SendReq(req tx.Request) {
	dev.reqChan <- req
}

func (dev *Device) UnmarshalJSON(data []byte) error {
	devicePrim := struct {
		PuluginPath string             `json:"pulgin_path"`
		BuffSize    uint16             `json:"buffsize"`
		Features    irdevctrl.Features `json:"features"`
		DeviceConf  json.RawMessage    `json:"device_config"`
	}{}

	if err := json.Unmarshal(data, &devicePrim); err != nil {
		return err
	}

	eventQueue := EventQueue{}
	if err := eventQueue.Init(devicePrim.PuluginPath, devicePrim.DeviceConf); err != nil {
		return err
	}

	dev.pluginPath = devicePrim.PuluginPath
	dev.buffSize = devicePrim.BuffSize
	dev.featurs = devicePrim.Features
	dev.eventQueue = eventQueue
	return nil
}

func (dev Device) Marshal() ([]byte, error) {
	devicePrim := struct {
		PuluginPath string             `json:"pulgin_path"`
		BuffSize    uint16             `json:"buffsize"`
		Features    irdevctrl.Features `json:"features"`
	}{}

	devicePrim.BuffSize = dev.buffSize
	devicePrim.Features = dev.featurs
	devicePrim.PuluginPath = dev.pluginPath

	return json.Marshal(devicePrim)
}
