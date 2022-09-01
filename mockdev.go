package main

import (
	"errors"
	"fmt"

	"github.com/NaKa2355/irdevctrl"
)

type ErrMockDev struct{}

func (dev ErrMockDev) ReceiveIRData() (irdevctrl.RawData, error) {
	return []irdevctrl.Pulse{irdevctrl.Pulse{irdevctrl.Micro, 10}, irdevctrl.Pulse{irdevctrl.Micro, 20}, irdevctrl.Pulse{irdevctrl.Micro, 30}}, nil
}

func (dev ErrMockDev) SendIRData(rawData irdevctrl.RawData) error {
	fmt.Println(rawData)
	return errors.New("test")
}

func (dev ErrMockDev) GetBufferSize() uint16 {
	var test irdevctrl.Controller
	test.GetSupportingFeatures()
	return 600
}

func (dev ErrMockDev) GetSupportingFeatures() irdevctrl.Features {
	return irdevctrl.Features{irdevctrl.Sending, irdevctrl.Receiving}
}

func (dev ErrMockDev) Drop() error {
	return errors.New("test")
}
