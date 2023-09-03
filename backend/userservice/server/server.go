package server

import (
	config "userservice/configs"
)

type Server struct {
	cfg    *config.Config
	router *Router
}

func NewServer(cfg *config.Config, router *Router) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
	}
}

func (s *Server) Init() {
	s.router.InitRouter().Run(s.cfg.Port)
}
