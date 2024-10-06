package config

import (
	"os"
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

type AppConfig struct {
	Mode        Mode   `default:"development"`
	ServiceName string `default:"go-api-server"`
}

func Load() *AppConfig {
	app := &AppConfig{}
	app.Mode = Mode(loadFromEnvironment("MODE", string(Production)))
	app.ServiceName = loadFromEnvironment("SERVICE_NAME", "go-api-server")
	// Load your config here from env or file
	return app
}

func loadFromEnvironment(envString, fallBackValue string) string {
	if val := os.Getenv(envString); val != "" {
		return val
	}
	return fallBackValue
}
