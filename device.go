package main

import "pirem/tx"

type Device struct {
	pluginPath string
	buffSize   uint16
	reqchan    chan<- tx.Request
}

func (dev *Device) Init(pluginPath string, buffSize uint16, reqchan chan tx.Request) {
	dev.pluginPath = pluginPath
	dev.buffSize = buffSize
	dev.reqchan = reqchan
}

func (dev Device) GetPluginPath() string {
	return dev.pluginPath
}

func (dev Device) GetBuffSize() uint16 {
	return dev.buffSize
}

func (dev Device) SendReq(req tx.Request) tx.Responce {
	dev.reqchan <- req
	return req.RecvResp()
}
