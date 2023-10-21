package rest

import "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/rest"

type App interface {
	GetServices()
}

type Server struct {
	conf	*rest.RestAPIConfig		
}

func (s *Server) Run() error {
	
}

func NewServer(conf	*rest.RestAPIConfig) *Server {
	return &Server{
		conf: conf,
	}
}