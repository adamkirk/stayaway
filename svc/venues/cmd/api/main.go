package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/adamkirk-stayaway/venues/internal/api"
	"github.com/adamkirk-stayaway/venues/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Stayaway Venues service",
	Long: `Blah`,
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			fx.Provide(api.NewServer),
			fx.Provide(
				fx.Annotate(
					buildConfig,
					fx.As(new(api.ApiServerConfig)),
				),
			),
			fx.Invoke(startServer),
		  ).Run()
	},
}

func startServer(lc fx.Lifecycle, srv *api.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
		  go srv.Start()
		  return nil
		},
		OnStop: func(ctx context.Context) error {
		  return srv.Shutdown(ctx)
		},
	  })
}

func buildConfig() *config.Config {
	c := &config.Config{}
	err := viper.Unmarshal(c)
	cobra.CheckErr(err)

	return c
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is WORKING_DIRECTORY/config.yaml)")
	rootCmd.PersistentFlags().String("log-level",  "info", "log level to use")
	rootCmd.PersistentFlags().String("log-format",  "json", "log format to use")
	rootCmd.PersistentFlags().Int("port",  8080, "port to serve API on")

	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("logging.format", rootCmd.PersistentFlags().Lookup("log-format"))
	viper.BindPFlag("api.server.port", rootCmd.PersistentFlags().Lookup("port"))
}

func initConfig() {
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
	// fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	  }
}