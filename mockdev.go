package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/NaKa2355/irdevctrl"
)

type ErrMockDev struct{}

func (dev ErrMockDev) ReceiveIRData() (irdevctrl.RawData, error) {
	fmt.Printf("receiving...")
	time.Sleep(5 * time.Second)
	return []irdevctrl.Pulse{irdevctrl.Pulse{irdevctrl.Micro, 10}, irdevctrl.Pulse{irdevctrl.Milli, 20}, irdevctrl.Pulse{irdevctrl.Micro, 30}}, nil
}

func (dev ErrMockDev) SendIRData(rawData irdevctrl.RawData) error {
	fmt.Printf("sending... %v\n", rawData)
	return nil
}

func (dev ErrMockDev) GetBufferSize() uint16 {
	return 600
}

func (dev ErrMockDev) GetSupportingFeatures() irdevctrl.Features {
	return irdevctrl.Features{irdevctrl.Sending, irdevctrl.Receiving}
}

func (dev ErrMockDev) Drop() error {
	return errors.New("test")
}
