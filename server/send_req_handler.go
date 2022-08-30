package server

import (
	"io"
	"net/http"
	"net/url"
	"pirem/reqjson"
	"strings"
)

func (s DaemonServer) sendHandler(w http.ResponseWriter, r *http.Request) {
	pathes := strings.Split(r.URL.Path, "/")

	if len(pathes) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(ErrInvaildURLPath))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(ErrInvaildMethod))
		return
	}

	buf := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	rawData, err := reqjson.JsonToIRRawData(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.ErrorToJson(err))
		return
	}

	dev_name, err := url.QueryUnescape(pathes[2])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	err = s.handlers.SendIRHandler(dev_name, rawData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}
