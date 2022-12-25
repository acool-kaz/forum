package server

import (
	"context"
	"forum/internal/config"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
}

func (s *Server) Start(cfg *config.Config, handlers http.Handler) error {
	s.srv = http.Server{
		Addr:         ":" + cfg.Handler.Addr,
		Handler:      handlers,
		WriteTimeout: time.Second * time.Duration(cfg.Handler.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.Handler.ReadTimeout),
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
