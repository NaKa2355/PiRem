package server

//デーモンのサーバー

import (
	"context"
	"fmt"
	"net/http"
	"pirem/irdata"
	"pirem/irdevice"
	"time"
)

type ServerHandlers struct {
	SendIRHandler     func(string, irdata.Data) error
	RecvIRDataHandler func(string) (irdata.Data, error)
	GetDevicesHandler func() (map[string]irdevice.Device, error)
	GetDeviceHandler  func(string) (irdevice.Device, error)
	ErrHandler        func(error) //サーバーでエラーを表示できなかった場合に呼び出される
}

type Server struct {
	server   http.Server
	handlers ServerHandlers
}

func (s *Server) New(handlers ServerHandlers, port uint16) {
	s.handlers = handlers

	mux := http.NewServeMux()
	mux.HandleFunc("/send/", s.sendHandler)
	mux.HandleFunc("/receive/", s.receiveHandler)
	mux.HandleFunc("/devices", s.getDevices)
	mux.HandleFunc("/devices/", s.getDevice)

	s.server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

func (s Server) Start() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.handlers.ErrHandler(err)
		}
	}()
	return nil
}

func (s Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
