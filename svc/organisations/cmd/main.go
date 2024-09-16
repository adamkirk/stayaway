package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	apicmd "github.com/adamkirk-stayaway/organisations/cmd/api"
	"github.com/adamkirk-stayaway/organisations/cmd/dbmigrate"
	"github.com/adamkirk-stayaway/organisations/cmd/dbping"
	"github.com/adamkirk-stayaway/organisations/internal/api"
	"github.com/adamkirk-stayaway/organisations/internal/config"
	"github.com/adamkirk-stayaway/organisations/internal/db"
	"github.com/adamkirk-stayaway/organisations/internal/repository"
	"github.com/adamkirk-stayaway/organisations/pkg/organisations"
	"github.com/adamkirk-stayaway/organisations/pkg/validation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cfgFile string
var appCfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Stayaway Organisations API service",
	Long: `Blah`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var apiServeCmd = &cobra.Command{
	Use: "api",
	Short: "Start the API server",
	Long:`Blah`,
	Run: func (cmd *cobra.Command, args []string) {
		apicmd.Handler(sharedOpts(), cmd, args)
	},
}

var dbCmd = &cobra.Command{
	Use: "db",
	Short: "DB related commands",
	Long:`Blah`,
	Run: func (cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dbInitCmd = &cobra.Command{
	Use: "ping",
	Short: "Ping the database",
	Long:`Blah`,
	Run: func (cmd *cobra.Command, args []string) {
		dbping.Handler(sharedOpts(), cmd, args)
	},
}

var dbMigrateCmd = &cobra.Command{
	Use: "migrate",
	Short: "Migrate the database",
	Long:`Blah`,
	Run: func (cmd *cobra.Command, args []string) {
		dbmigrate.Handler(sharedOpts(), cmd, args)
	},
}

func sharedOpts() []fx.Option {
	opts := []fx.Option{
		fx.Provide(
			fx.Annotate(
				buildConfig,
				fx.As(new(api.ApiServerConfig)),
				fx.As(new(api.OrganisationsV1ControllerConfig)),
				fx.As(new(db.MongoConfig)),
				fx.As(new(db.MongoDbMigratorConfig)),
			),
		),
		fx.Provide(
			fx.Annotate(
				api.NewServer,
				fx.ParamTags(`group:"apiControllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				api.NewOrganisationsV1Controller,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"apiControllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				repository.NewMongoDbOrganisations,
				fx.As(new(api.OrganisationsRepo)),
				fx.As(new(organisations.GetHandlerRepo)),
				fx.As(new(organisations.ListHandlerRepo)),
				fx.As(new(organisations.CreateHandlerRepo)),
				fx.As(new(organisations.DeleteHandlerRepo)),
				fx.As(new(organisations.UpdateHandlerRepo)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewGetHandler,
				fx.As(new(api.OrganisationsGetHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewListHandler,
				fx.As(new(api.OrganisationsListHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewCreateHandler,
				fx.As(new(api.OrganisationsCreateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewDeleteHandler,
				fx.As(new(api.OrganisationsDeleteHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewUpdateHandler,
				fx.As(new(api.OrganisationsUpdateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				validation.NewValidator,
				fx.As(new(organisations.Validator)),
			),
		),
		fx.Provide(api.NewValidationMapper),
	}

	if ! appCfg.DbDriver().IsKnown() {
		slog.Error("Unknown db driver", "driver", string(appCfg.DbDriver()))
		os.Exit(1)
	}

	// Register difference implementations based on configured driver
	if appCfg.DbDriver().IsMongoDb() {
		opts = append(opts, []fx.Option{
			fx.Provide(db.NewMongoDbConnector),
			fx.Provide(
				fx.Annotate(
					db.NewMongoDbConnector,
					fx.As(new(repository.MongoDbConnector)),
				),
			),
			fx.Provide(
				fx.Annotate(
					db.NewMongoDbPinger,
					fx.As(new(db.Pinger)),
				),
			),
			fx.Provide(
				fx.Annotate(
					db.NewMongoDbMigrator,
					fx.As(new(db.Migrator)),
				),
			),
		}...)
	}
	return opts
}

func buildConfig() *config.Config {
	if appCfg != nil {
		return appCfg
	}

	c := config.NewDefault()
	err := viper.Unmarshal(c)
	cobra.CheckErr(err)

	appCfg = c
	return appCfg
}

func init() {
	cobra.OnInitialize(bootstrap)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is WORKING_DIRECTORY/config.yaml)")
	rootCmd.PersistentFlags().String("log-level",  "info", "log level to use")
	rootCmd.PersistentFlags().String("log-format",  "json", "log format to use")
	rootCmd.PersistentFlags().Int("port",  8080, "port to serve API on")

	dbCmd.AddCommand(dbInitCmd)
	dbCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(apiServeCmd)
	rootCmd.AddCommand(dbCmd)

	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("api.server.port", rootCmd.PersistentFlags().Lookup("port"))
}

func bootstrap() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		currentDir, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Tell viper to replace . in nested path with underscores
	// e.g. logging.level becomes LOGGING_LEVEL
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetEnvPrefix("stayaway")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	cobra.CheckErr(err)

	cfg := buildConfig()

	l := slog.Level(slog.LevelInfo)
	err = l.UnmarshalText([]byte(cfg.Logging.Level))

	cobra.CheckErr(err)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level: l,
	}

	var logger *slog.Logger

	if (cfg.LogFormat() == "text") {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	logger = logger.With(slog.String("log_type", "app"))
	slog.SetDefault(logger)
	// fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	  }
}