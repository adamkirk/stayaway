package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	apicmd "github.com/adamkirk-stayaway/organisations/cmd/api"
	dbmigrate "github.com/adamkirk-stayaway/organisations/cmd/db_migrate"
	dbping "github.com/adamkirk-stayaway/organisations/cmd/db_ping"
	municipalitiessync "github.com/adamkirk-stayaway/organisations/cmd/municipalities_sync"
	"github.com/adamkirk-stayaway/organisations/internal/api"
	v1 "github.com/adamkirk-stayaway/organisations/internal/api/v1"
	"github.com/adamkirk-stayaway/organisations/internal/config"
	"github.com/adamkirk-stayaway/organisations/internal/db"
	"github.com/adamkirk-stayaway/organisations/internal/domain/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
	"github.com/adamkirk-stayaway/organisations/internal/mutex"
	"github.com/adamkirk-stayaway/organisations/internal/repository"
	"github.com/adamkirk-stayaway/organisations/internal/validation"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cfgFile string
var appCfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Stayaway Organisations API service",
	Long:  `Blah`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var apiServeCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		apicmd.Handler(sharedOpts(), cmd, args)
	},
}

var municipalitiesCmd = &cobra.Command{
	Use:   "municipalities",
	Short: "Municipalities commands",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var municipalitiesSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync municipalities",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		municipalitiessync.Handler(sharedOpts(), cmd, args)
	},
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "DB related commands",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dbInitCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the database",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		dbping.Handler(sharedOpts(), cmd, args)
	},
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		dbmigrate.Handler(sharedOpts(), cmd, args)
	},
}

func newFs() afero.Fs {
	return afero.NewOsFs()
}

