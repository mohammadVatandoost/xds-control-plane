package rest

import (
	"net/http"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/rest"
)

type Server struct {
	conf *rest.RestAPIConfig
	app  App
	mux  *http.ServeMux
}

func NewServer(conf *rest.RestAPIConfig, app App) *Server {
	return &Server{
		conf: conf,
		app:  app,
		mux:  http.NewServeMux(),
	}
}
