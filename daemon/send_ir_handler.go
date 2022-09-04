package daemon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pirem/defs"
	"pirem/irdata"
	"pirem/server"
)

func (d Daemon) sendIRReqWrapper(handler func(irdata.Data, string) error, devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		req, err := io.ReadAll(r.Body)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
			return
		}

		irData := irdata.Data{}
		err = json.Unmarshal(req, &irData)
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusBadRequest)
			return
		}

		err = handler(irData, pathParam[devParamKey])
		if err != nil {
			d.errHandler(err)
			sendError(err, w, http.StatusInternalServerError)
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
