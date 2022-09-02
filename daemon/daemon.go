package daemon

import (
	"fmt"
	"pirem/defs"
	"pirem/irdata"
	"pirem/irdevice"
	"pirem/server"
	"time"
)

type Daemon struct {
	Devices    map[string]irdevice.Device
	server     server.Server
	ErrHandler func(error)
}

func (d *Daemon) Init() {
	d.Devices = map[string]irdevice.Device{}
}

func (d *Daemon) AddDevice(name string, dev irdevice.Device) {
	d.Devices[name] = dev
}

func (d Daemon) sendIRHandler(devName string, irData irdata.Data) error {
	dev, exist := d.Devices[devName]
	if !exist {
		return fmt.Errorf("the device, \"%s\" not found: %s", devName, defs.ErrInvaildInput)
	}

	return dev.SendIR(irData)
}

func (d Daemon) receiveHandler(devName string) (irdata.Data, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return irdata.Data{}, fmt.Errorf("the device, \"%s\" not found: %s", devName, defs.ErrInvaildInput)
	}

	return dev.ReceiveIR()
}

func (d Daemon) getDevicesHandler() (map[string]irdevice.Device, error) {
	return d.Devices, nil
}

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.Devices[devName]
	if !exist {
		return dev, fmt.Errorf("the device, \"%s\" not found: %s", devName, defs.ErrInvaildInput)
	}
	return dev, nil
}

func (d Daemon) Stop() error {
	if err := d.server.Stop(15 * time.Second); err != nil {
		return err
	}

	for _, dev := range d.Devices {
		dev.Drop()
	}

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
