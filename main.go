package main

import (
	"encoding/json"
	"fmt"
	"pirem/irdata"
)

func main() {
	jsonData := []byte(`{"type":"raw","ir_data":{"data": ["10us","20us","30us"]}}`)
	data := irdata.Data{}
	json.Unmarshal(jsonData, &data)
	fmt.Println(data)
	jsonData, _ = json.Marshal(data)
	fmt.Println(string(jsonData))
	/*
		var dataType ir.DataType = ir.Raw
		jsonData, _ := json.Marshal(dataType)
		fmt.Println(string(jsonData))
	*/
	/*
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
	*/
}
