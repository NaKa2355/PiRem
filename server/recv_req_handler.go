package server

//クライアントから赤外線の受信命令が来た時のハンドラ

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (s Server) receiveHandler(w http.ResponseWriter, r *http.Request) {
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

	irData, err := s.handlers.RecvIRDataHandler(dev_name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}

	resp, err := json.Marshal(irData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(s.ErrorToJson(err))
		return
	}
	w.Write(resp)
}
