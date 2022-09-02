package server

//クライアントから赤外線の送信要求が来た時のサーバーのハンドラ

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pirem/defs"
	"pirem/irdata"
	"strings"
)

func (s Server) sendHandler(w http.ResponseWriter, r *http.Request) {
	pathes := strings.Split(r.URL.Path, "/")
	w.Header().Set("Content-Type", "text/json")
	if len(pathes) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("url path must be like this (/send/device_name): %s", defs.ErrInvaildURLPath)
		w.Write(s.ErrorToJson(err))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf("use POST method to send IR: %s", defs.ErrInvaildMethod)
		w.Write(s.ErrorToJson(err))
		return
	}

	buf := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	var irData irdata.Data
	err = json.Unmarshal(buf, &irData)
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

	err = s.handlers.SendIRHandler(dev_name, irData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}
	w.WriteHeader(http.StatusOK)
}
