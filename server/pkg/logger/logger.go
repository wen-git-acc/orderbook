package logger

import (
	"log/slog"
	"os"
)

type LoggerClientInterface interface {
	GetLoggerWithProfile(profileName string) LoggerClientInterface
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
}

type LoggerClient struct {
	name   string
	logger *slog.Logger
}

func NewLoggerClient(serviceName string) LoggerClientInterface {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger = logger.With("service_name", serviceName)
	return &LoggerClient{
		name:   "",
		logger: logger,
	}
}

func (l *LoggerClient) GetLoggerWithProfile(name string) LoggerClientInterface {
	return &LoggerClient{
		name:   name,
		logger: l.logger,
	}
}

func (l *LoggerClient) Debug(message string) {
	l.logger.Debug(message, slog.String("name", l.name))
}

func (l *LoggerClient) Info(message string) {
	l.logger.Info(message, slog.String("name", l.name))
}

func (l *LoggerClient) Warn(message string) {
	l.logger.Warn(message, slog.String("name", l.name))
}

func (l *LoggerClient) Error(message string) {
	l.logger.Error(message, slog.String("name", l.name))
}
