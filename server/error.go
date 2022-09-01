package server

import (
	"errors"
	"pirem/json/response"
)

var (
	ErrInvaildURLPath = errors.New("invaild URL path")
	ErrInvaildMethod  = errors.New("invaild http method")
)

func (s DaemonServer) ErrorToJson(inputErr error) []byte {
	json_data, err := response.ErrorToJson(inputErr)
	if err != nil {
		s.handlers.ErrHandler(inputErr)
		s.handlers.ErrHandler(err)
		return []byte("")
	}
	return json_data
}
