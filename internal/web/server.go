package web

import (
	"context"
	"net/http"
)

type Server struct {
	srv *http.Server
	mux *http.ServeMux
}

func New(opts ...Option) *Server {
	c := &Server{
		srv: &http.Server{
			Addr: ":8080",
		},
		mux: http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(*Server)

func (s *Server) Start() error {
	s.mux = http.NewServeMux()
	s.initRoutes()

	s.srv.Handler = s.mux
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(context.Background())
}

func (s *Server) initRoutes() {
	s.mux.HandleFunc("/", getNetworkStat)
}
