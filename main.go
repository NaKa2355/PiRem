package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"pirem/daemon"
	"syscall"
)

const configFilePath = "/etc/piremd/config.json"

func errHandler(err error) {
	log.Println(err)
}

func main() {
	log.Println("starting daemon...")

	jsonConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Printf("faild to open file(%s): %s", configFilePath, err)
		os.Exit(1)
	}

	config := daemon.Config{}

	err = json.Unmarshal(jsonConfig, &config)
	if err != nil {
		log.Printf("faild to parse config(%s): %s", configFilePath, err)
	}

	daemon := daemon.NewDaemon(config.ServerPort, errHandler)
	log.Println("daemon started")

	for devName, dev := range config.Devices {
		dev.StartDispatcher()
		daemon.AddDevice(devName, dev)
	}

	daemon.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)
	daemon.Stop()
	log.Printf("daemon stopped")

}