func sharedOpts() []fx.Option {
	opts := []fx.Option{
		fx.Provide(
			fx.Annotate(
				buildConfig,
				fx.As(new(api.ApiServerConfig)),
				fx.As(new(v1.OrganisationsControllerConfig)),
				fx.As(new(v1.VenuesControllerConfig)),
				fx.As(new(v1.MunicipalitiesControllerConfig)),
				fx.As(new(db.MongoConfig)),
				fx.As(new(db.MongoDbMigratorConfig)),
				fx.As(new(municipalities.SyncHandlerConfig)),
				fx.As(new(db.RedisConnectorConfig)),
			),
		),
		fx.Provide(api.NewServer),
		fx.Provide(
			fx.Annotate(
				api.NewV1Api,
				fx.ParamTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewOrganisationsController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewVenuesController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewMunicipalitiesController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				v1.NewVenueAccommodationTemplatesController,
				fx.As(new(api.Controller)),
				fx.ResultTags(`group:"api.v1.controllers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewGetHandler,
				fx.As(new(v1.OrganisationsGetHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewListHandler,
				fx.As(new(v1.OrganisationsListHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewCreateHandler,
				fx.As(new(v1.OrganisationsCreateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewDeleteHandler,
				fx.As(new(v1.OrganisationsDeleteHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				organisations.NewUpdateHandler,
				fx.As(new(v1.OrganisationsUpdateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				validation.NewValidator,
				fx.As(new(organisations.Validator)),
				fx.As(new(venues.Validator)),
				fx.As(new(municipalities.Validator)),
				fx.As(new(accommodations.Validator)),
				fx.ParamTags(`group:"validationExtensions"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewValidationExtension,
				fx.As(new(validation.Extension)),
				fx.ResultTags(`group:"validationExtensions"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewValidationExtension,
				fx.As(new(validation.Extension)),
				fx.ResultTags(`group:"validationExtensions"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewCreateHandler,
				fx.As(new(v1.VenuesCreateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewListHandler,
				fx.As(new(v1.VenuesListHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewGetHandler,
				fx.As(new(v1.VenuesGetHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewDeleteHandler,
				fx.As(new(v1.VenuesDeleteHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				venues.NewUpdateHandler,
				fx.As(new(v1.VenuesUpdateHandler)),
			),
		),

		fx.Provide(
			fx.Annotate(
				municipalities.NewListHandler,
				fx.As(new(v1.MunicipalitiesListHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewCreateVenueTemplateHandler,
				fx.As(new(v1.VenueAccommodationTemplateCreateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewGetVenueTemplateHandler,
				fx.As(new(v1.VenueAccommodationTemplateGetHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewListVenueTemplatesHandler,
				fx.As(new(v1.VenueAccommodationTemplatesListHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewDeleteVenueTemplateHandler,
				fx.As(new(v1.VenueAccommodationTemplateDeleteHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				accommodations.NewUpdateVenueTemplateHandler,
				fx.As(new(v1.VenueAccommodationTemplateUpdateHandler)),
			),
		),
		fx.Provide(
			fx.Annotate(
				municipalities.NewSyncHandler,
			),
		),
		fx.Provide(
			fx.Annotate(
				newFs,
				fx.As(new(afero.Fs)),
			),
		),
		fx.Provide(
			fx.Annotate(
				db.NewRedisConnector,
				fx.As(new(mutex.RedisConnector)),
			),
		),
		fx.Provide(
			fx.Annotate(
				mutex.NewRedisMutex,
				fx.As(new(organisations.DistributedMutex)),
			),
		),
		fx.Provide(validation.NewValidationMapper),
	}

	if !appCfg.DbDriver().IsKnown() {
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
			fx.Provide(
				fx.Annotate(
					repository.NewMongoDbOrganisations,
					fx.As(new(organisations.GetHandlerRepo)),
					fx.As(new(organisations.ListHandlerRepo)),
					fx.As(new(organisations.CreateHandlerRepo)),
					fx.As(new(organisations.DeleteHandlerRepo)),
					fx.As(new(organisations.UpdateHandlerRepo)),
				),
			),
			fx.Provide(
				fx.Annotate(
					repository.NewMongoDbVenues,
					fx.As(new(venues.CreateHandlerRepo)),
					fx.As(new(venues.ListHandlerRepo)),
					fx.As(new(venues.GetHandlerRepo)),
					fx.As(new(venues.DeleteHandlerRepo)),
					fx.As(new(venues.UpdateHandlerRepo)),
					fx.As(new(accommodations.CreateVenueTemplateHandlerVenuesRepo)),
					fx.As(new(accommodations.GetVenueTemplateHandlerVenuesRepo)),
					fx.As(new(accommodations.DeleteVenueTemplateHandlerVenuesRepo)),
					fx.As(new(accommodations.ListVenueTemplatesHandlerVenuesRepo)),
					fx.As(new(accommodations.UpdateVenueTemplateHandlerVenuesRepo)),
				),
			),
			fx.Provide(
				fx.Annotate(
					repository.NewMongoDbMunicipalities,
					fx.As(new(municipalities.SyncHandlerRepo)),
					fx.As(new(municipalities.ListHandlerRepo)),
				),
			),
			fx.Provide(
				fx.Annotate(
					repository.NewMongoDbVenueAccommodationTemplates,
					fx.As(new(accommodations.CreateVenueTemplateHandlerRepo)),
					fx.As(new(accommodations.GetVenueTemplateHandlerRepo)),
					fx.As(new(accommodations.DeleteVenueTemplateHandlerRepo)),
					fx.As(new(accommodations.ListVenueTemplatesHandlerRepo)),
					fx.As(new(accommodations.UpdateVenueTemplateHandlerRepo)),
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
	rootCmd.PersistentFlags().String("log-level", "info", "log level to use")
	rootCmd.PersistentFlags().String("log-format", "json", "log format to use")
	rootCmd.PersistentFlags().Int("port", 8080, "port to serve API on")

	municipalitiesCmd.AddCommand(municipalitiesSyncCmd)

	dbCmd.AddCommand(dbInitCmd)
	dbCmd.AddCommand(dbMigrateCmd)

	rootCmd.AddCommand(apiServeCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(municipalitiesCmd)

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
		Level:     l,
	}

	var logger *slog.Logger

	if cfg.LogFormat() == "text" {
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
