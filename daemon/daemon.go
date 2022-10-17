package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/irdevice"
	"pirem/server"
	"time"
)

type Daemon struct {
	devices    irdevice.Devices
	server     *server.Server
	errHandler func(error)
}

// デーモンに管理するデバイスを追加
func (d Daemon) AddDevice(name string, dev *irdevice.Device) error {
	d.devices[name] = dev
	return nil
}

// エラーをjsonにエンコードしてサーバーに送信
func sendError(inputErr error, w http.ResponseWriter, statusCode int, errHandler func(error)) {
	errJson := struct {
		Err string `json:"error"`
	}{}
	errJson.Err = inputErr.Error()

	resp, err := json.Marshal(errJson)
	if err != nil {
		errHandler(inputErr)
		errHandler(err)
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

func NewDaemon(serverPort uint16, errHandler func(error)) *Daemon {
	d := Daemon{
		devices:    irdevice.Devices{},
		errHandler: errHandler,
	}
	d.server = server.NewServer(uint32(serverPort), d.errHandler)

	d.server.AddHandler("GET", "/devices", respWrapper(getDevsReqWrapper(d.getDevicesHandler), d.errHandler))
	d.server.AddHandler("GET", "/devices/:deviceName", respWrapper(getDevReqWrapper(d.getDeviceHandler, "deviceName"), d.errHandler))
	d.server.AddHandler("GET", "/receive/:deviceName", respWrapper(recvIRReqWrapper(d.receiveIRHandler, "deviceName"), d.errHandler))
	d.server.AddHandler("POST", "/send/:deviceName", respWrapper(sendIRReqWrapper(d.sendIRHandler, "deviceName"), d.errHandler))
	return &d
}

func (d Daemon) Start() error {
	d.server.Start()
	return nil
}
