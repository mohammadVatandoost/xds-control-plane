package rest

import (
	"log/slog"
	"net/http"
)

type App interface {
	GetNodes() ([]byte, error)
	GetResources() ([]byte, error)
}

func (s *Server) Run() error {
	s.routes()
	slog.Info("running rest api server", "address", s.conf.String())
	return http.ListenAndServe(s.conf.String(), s.mux)
}
