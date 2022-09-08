package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/irdevice"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func getDevsReqWrapper(handler func() (irdevice.Devices, error), errHandler func(error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		devices, err := handler()
		if err != nil {
			sendError(err, w, http.StatusInternalServerError, errHandler)
			return
		}

		resp, err := json.Marshal(devices)
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

func (d Daemon) getDevicesHandler() (irdevice.Devices, error) {
	return d.devices, nil
}
