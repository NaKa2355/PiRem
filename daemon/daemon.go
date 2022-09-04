package daemon

import (
	"net/http"
	"pirem/irdevice"
	"pirem/server"
	"time"
)

type Daemon struct {
	devices    map[string]irdevice.Device
	server     *server.Server
	errHandler func(error)
}

func (d *Daemon) AddDevice(name string, dev irdevice.Device) {
	d.devices[name] = dev
}

func sendError(err error, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}

func (d Daemon) Stop() error {
	if err := d.server.Stop(15 * time.Second); err != nil {
		return err
	}

	for _, dev := range d.devices {
		dev.Drop()
	}

	return nil
}

func NewDaemon(serverPort uint16, errHandler func(error)) Daemon {
	d := Daemon{
		devices:    make(map[string]irdevice.Device),
		errHandler: errHandler,
	}
	d.server = server.NewServer(uint32(serverPort), d.errHandler)

	d.server.AddHandler("GET", "/devices", d.getDevsReqWrapper(d.getDevicesHandler))
	d.server.AddHandler("GET", "/devices/:deviceName", d.getDevReqWrapper(d.getDeviceHandler, "deviceName"))
	d.server.AddHandler("GET", "/receive/:deviceName", d.recvIRReqWrapper(d.receiveIRHandler, "deviceName"))
	d.server.AddHandler("POST", "/send/:deviceName", d.sendIRReqWrapper(d.sendIRHandler, "deviceName"))
	return d
}

func (d Daemon) Start() {
	d.server.Start()
}
