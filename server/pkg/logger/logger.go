package logger

import (
	"log/slog"
	"os"
)

type LoggerClientInterface interface {
	GetLoggerWithProfile(profileName string) LoggerClientInterface
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
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

func (l *LoggerClient) Debug(message string, args ...interface{}) {
	l.logger.Debug(message, append([]interface{}{slog.String("name", l.name)}, args...)...)
}

func (l *LoggerClient) Info(message string, args ...interface{}) {
	l.logger.Info(message, append([]interface{}{slog.String("name", l.name)}, args...)...)
}

func (l *LoggerClient) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, append([]interface{}{slog.String("name", l.name)}, args...)...)
}

func (l *LoggerClient) Error(message string, args ...interface{}) {
	l.logger.Error(message, append([]interface{}{slog.String("name", l.name)}, args...)...)
}
