package database

import (
	"context"
	"template/go-api-server/pkg/logger"
	"template/go-api-server/storage"
)

type IExampleDao interface {
	InitExampleDao(logger logger.LoggerClientInterface)
	Read(query string, dest *ExampleModel, args ...interface{}) error
}

type ExampleDao struct {
	logger logger.LoggerClientInterface
}

var ExampleDaoImpl IExampleDao = &ExampleDao{}

func (dao *ExampleDao) InitExampleDao(logger logger.LoggerClientInterface) {
	dao.logger = logger.GetLoggerWithProfile("database.example")
}

func (dao *ExampleDao) Read(query string, dest *ExampleModel, args ...interface{}) error {
	dao.logger.Info("Executing read query", "query", query, "args", args)
	ctx := context.Background()
	err := storage.PgStorage.ExecuteRead(ctx, query, dest, args...)
	if err != nil {
		dao.logger.Error("Failed to execute read query", "query", query, "args", args, "error", err)
		return err
	}
	dao.logger.Info("Successfully executed read query", "query", query, "args", args)
	return nil
}
