package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pirem/defs"
	"pirem/irdata"
	"pirem/irdevice"
	"pirem/server"
	"time"
)

type Handler struct {
	errHandler func(error)
	handler    func(interface{}) (interface{}, error)
}

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

func (d Daemon) recvIRReqWrapper(handler func(string) (irdata.Data, error), devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		irData, err := handler(pathParam[devParamKey])
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(irData)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) sendIRReqWrapper(handler func(irdata.Data, string) error, devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		req, err := io.ReadAll(r.Body)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}
		irData := irdata.Data{}
		json.Unmarshal(req, &irData)

		err = handler(irData, pathParam[devParamKey])
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return f
}

func (d Daemon) getDevsReqWrapper(handler func() (irdevice.Devices, error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		devices, err := handler()
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(devices)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) getDevReqWrapper(handler func(string) (irdevice.Device, error), devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		device, err := handler(pathParam[devParamKey])
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(device)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) sendIRHandler(irdata irdata.Data, devName string) error {
	dev, exist := d.devices[devName]
	if !exist {
		return fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}

	return dev.SendIR(irdata)
}

func (d Daemon) receiveIRHandler(devName string) (irdata.Data, error) {
	irdata := irdata.Data{}
	dev, exist := d.devices[devName]
	if !exist {
		return irdata, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}

	return dev.ReceiveIR()
}

func (d Daemon) getDevicesHandler() (irdevice.Devices, error) {
	return d.devices, nil
}

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.devices[devName]
	if !exist {
		return dev, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}
	return dev, nil
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

func NewDaemon(server_port uint16, errHandler func(error)) {
	d := Daemon{
		devices:    make(map[string]irdevice.Device),
		errHandler: errHandler,
	}
	d.server = server.NewServer(uint32(server_port), d.errHandler)

	d.server.AddHandler("GET", "/devices", d.getDevsReqWrapper(d.getDevicesHandler))
	d.server.AddHandler("GET", "/devices/:deviceName", d.getDevReqWrapper(d.getDeviceHandler, "deviceName"))
	d.server.AddHandler("GET", "/receive/:deviceName", d.recvIRReqWrapper(d.receiveIRHandler, "deviceName"))
	d.server.AddHandler("POST", "/send/:deviceName", d.sendIRReqWrapper(d.sendIRHandler, "deviceName"))
}

func (d Daemon) Start() {
	d.server.Start()
}
