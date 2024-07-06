package config

type ConfigLogging struct {
	Level string
	Format string
}

type ConfigApiServer struct {
	Port int
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