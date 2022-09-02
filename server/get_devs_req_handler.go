package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pirem/defs"
)

func (s Server) getDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("use GET method to get devices information: %s", defs.ErrInvaildMethod)
		w.Write(s.ErrorToJson(err))
		return
	}

	devs, err := s.handlers.GetDevicesHandler()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	resp, err := json.Marshal(devs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
