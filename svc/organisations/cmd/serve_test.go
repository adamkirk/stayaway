package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/adamkirk-stayaway/organisations/internal/api"
	v1 "github.com/adamkirk-stayaway/organisations/internal/api/v1"
	"github.com/adamkirk-stayaway/organisations/internal/config"
	"github.com/adamkirk-stayaway/organisations/internal/db"
	"github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"
	"github.com/adamkirk-stayaway/organisations/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	fxtest "go.uber.org/fx/fxtest"
	"gopkg.in/yaml.v2"
)

var apiConfig string = `
logging:
  level: debug
  format: json

api:
  server:
    debug_errors_enabled: true
    port: 9999
    access_log:
      format: json
      enabled: false

db:
  driver: "mongodb"
  mongodb:
    ## Add a test container for this
    uri: mongodb://root:mongo-root-pass@mongo:27017/
    database: "organisations"

## Add a test container for this
redis:
  host: "redis:6379"
  password: "redis-pass"
  db: 0
  connection_retries: 3

municipalities:
  sync:
    max_processes: 100
    batch_size: 1000
    countries:
      - "United Kingdom"
`

func buildApiTestConfig() *config.Config {
	cfg := config.NewDefault()

	err := yaml.Unmarshal([]byte(apiConfig), cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}

func fxRun(t *testing.T, invokable any) {
	cfg := buildApiTestConfig()

	mockConfig := func() *config.Config {
		return cfg
	}

	// Okay so fx.Replace should be enough to replace thee config but it doesn't
	// seem to work. This gets everything but the config providers...need to 
	// figure out why fx.Replace ain't working, but wanna test this theory first...
	opts := sharedOpts(cfg)[2:]
	opts = append(
		opts,
		fx.Provide(
			mockConfig,
		),
		fx.Provide(
			fx.Annotate(
				mockConfig,
				fx.As(new(api.ApiServerConfig)),
				fx.As(new(v1.OrganisationsControllerConfig)),
				fx.As(new(v1.VenuesControllerConfig)),
				fx.As(new(municipalities.Config)),
				fx.As(new(db.RedisConnectorConfig)),
				fx.As(new(repository.MongoDBRepositoryConfig)),
			),
		),
		fx.NopLogger,
		fx.Invoke(invokable),
	)
	app := fxtest.New(t, opts...,)

	defer app.RequireStop()
	app.RequireStart()

}

// Basic test that the api starts listening, uses the startup probe.
func TestApiListens(t *testing.T) {
	fxRun(t, func (server *api.Server) {
		go server.Start()
		defer server.Shutdown(context.TODO())

		// Just wait for the server to be ready
		for i :=0; i< 5; i++{
			_, err := http.Get("http://localhost:9999/blah")

			if err == nil || i >= 5{
				break
			}

			time.Sleep(100 * time.Millisecond)
			i++
		}
		// ideally need to wait for the server to start, not sure the best way
		// although it's not causing an issue right now...

		res, err := http.Get("http://localhost:9999/api/v1/_probes/startup")
		require.Nil(t, err)
		assert.Equal(t, 204, res.StatusCode)
	})
}