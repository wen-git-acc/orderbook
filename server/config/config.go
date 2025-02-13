package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Mode string

const (
	Production     Mode = "production"
	Staging        Mode = "staging"
	Development    Mode = "development"
	Test           Mode = "test"
	AcceptanceTest Mode = "acceptance-test"
	LoadTest       Mode = "load-test"
)

type PgSqlMasterSlave struct {
	Master *PgSqlConfig `json:"master"`
	Slave  *PgSqlConfig `json:"slave"`
}

type PgSqlConfig struct {
	PostgresDbUser      string `json:"postgresDbUser"`
	PostgresDbPassword  string `json:"postgresDbPassword"`
	PostgresDbHost      string `json:"postgresDbHost"`
	PostgresDbPort      int    `json:"postgresDbPort"`
	PostgresDbName      string `json:"postgresDbName"`
	PostgresMaxOverflow int    `json:"postgresMaxOverflow"`
	PostgresPoolSize    int    `json:"postgresPoolSize"`
	PostgresPoolTimeout int    `json:"postgresPoolTimeout"`
	PostgresPoolRecycle int    `json:"postgresPoolRecycle"`
}

type TarantoolConfig struct {
	Host     string `json:"tarantoolHost"`
	Port     string `json:"tarantoolPort"`
	User     string `json:"tarantoolUser"`
	Password string `json:"tarantoolPassword"`
}

type AppConfig struct {
	Mode        Mode   `default:"development"`
	ServiceName string `default:"go-api-server"`
	PgSql       *PgSqlMasterSlave
	Tarantool   *TarantoolConfig
}

func Load() *AppConfig {
	app := &AppConfig{}

	// Load Environment Mode
	app.Mode = Mode(loadFromEnvironment("MODE", string(Development)))

	if app.Mode == Development {
		loadDotEnvFile(".env.development")
	}

	// Load Service Name
	app.ServiceName = loadFromEnvironment("SERVICE_NAME", "go-api-server")

	// Load pgsql config
	loadPgConfig(app, app.Mode)
	loadTarantoolConfig(app, app.Mode)

	return app
}

func loadFromEnvironment(envString, fallBackValue string) string {
	if val := os.Getenv(envString); val != "" {
		return val
	}
	return fallBackValue
}

func loadDotEnvFile(fileName string) {
	// For local development, you can store your env variables in the .env file
	// Load the .env file by calling 'godotenv.Load()'
	// Example: godotenv.Load()

	pwd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get current working directory", "error", err)
		return
	}

	if err := godotenv.Load(filepath.Join(pwd, fmt.Sprintf("../../env/%s", fileName))); err != nil {
		slog.Error("Failed to load environment file", err)
	}
}

func loadPgConfig(appCfg *AppConfig, mode Mode) {
	pgSqlConfg := &PgSqlConfig{}

	// Unmarshal the JSON string into a struct
	if err := json.Unmarshal([]byte(loadFromEnvironment("MASTER_POSTGRES_CONFIG", "{}")), pgSqlConfg); err != nil {
		slog.Error("Error unmarshaling pgsql config: %v", err)
	}

	// Load the password from the local environment
	if mode == Development {
		pgSqlConfg.PostgresDbPassword = loadFromEnvironment("POSTGRES_DB_PASSWORD", "postgres_password")
	}

	appCfg.PgSql = &PgSqlMasterSlave{
		Master: pgSqlConfg,
	}
}

func loadTarantoolConfig(appCfg *AppConfig, mode Mode) {
	tarantoolConfig := &TarantoolConfig{}

	if err := json.Unmarshal([]byte(loadFromEnvironment("TARANTOOL_CONFIG", "{}")), tarantoolConfig); err != nil {
		slog.Error("Error unmarshaling tarantool config: %v", err)
	}

	if mode == Development {
		tarantoolConfig.Password = loadFromEnvironment("TARANTOOL_PASSWORD", "123456")
	}

	appCfg.Tarantool = tarantoolConfig
}
