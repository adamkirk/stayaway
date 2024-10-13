// To override the name of a config field in the yaml file,  you need to use the
// mapstructure tag instead of the yaml tag as you may expect. Viper starts by
// unmarshalling the yaml to map[string]interface. access_log is one example
package config

type DbDriver string

const DbDriverMongoDb DbDriver = "mongodb"

var availableDbDrivers = []DbDriver{
	DbDriverMongoDb,
}

func (val DbDriver) IsKnown() bool {
	for _, chosen := range availableDbDrivers {
		if chosen == val {
			return true
		}
	}

	return false
}

func (val DbDriver) IsMongoDb() bool {
	return val == DbDriverMongoDb
}

type ConfigLogging struct {
	Level  string
	Format string
}

type ConfigApiServerAccessLog struct {
	Enabled bool
	Format  string
}

type ConfigApiServer struct {
	DebugErrorsEnabled bool `mapstructure:"debug_errors_enabled"`
	Port               int
	AccessLog          ConfigApiServerAccessLog `mapstructure:"access_log"`
}

type ConfigApi struct {
	Server ConfigApiServer
}

type ConfigDbMongoDb struct {
	Uri                string
	Database           string
	ConnectionRetries  int    `mapstructure:"connection_retries"`
	MigrationsDatabase string `mapstructure:"migrations_database"`
}

type ConfigDb struct {
	Driver  DbDriver
	MongoDb ConfigDbMongoDb
}

type ConfigRedis struct {
	Host              string
	Password          *string
	Db                int
	ConnectionRetries int `mapstructure:"connection_retries"`
}

type ConfigMunicipalitiesSync struct {
	MaxProcesses int `mapstructure:"max_processes"`
	BatchSize    int `mapstructure:"batch_size"`
	Countries    []string
}

type ConfigMunicipalities struct {
	Sync ConfigMunicipalitiesSync
}

type Config struct {
	Logging        ConfigLogging
	Api            ConfigApi
	Db             ConfigDb
	Redis          ConfigRedis
	Municipalities ConfigMunicipalities
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

func (c *Config) DbDriver() DbDriver {
	return c.Db.Driver
}

func (c *Config) MongoDbUri() string {
	return c.Db.MongoDb.Uri
}

func (c *Config) MongoDbDatabase() string {
	return c.Db.MongoDb.Database
}

func (c *Config) MongoDbMigrationsDatabase() string {
	return c.Db.MongoDb.MigrationsDatabase
}

func (c *Config) MongoDbConnectionRetries() int {
	return c.Db.MongoDb.ConnectionRetries
}

func (c *Config) MunicipalitiesSyncBatchSize() int {
	return c.Municipalities.Sync.BatchSize
}

func (c *Config) MunicipalitiesSyncMaxProcesses() int {
	return c.Municipalities.Sync.MaxProcesses
}

func (c *Config) MunicipalitiesSyncCountries() []string {
	return c.Municipalities.Sync.Countries
}

func (c *Config) RedisHost() string {
	return c.Redis.Host
}

func (c *Config) RedisPassword() *string {
	return c.Redis.Password
}

func (c *Config) RedisDb() int {
	return c.Redis.Db
}

func (c *Config) RedisConnectionRetries() int {
	return c.Redis.ConnectionRetries
}

func NewDefault() *Config {
	return &Config{
		Logging: ConfigLogging{
			Level:  "info",
			Format: "json",
		},
		Api: ConfigApi{
			Server: ConfigApiServer{
				DebugErrorsEnabled: false,
				Port:               8080,
				AccessLog: ConfigApiServerAccessLog{
					Enabled: true,
					Format:  "json",
				},
			},
		},
		Db: ConfigDb{
			Driver: DbDriverMongoDb,
			MongoDb: ConfigDbMongoDb{
				Uri:                "",
				Database:           "organisations",
				ConnectionRetries:  3,
				MigrationsDatabase: "migrations",
			},
		},
		Redis: ConfigRedis{
			Host:              "",
			Password:          nil,
			Db:                0,
			ConnectionRetries: 3,
		},
		Municipalities: ConfigMunicipalities{
			Sync: ConfigMunicipalitiesSync{
				MaxProcesses: 10,
				BatchSize:    100,
				Countries:    []string{},
			},
		},
	}
}
