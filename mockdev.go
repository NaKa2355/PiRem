package main

import (
	"errors"
	"fmt"

	"github.com/NaKa2355/ir"
)

type ErrMockDev struct{}

func (dev ErrMockDev) ReceiveIRData() (ir.RawData, error) {
	return []ir.Pulse{ir.Pulse{ir.Micro, 10}, ir.Pulse{ir.Micro, 20}, ir.Pulse{ir.Micro, 30}}, nil
}

func (dev ErrMockDev) SendIRData(rawData ir.RawData) error {
	fmt.Println(rawData)
	return errors.New("test")
}

func (dev ErrMockDev) GetBufferSize() uint16 {
	return 600
}

func (dev ErrMockDev) Drop() error {
	return errors.New("test")
}
