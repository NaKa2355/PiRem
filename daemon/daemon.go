package daemon

import (
	"encoding/json"
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

func (d Daemon) sendError(inputErr error, w http.ResponseWriter, statusCode int) {
	errJson := struct {
		Err string `json:"error"`
	}{}
	errJson.Err = inputErr.Error()

	resp, err := json.Marshal(errJson)
	if err != nil {
		d.errHandler(inputErr)
		d.errHandler(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(resp)
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
