package main

import "pirem/tx"

type Device struct {
	pluginPath string
	buffSize   uint16
	reqchan    chan<- tx.Request
}

func NewDev(pluginPath string, buffSize uint16, reqchan chan tx.Request) Device {
	dev := Device{}
	dev.pluginPath = pluginPath
	dev.buffSize = buffSize
	dev.reqchan = reqchan
	return dev
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

func (dev Device) RemoveDev() tx.Responce {
	respChan := make(<-chan tx.ResultResp)
	req := tx.RemoveDevReq{RespChan: respChan}
	dev.reqchan <- req
	return req.RecvResp()
}
