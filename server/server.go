package server

//デーモンのサーバー

import (
	"fmt"
	"net/http"
	"pirem/irdata"
	"pirem/irdevice"
)

type ServerHandlers struct {
	SendIRHandler     func(string, irdata.Data) error
	RecvIRDataHandler func(string) (irdata.Data, error)
	GetDevicesHandler func() (map[string]irdevice.Device, error)
	GetDeviceHandler  func(string) (irdevice.Device, error)
	ErrHandler        func(error)
}

type DaemonServer struct {
	handlers ServerHandlers
	port     uint16
}

func (s *DaemonServer) New(handlers ServerHandlers, port uint16) {
	s.handlers = handlers
	s.port = port
}

func (s DaemonServer) Start() {
	http.HandleFunc("/send/", s.sendHandler)
	http.HandleFunc("/receive/", s.receiveHandler)
	http.HandleFunc("/devices", s.getDevices)
	http.HandleFunc("/devices/", s.getDevice)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
