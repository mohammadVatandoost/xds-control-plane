package rest

import (
	"log/slog"
	"net/http"
)

func (s *Server) GetResources(w http.ResponseWriter, req *http.Request) {
	d, err := s.app.GetResources()
	if err != nil {
		slog.Error("couldn't get the resources", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(d)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetNodes(w http.ResponseWriter, req *http.Request) {
	d, err := s.app.GetNodes()
	if err != nil {
		slog.Error("couldn't get the nodes", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(d)
	w.WriteHeader(http.StatusOK)
}
