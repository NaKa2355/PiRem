package newserver

import "net/http"

type Handler struct {
	method  string
	path    Path // exmple: "/users/:userId"
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}
