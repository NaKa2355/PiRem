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

func errHandler(err error) {
	fmt.Println(err)
}

func main() {
	fmt.Println("starting daemon...")
	daemon := daemon.NewDaemon(8080, errHandler)
	fmt.Println("daemon started")

	for i := 0; i < 2; i++ {
		mockdev := ErrMockDev{}
		dev := irdevice.Device{}
		dev.InitMock("/test", 10*time.Second, mockdev)
		dev.StartDispatcher()
		daemon.AddDevice(fmt.Sprintf("test %d", i), dev)
	}

	daemon.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)
	daemon.Stop()
	log.Printf("daemon stopped")

}
