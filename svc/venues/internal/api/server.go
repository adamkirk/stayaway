package api

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

type ApiServerConfig interface {
	ApiServerPort() int
}

type Server struct {
	cfg ApiServerConfig
	e *echo.Echo
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.ApiServerPort()))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func NewServer(cfg ApiServerConfig) *Server {
	return &Server{
		cfg: cfg,
		e: echo.New(),
	}
}
