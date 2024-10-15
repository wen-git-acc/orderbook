package service

import (
	"context"

	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/utils"
	"github.com/wen-git-acc/orderbook/storage"
	"github.com/wen-git-acc/orderbook/storage/database"
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
