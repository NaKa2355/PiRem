package server

import (
	"fmt"
	"net/http"
	"pirem/irdevice"

	"github.com/NaKa2355/ir"
)

type ServerHandlers struct {
	SendIRHandler func(string, ir.RawData) error
	ReceiveIRData func(string) (ir.RawData, error)
	GetDevices    func() (map[string]irdevice.Device, error)
	GetDevice     func(string) (irdevice.Device, error)
	ErrHandler    func(error)
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
