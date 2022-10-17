package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pirem/defs"
	"pirem/irdevice"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func getDevReqWrapper(handler func(string) (*irdevice.Device, error), paramKey string) server.ReqHandlerFunc {
	f := func(r *http.Request, pathParam map[string]string) ([]byte, error) {
		device, err := handler(pathParam[paramKey])
		if err != nil {
			return []byte(""), err
		}

		return json.Marshal(device)
	}
	return f
}

func (d Daemon) getDeviceHandler(devName string) (*irdevice.Device, error) {
	dev, exist := d.devices[devName]
	if !exist {
		return dev, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}
	return dev, nil
}
