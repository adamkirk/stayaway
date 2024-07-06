package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
	RegisterRoutes(g *echo.Group)
}

type ApiServerConfig interface {
	ApiServerPort() int
	ApiServerAccessLogEnabled() bool
	ApiServerAccessLogFormat() string
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

// setupLoggingMiddleare adds a curtom logger to the given echo server
// We're purposely using slgo here so that it will default to whatever type of
// logger we initially used .e.g JSON or TEXT
func setupLoggingMiddleware(cfg ApiServerConfig, e *echo.Echo) {
	if !cfg.ApiServerAccessLogEnabled() {
		slog.Info("access logs disabled")
		return
	}
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// logger := slog.Default().With(slog.String("source", "access_log"))
	var logger *slog.Logger
	if cfg.ApiServerAccessLogFormat() == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	registeredRoutes := e.Routes()

	// See https://echo.labstack.com/docs/middleware/logger#examples
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogLatency:   true,
		LogMethod:    true,
		LogRequestID: true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

			level := slog.LevelInfo
			errorMsg := "nil"

			if v.Error != nil {
				errorMsg = v.Error.Error()
				level = slog.LevelError
			}

			var routeName string = "nil"

			for _, r := range registeredRoutes {
				if r.Method == c.Request().Method && r.Path == c.Path() {
					routeName = r.Name
					break
				}
			}

			logger.LogAttrs(context.Background(), level, "REQUEST",
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("request-id", v.RequestID),
				slog.String("log_type", "access"),
				slog.String("route", routeName),
				// Convert to milliseconds
				slog.Float64("duration", float64(v.Latency.Microseconds())/1000),
				slog.String("err", errorMsg),
			)
			return nil
		},
	}))
}

func NewServer(apiControllers []Controller, cfg ApiServerConfig) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())

	api := e.Group("/api")

	for _, c := range apiControllers {
		c.RegisterRoutes(api)
	}

	setupLoggingMiddleware(cfg, e)

	srv := &Server{
		cfg: cfg,
		e:   e,
	}

	return srv
}
