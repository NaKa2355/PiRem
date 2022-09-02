package daemon

import (
	"fmt"
	"pirem/defs"
	ir "pirem/irdata"
	"pirem/irdevice"
	"pirem/irdevice/tx"
	"pirem/server"
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

func (d Daemon) sendIRHandler(devName string, rawData ir.Data) error {
	dev, exist := d.Devices[devName]
	if !exist {
		return fmt.Errorf("the device, \"%s\" is not found: %s", devName, defs.ErrInvaildInput)
	}
	respChan := make(chan tx.ResultResp)
	req := tx.SendIRReq{RespChan: respChan, Param: rawData}
	dev.SendReq(req)
	resp := <-respChan
	return resp.Err
}

func (d Daemon) receiveHandler(devName string) (ir.Data, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return ir.Data{}, fmt.Errorf("the device, \"%s\" is not found: %s", devName, defs.ErrInvaildInput)
	}

	respChan := make(chan tx.ResultIRDataResp)
	req := tx.ReceiveIRReq{RespChan: respChan}
	dev.SendReq(req)
	resp, ok := <-respChan
	if !ok {
		return resp.Value, fmt.Errorf("no reply from the device: %s", defs.ErrInternal)
	}
	return resp.Value, resp.Err
}

func (d Daemon) getDevicesHandler() (map[string]irdevice.Device, error) {
	return d.Devices, nil
}

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return dev, fmt.Errorf("the device, \"%s\" is not found: %s", devName, defs.ErrInvaildInput)
	}
	return dev, nil
}

func (d Daemon) Stop() error {
	return nil
}

func (d Daemon) Start(server_port uint16) {
	handler := server.ServerHandlers{
		SendIRHandler:     d.sendIRHandler,
		RecvIRDataHandler: d.receiveHandler,
		GetDevicesHandler: d.getDevicesHandler,
		GetDeviceHandler:  d.getDeviceHandler,
		ErrHandler: func(err error) {
			d.ErrHandler(err)
		},
	}
	daemonServer := server.Server{}
	daemonServer.New(handler, server_port)
	daemonServer.Start()
}
