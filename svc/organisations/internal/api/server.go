package api

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouteHandler func(e echo.Context) error

type Route struct {
	Handler RouteHandler
	Path    string
	Name    string
	Method  string
}

type Routes []Route

type Controller interface {
	RegisterRoutes(g huma.API)
}

type ApiServerConfig interface {
	ApiServerPort() int
	ApiServerAccessLogEnabled() bool
	ApiServerAccessLogFormat() string
	ApiServerDebugErrorsEnabled() bool
}

type Server struct {
	cfg ApiServerConfig
	e   *echo.Echo
}

func (s *Server) Start() error {
	slog.Info(fmt.Sprintf("listening on port %d", s.cfg.ApiServerPort()))
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.ApiServerPort()))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

// @title Stayaway - Organisations
// @version 1.0
// @description This is an API for managing organisations in the stayaway ecosystem.

// @contact.name API Support
// @contact.url https://github.com/adamkirk/stayaway
// @contact.email adamkirk@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host organisations.stayaway.test
// @BasePath /api
func NewServer(v1Api *V1Api, cfg ApiServerConfig) *Server {
	e := echo.New()
	
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())

	if cfg.ApiServerAccessLogEnabled() {
		e.Use(buildLoggingMiddleware(cfg.ApiServerAccessLogFormat()))
	}
	
	for _, c := range v1Api.Controllers() {
		apiBase := fmt.Sprintf("/api/%s", v1Api.Version())
		api := e.Group(apiBase)
		cfg := huma.DefaultConfig("Organisations", v1Api.Version())
		// Needed to get the docs displaying properly.
		cfg.OpenAPI.Servers = []*huma.Server{
			{
				URL: apiBase,
			},
		}

		hg := humaecho.NewWithGroup(e, api, cfg)
		// hg.Use(v1Api.Middleware(cfg)...)
		c.RegisterRoutes(hg)
	}

	// e.GET("/openapi/*", echoSwagger.WrapHandler)

	srv := &Server{
		cfg: cfg,
		e:   e,
	}

	return srv
}
