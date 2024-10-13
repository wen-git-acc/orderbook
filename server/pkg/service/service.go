package service

import (
	"context"
	"template/go-api-server/config"
	"template/go-api-server/pkg/logger"
	"template/go-api-server/pkg/utils"
	"template/go-api-server/storage"
	"template/go-api-server/storage/database"
)

type ServiceClientConfig struct {
	Mode   config.Mode
	Ctx    context.Context
	AppCfg *config.AppConfig
}

type ServiceClient struct {
	Ctx          context.Context
	Mode         config.Mode
	Services     *Services
	DatabaseDaos *DatabaseDaos
	Logger       logger.LoggerClientInterface
}

// All dependencies package register here
type Services struct {
	Utils utils.UtilsClientInterface
}

// All database daos register here
type DatabaseDaos struct {
	ExampleDao database.IExampleDao
}

func NewServiceClient(serviceClientConfig *ServiceClientConfig) *ServiceClient {

	serviceName := serviceClientConfig.AppCfg.ServiceName
	return &ServiceClient{
		Ctx:      serviceClientConfig.Ctx,
		Mode:     serviceClientConfig.Mode,
		Services: &Services{},
		Logger:   logger.NewLoggerClient(serviceName),
	}
}

func (s *ServiceClient) RegisterDbPackage(cfg *config.PgSqlConfig) *ServiceClient {
	s.Logger.Info("Registering db package")

	// Must register before dao init, as it creates database connection.
	storage.PgStorage.InitCfg(s.Mode, cfg, s.Logger)

	database.ExampleDaoImpl.InitExampleDao(s.Logger)

	s.DatabaseDaos = &DatabaseDaos{
		ExampleDao: database.ExampleDaoImpl,
	}

	return s
}

func (s *ServiceClient) RegisterUtilsPackage() *ServiceClient {
	s.Logger.Info("Registering utils package")

	s.Services.Utils = utils.NewUtilsClient(&utils.UtilsClientOptions{
		Logger: s.Logger.GetLoggerWithProfile("utils")})

	return s
}
