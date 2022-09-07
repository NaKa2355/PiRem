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
func (d Daemon) getDevReqWrapper(handler func(string) (irdevice.Device, error), devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		device, err := handler(pathParam[devParamKey])
		if err != nil {
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(device)
		if err != nil {
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) getDeviceHandler(devName string) (irdevice.Device, error) {
	dev, exist := d.devices[devName]
	if !exist {
		return dev, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}
	return dev, nil
}
