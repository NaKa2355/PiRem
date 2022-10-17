package server

import "net/http"

type RespHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type ReqHandlerFunc func(*http.Request, map[string]string) ([]byte, error)

type Handler struct {
	method  string
	path    Path // exmple: "/users/:userId"
	handler RespHandlerFunc
}
