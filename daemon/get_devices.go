package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/irdevice"
	"pirem/server"
)

func (d Daemon) getDevReqWrapper(handler func(string) (irdevice.Device, error), devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		device, err := handler(pathParam[devParamKey])
		if err != nil {
			d.errHandler(err)
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(device)
		if err != nil {
			d.errHandler(err)
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) getDevsReqWrapper(handler func() (irdevice.Devices, error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		devices, err := handler()
		if err != nil {
			d.errHandler(err)
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(devices)
		if err != nil {
			d.errHandler(err)
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
