package daemon

import (
	"encoding/json"
	"net/http"
	"pirem/server"
)

type Responce struct {
	Body    string `json:"body"`
	Message string `json:"message"`
}

func respWrapper(handler server.ReqHandlerFunc, errHandler func(error)) server.RespHandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		body, err := handler(r, pathParam)

		resp := Responce{}
		if err != nil {
			resp.Body = ""
			resp.Message = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			resp.Body = string(body)
			resp.Message = "success"
			w.WriteHeader(http.StatusOK)
		}

		strResp, err := json.Marshal(resp)

		if err != nil {
			errHandler(err)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(strResp)
	}
	return f
}