package server

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type Handler struct {
	method  string
	path    Path // exmple: "/users/:userId"
	handler HandlerFunc
}
