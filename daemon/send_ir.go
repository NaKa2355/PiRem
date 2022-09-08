package daemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil" //for go 1.11
	"net/http"
	"pirem/defs"
	"pirem/irdata"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func sendIRReqWrapper(handler func(irdata.Data, string) error, paramKey string, errHandler func(error)) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendError(err, w, http.StatusInternalServerError, errHandler)
			return
		}

		irData := irdata.Data{}
		err = json.Unmarshal(req, &irData)
		if err != nil {
			sendError(err, w, http.StatusBadRequest, errHandler)
			return
		}

		err = handler(irData, pathParam[paramKey])
		if err != nil {
			sendError(err, w, http.StatusInternalServerError, errHandler)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return f
}

func (d Daemon) sendIRHandler(irdata irdata.Data, devName string) error {
	dev, exist := d.devices[devName]
	if !exist {
		return fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}

	return dev.SendIR(irdata)
}
