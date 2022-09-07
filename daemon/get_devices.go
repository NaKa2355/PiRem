package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/irdevice"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func (d Daemon) getDevsReqWrapper(handler func() (irdevice.Devices, error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		devices, err := handler()
		if err != nil {
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(devices)
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

func (d Daemon) getDevicesHandler() (irdevice.Devices, error) {
	return d.devices, nil
}
