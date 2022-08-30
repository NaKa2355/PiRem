package irdevice

import "pirem/tx"

type Device struct {
	PluginPath string `json:"plugin_path"`
	BuffSize   uint16 `json:"buf_size"`
	reqchan    chan<- tx.Request
}

func (dev *Device) Init(pluginPath string, buffSize uint16, reqchan chan tx.Request) {
	dev.PluginPath = pluginPath
	dev.BuffSize = buffSize
	dev.reqchan = reqchan
}

func (dev Device) GetPluginPath() string {
	return dev.PluginPath
}

func (dev Device) GetBuffSize() uint16 {
	return dev.BuffSize
}

func (dev Device) SendReq(req tx.Request) {
	dev.reqchan <- req
}
