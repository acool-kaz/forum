package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
}

func (s *Server) Start(port string, handlers http.Handler) error {
	s.srv = http.Server{
		Addr:         port,
		Handler:      handlers,
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 3,
	}
	fmt.Printf("Server starting http://localhost%s\n", port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
