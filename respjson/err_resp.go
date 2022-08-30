package respjson

import "encoding/json"

type ErrorResponce struct {
	Err string `json:"error"`
}

func ErrorToJson(err error) ([]byte, error) {
	errResp := ErrorResponce{Err: err.Error()}
	return json.Marshal(errResp)
}
