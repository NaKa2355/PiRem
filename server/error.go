package server

import (
	"errors"
	"pirem/respjson"
)

var (
	ErrInvaildURLPath = errors.New("invaild URL path")
	ErrInvaildMethod  = errors.New("invaild http method")
)

func (s DaemonServer) ErrorToJson(inputErr error) []byte {
	json_data, err := respjson.ErrorToJson(inputErr)
	if err != nil {
		s.handlers.ErrHandler(inputErr)
		s.handlers.ErrHandler(err)
		return []byte("")
	}
	return json_data
}
