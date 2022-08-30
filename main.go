package main

import (
	"pirem/daemon"
	"pirem/irdevice"
	"pirem/tx"
)

func main() {
	mockdev := ErrMockDev{}
	devctrl := DevController{}
	devctrl.Init(mockdev)
	reqChan := make(chan tx.Request, 10)
	go devctrl.Start(reqChan)

	dev := irdevice.Device{}
	dev.Init("test.so", 600, reqChan)

	daemon := daemon.Daemon{}
	daemon.Init()
	daemon.AddDevice("airer", dev)
	daemon.Start(8080)
}
