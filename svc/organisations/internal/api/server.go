package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
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

func translateErrToHttpErr(err error) HttpError {
	switch t := err.(type) {
		default:
			return nil
		case model.ErrNotFound:
			return ErrNotFound{
				ResourceName: t.ResourceName,
			}
	}
}

func handleValidationError(ctx echo.Context, errs validation.ValidationError) {
	respErrors := map[string][]string{}

	for _, err := range errs.Errs {
		respErrors[err.Key] = err.Errors
	}

	respBody := map[string]any{
		"errors": respErrors,
	}

	ctx.JSON(422, respBody)
}

func setupErrorHandlingMiddleware(cfg ApiServerConfig, e *echo.Echo) {
	e.Use(func (next echo.HandlerFunc) echo.HandlerFunc {
		return func (ctx echo.Context) error {
			err := next(ctx); 
			if err == nil {
				return nil
			}

			if err, ok := err.(validation.ValidationError); ok {
				handleValidationError(ctx, err)
				return nil
			}

			respBody := map[string]any{}

			var httpErr HttpError

			if translated := translateErrToHttpErr(err); translated != nil {
				httpErr = translated
			} else {
				translated, ok := err.(HttpError);

				if ok {
					httpErr = translated
				}
			}

			if httpErr != nil {
				respBody["message"] = err.Error()

				debuggableErr, ok := err.(HttpDebuggableError)

				if ok && cfg.ApiServerDebugErrorsEnabled() {
					respBody["debug"] = map[string]any{
						"error": debuggableErr.DebugError(),
					}
				}

				ctx.JSON(httpErr.HttpStatusCode(), respBody)

				return nil
			}

			if cfg.ApiServerDebugErrorsEnabled() {
				respBody["debug"] = map[string]any{
					"error": err.Error(),
				}
			}

			respBody["message"] = "internal server error"
			ctx.JSON(500, respBody)

			return nil
		}
	})
}

func bindRequest(req any, ctx echo.Context) error {
	if reflect.ValueOf(req).Kind() != reflect.Ptr {
		slog.Error("cannot bind to non pointer", "path", ctx.Path())

		return errors.New("cannot bind request to non pointer value")
	}

	if err := ctx.Bind(req); err != nil {
		return ErrBadRequest{
			Message: "failed to parse request",
			DebugMessage: err.Error(),
		}
	}

	return nil
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
	setupErrorHandlingMiddleware(cfg, e)

	srv := &Server{
		cfg: cfg,
		e:   e,
	}

	return srv
}
