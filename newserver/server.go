package newserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvaildURLPath = errors.New("invaild URL path")
	ErrInvaildMathod  = errors.New("invaild http method")
)

type Server struct {
	server       http.Server
	handlers     []Handler
	errorHandler func(error)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	reqPath := NewPath(r.URL.Path)
	pathParam := map[string]string{}
	var err error
	handler := Handler{}

	for _, handler = range s.handlers {
		pathParam, err = reqPath.RoutePath(handler.path)
		if err == nil {
			break
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	if handler.method != r.Method {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	handler.handler(w, r, pathParam)
}

func NewServer(port uint32, errorHandler func(error)) *Server {
	s := Server{}
	s.handlers = make([]Handler, 0, 3)
	s.errorHandler = errorHandler

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)

	s.server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &s
}

func (s *Server) AddHandler(method string, path string, handler func(http.ResponseWriter, *http.Request, map[string]string)) {
	h := Handler{
		method:  method,
		path:    strings.Split(path, "/"),
		handler: handler,
	}
	s.handlers = append(s.handlers, h)
}

func (s *Server) Start() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			s.errorHandler(err)
		}
	}()
	return nil
}

func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
