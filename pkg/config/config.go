package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig is configuration for Server
type AppConfig struct {
	// RestServerPort is the TCP port for REST API
	RestServerPort string
	// GRPCServerPort is the TCP port for GRPC API
	GRPCServerPort string
	// LogLevel is global log level: DEBUG, INFO, WARN, ERROR
	LogLevel string
	// DatabaseHost is the app main database host address
	DatabaseHost string
	// DatabaseUser is the app main database username
	DatabaseUser string
	// DatabasePassword is the app main database password
	DatabasePassword string
	// DatabaseName is the app main database name
	DatabaseName string
	// DatabasePort is the app main database address port
	DatabasePort string
	// DatabaseTimezone is the app main database timezone configuration
	DatabaseTimezone string
}

func replaceIfEmpty(v1 *string, v2 string) {
	if *v1 == "" {
		*v1 = v2
	}
}

func (c *AppConfig) Complete(config *AppConfig) *AppConfig {
	replaceIfEmpty(&c.RestServerPort, config.RestServerPort)
	replaceIfEmpty(&c.GRPCServerPort, config.GRPCServerPort)
	replaceIfEmpty(&c.LogLevel, config.LogLevel)
	replaceIfEmpty(&c.DatabaseHost, config.DatabaseHost)
	replaceIfEmpty(&c.DatabaseUser, config.DatabaseUser)
	replaceIfEmpty(&c.DatabasePassword, config.DatabasePassword)
	replaceIfEmpty(&c.DatabaseName, config.DatabaseName)
	replaceIfEmpty(&c.DatabasePort, config.DatabasePort)
	replaceIfEmpty(&c.DatabaseTimezone, config.DatabaseTimezone)
	return c
}

var DefaultConfig = &AppConfig{
	RestServerPort:   "3000",
	GRPCServerPort:   "3001",
	LogLevel:         "INFO",
	DatabaseHost:     "localhost",
	DatabaseUser:     "postgres",
	DatabasePassword: "",
	DatabaseName:     "question-and-answers",
	DatabasePort:     "5432",
	DatabaseTimezone: "America/Sao_Paulo",
}

func LoadEnvironmentConfiguration() *AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load() // The Original .env

	return &AppConfig{
		RestServerPort:   os.Getenv("REST_PORT"),
		GRPCServerPort:   os.Getenv("GRPC_PORT"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseTimezone: os.Getenv("DB_TIMEZONE"),
	}
}

func LoadParamConfig() *AppConfig {
	var cfg AppConfig
	flag.StringVar(&cfg.RestServerPort, "rest-port", "", "TCP port for REST API server to bind")
	flag.StringVar(&cfg.GRPCServerPort, "grpc-port", "", "TCP port for GRPC API server to bind")
	flag.StringVar(&cfg.LogLevel, "log-level", "", "Global log level")
	flag.StringVar(&cfg.DatabaseHost, "db-host", "", "Main database host")
	flag.StringVar(&cfg.DatabaseUser, "db-user", "", "Main database user")
	flag.StringVar(&cfg.DatabasePassword, "db-password", "", "Main database password")
	flag.StringVar(&cfg.DatabaseName, "db-name", "", "Main database name")
	flag.StringVar(&cfg.DatabasePort, "db-port", "", "Main database port")
	flag.StringVar(&cfg.DatabaseTimezone, "db-tz", "", "Main database timezone")
	flag.Parse()
	return &cfg
}

// LoadAppConfiguration loads the app configuration from respecting the following load order:
// - app run parameter,
// - dotenv (.env) config files (depends on the APP_ENV variable)
// - system environment variables
// - default values
func LoadAppConfiguration() *AppConfig {
	paramConfig := LoadParamConfig()
	envConfig := LoadEnvironmentConfiguration()
	return paramConfig.Complete(envConfig).Complete(DefaultConfig)
}
