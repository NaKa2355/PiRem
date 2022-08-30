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
	}

	dev, err := s.handlers.GetDevice(dev_name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
	}

	resp, err := respjson.DeviceToJson(dev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
