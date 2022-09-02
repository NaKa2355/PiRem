package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pirem/daemon"
	"pirem/irdevice"
	"syscall"
	"time"
)

func main() {

	daemon := daemon.Daemon{}
	daemon.Init()

	for i := 0; i < 2; i++ {
		mockdev := ErrMockDev{}
		dev := irdevice.Device{}
		dev.InitMock("/test", 10*time.Second, mockdev)
		dev.StartDispatcher()
		daemon.AddDevice(fmt.Sprintf("test %d", i), dev)
	}

	daemon.ErrHandler = func(err error) {
		log.Println(err.Error())
	}

	daemon.Start(8080)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)
	daemon.Stop()
	log.Printf("server stopped")

}
