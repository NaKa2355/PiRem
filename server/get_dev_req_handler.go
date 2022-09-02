package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pirem/defs"
	"strings"
)

func (s Server) getDevice(w http.ResponseWriter, r *http.Request) {
	pathes := strings.Split(r.URL.Path, "/")
	w.Header().Set("Content-Type", "text/json")

	if len(pathes) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("url path must be like this (/devices/device_name): %s", defs.ErrInvaildURLPath)
		w.Write(s.ErrorToJson(err))
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("use GET method to get device information: %s", defs.ErrInvaildMethod)
		w.Write(s.ErrorToJson(err))
		return
	}

	dev_name, err := url.QueryUnescape(pathes[2])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	dev, err := s.handlers.GetDeviceHandler(dev_name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	resp, err := json.Marshal(dev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
