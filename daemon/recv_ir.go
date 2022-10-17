package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pirem/defs"
	"pirem/irdata"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func recvIRReqWrapper(handler func(string) (irdata.Data, error), paramKey string) server.ReqHandlerFunc {
	f := func(r *http.Request, pathParam map[string]string) ([]byte, error) {
		irData, err := handler(pathParam[paramKey])
		if err != nil {
			return []byte("{}"), err
		}
		return json.Marshal(irData)
	}
	return f
}

func (d Daemon) receiveIRHandler(devName string) (irdata.Data, error) {
	irdata := irdata.Data{}
	dev, exist := d.devices[devName]
	if !exist {
		return irdata, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}

	return dev.ReceiveIR()
}
