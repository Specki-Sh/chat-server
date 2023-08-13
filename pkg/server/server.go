package server

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Addr           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func NewServer(c *Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           c.Addr,
			MaxHeaderBytes: c.MaxHeaderBytes,
			Handler:        handler,
			ReadTimeout:    c.ReadTimeout,
			WriteTimeout:   c.WriteTimeout,
		},
	}
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
