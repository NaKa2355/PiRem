package server

import (
	"net/http"
	"pirem/respjson"
)

func (s DaemonServer) getDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(ErrInvaildMethod))
		return
	}

	devs, err := s.handlers.GetDevices()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	resp, err := respjson.DevicesToJson(devs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
