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
func getDevReqWrapper(handler func(string) (*irdevice.Device, error), paramKey string, errHandler func(error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		device, err := handler(pathParam[paramKey])
		if err != nil {
			sendError(err, w, http.StatusInternalServerError, errHandler)
			return
		}

		resp, err := json.Marshal(device)
		if err != nil {
			sendError(err, w, http.StatusInternalServerError, errHandler)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
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
