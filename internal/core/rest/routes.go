package rest

func (s *Server) routes() {
	s.mux.HandleFunc("/resources", s.GetResources)
	s.mux.HandleFunc("/nodes", s.GetNodes)
}
