package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/irdevice"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func getDevsReqWrapper(handler func() (irdevice.Devices, error)) server.ReqHandlerFunc {
	f := func(r *http.Request, pathParam map[string]string) ([]byte, error) {
		devices, err := handler()
		if err != nil {
			return []byte(""), err
		}
		return json.Marshal(devices)
	}
	return f
}

func (d Daemon) getDevicesHandler() (irdevice.Devices, error) {
	return d.devices, nil
}
