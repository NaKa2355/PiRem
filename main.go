package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pirem/server"
	"syscall"
	"time"
)

func main() {

	s := server.NewServer(8080, func(err error) {
		fmt.Println(err)
	})
	s.AddHandler("GET", "/hello/:userId", func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		fmt.Printf("path param: %+v\n", pathParam)
		//fmt.Printf("----body----\n%s----", string(body))
	})
	s.AddHandler("GET", "/hello", func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		fmt.Printf("hello\n")
		//fmt.Printf("----body----\n%s----", string(body))
	})
	s.Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)
	if err := s.Stop(10 * time.Second); err != nil {
		fmt.Println(err)
	}

	/*
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
	*/
}
