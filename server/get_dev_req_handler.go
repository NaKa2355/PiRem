package server

import (
	"fmt"
	"net/http"
	"net/url"
	"pirem/respjson"
	"strings"
)

func (s DaemonServer) getDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	pathes := strings.Split(r.URL.Path, "/")
	w.Header().Set("Content-Type", "text/json")

	if len(pathes) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(ErrInvaildURLPath))
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(ErrInvaildMethod))
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

	resp, err := respjson.DeviceToJson(dev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
