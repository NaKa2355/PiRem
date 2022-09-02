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
	/*
		jsonData := []byte(`{"type":"raw","ir_data":{"data": ["10us","20us","30us"]}}`)
		data := irdata.Data{}
		json.Unmarshal(jsonData, &data)
		fmt.Println(data)
		jsonData, _ = json.Marshal(data)
		fmt.Println(string(jsonData))

			var dataType ir.DataType = ir.Raw
			jsonData, _ := json.Marshal(dataType)
			fmt.Println(string(jsonData))
	*/
	daemon := daemon.Daemon{}
	daemon.Init()

	for i := 0; i < 2; i++ {
		mockdev := ErrMockDev{}
		dev := irdevice.Device{}
		dev.InitMock("/test", 10*time.Second, mockdev)

		//jsonData, _ := json.Marshal(dev)
		//fmt.Println(string(jsonData))
		dev.GenerateEventQueue()
		daemon.AddDevice(fmt.Sprintf("test %d", i), dev)
	}
	daemon.ErrHandler = func(err error) {
		println(err.Error())
	}

	daemon.Start(8080)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)
	daemon.Stop()
	log.Printf("server stopped")
}
