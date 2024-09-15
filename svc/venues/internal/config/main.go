// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

type ConfigLogging struct {
	Level string
	Format string
}

type ConfigApiServerAccessLog struct {
	Enabled bool
	Format string
}

type ConfigApiServer struct {
	DebugErrorsEnabled bool `mapstructure:"debug_errors_enabled"`
	Port int
	AccessLog ConfigApiServerAccessLog `mapstructure:"access_log"`
}

type ConfigApi struct {
	Server ConfigApiServer
}

type Config struct {
	Logging ConfigLogging
	Api ConfigApi
}

func (c *Config) LogLevel() string {
	return c.Logging.Level
}

func (c *Config) LogFormat() string {
	return c.Logging.Format
}

func (c *Config) ApiServerPort() int {
	return c.Api.Server.Port
}

func (c *Config) ApiServerAccessLogEnabled() bool {
	return c.Api.Server.AccessLog.Enabled
}

func (c *Config) ApiServerAccessLogFormat() string {
	return c.Api.Server.AccessLog.Format
}

func (c *Config) ApiServerDebugErrorsEnabled() bool {
	return c.Api.Server.DebugErrorsEnabled
}