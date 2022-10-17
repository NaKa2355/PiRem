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
func sendIRReqWrapper(handler func(irdata.Data, string) error, paramKey string) server.ReqHandlerFunc {
	f := func(r *http.Request, pathParam map[string]string) ([]byte, error) {
		body := []byte("")
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return body, err
		}

		irData := irdata.Data{}
		err = json.Unmarshal(req, &irData)
		if err != nil {
			return body, err
		}

		err = handler(irData, pathParam[paramKey])
		if err != nil {
			return body, err
		}
		return body, nil
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
