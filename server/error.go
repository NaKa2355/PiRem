package server

import (
	"encoding/json"
)

func (s Server) ErrorToJson(inputErr error) []byte {
	json_data, err := json.Marshal(inputErr.Error())
	if err != nil {
		s.handlers.ErrHandler(inputErr)
		s.handlers.ErrHandler(err)
		return []byte("")
	}

	s.handlers.ErrHandler(inputErr)

	return json_data
}
