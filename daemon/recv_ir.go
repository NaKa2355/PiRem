package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pirem/defs"
	"pirem/irdata"
	"pirem/server"
)

// net/httpのハンドラ関数をラップして扱いやすくする
func (d Daemon) recvIRReqWrapper(handler func(string) (irdata.Data, error), devParamKey string) server.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request, pathParam map[string]string) {
		irData, err := handler(pathParam[devParamKey])
		if err != nil {
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(irData)
		if err != nil {
			d.sendError(err, w, http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
	return f
}

func (d Daemon) receiveIRHandler(devName string) (irdata.Data, error) {
	irdata := irdata.Data{}
	dev, exist := d.devices[devName]
	if !exist {
		return irdata, fmt.Errorf("no such a device: %s", defs.ErrInvaildInput)
	}

	return dev.ReceiveIR()
}
