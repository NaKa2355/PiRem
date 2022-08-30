package daemon

import (
	"pirem/irdevice"
	"pirem/server"
	"pirem/tx"

	"github.com/NaKa2355/ir"
)

type Daemon struct {
	Devices    map[string]irdevice.Device
	ErrHandler func(error)
}

func (d *Daemon) Init() {
	d.Devices = map[string]irdevice.Device{}
}

func (d *Daemon) AddDevice(name string, dev irdevice.Device) {
	d.Devices[name] = dev
}

func (d Daemon) sendIRHandler(devName string, rawData ir.RawData) error {
	dev, exist := d.Devices[devName]
	if !exist {
		return ErrDevNotFound
	}
	respChan := make(chan tx.ResultResp)
	req := tx.SendIRReq{RespChan: respChan, Param: rawData}
	dev.SendReq(req)
	resp := <-respChan
	return resp.Err
}

func (d Daemon) receiveHandler(devName string) (ir.RawData, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return nil, ErrDevNotFound
	}

	respChan := make(chan tx.ResultIRRawDataResp)
	req := tx.ReceiveIRReq{RespChan: respChan}
	dev.SendReq(req)
	resp := <-respChan
	return resp.Value, resp.Err
}

func (d Daemon) getDevicesHandler() (map[string]irdevice.Device, error) {
	return d.Devices, nil
}

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return dev, ErrDevNotFound
	}
	return dev, nil
}

func (d Daemon) Start(server_port uint16) {
	handler := server.ServerHandlers{
		SendIRHandler: d.sendIRHandler,
		ReceiveIRData: d.receiveHandler,
		GetDevices:    d.getDevicesHandler,
		GetDevice:     d.getDeviceHandler,
		ErrHandler: func(err error) {
			d.ErrHandler(err)
		},
	}
	daemonServer := server.DaemonServer{}
	daemonServer.New(handler, server_port)
	daemonServer.Start()
}
