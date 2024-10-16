package service

import (
	"context"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/matching_engine.go"
	tarantoolPkg "github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"

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
	Utils          utils.UtilsClientInterface
	Tarantool      tarantoolPkg.TarantoolClientInterface
	MatchingEngine matching_engine.MatchingEngineInterface
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

func (s *ServiceClient) RegisterTarantoolPackage(conn *tarantool.Connection) *ServiceClient {
	s.Logger.Info("Registering tarantool package")

	s.Services.Tarantool = tarantoolPkg.NewTarantoolClient(&tarantoolPkg.TarantoolClientOptions{
		Logger: s.Logger.GetLoggerWithProfile("tarantool"),
		Conn:   conn,
		Utils:  s.Services.Utils,
	})

	return s
}

func (s *ServiceClient) RegisterMatchingEnginePackage() *ServiceClient {
	s.Logger.Info("Registering matching engine package")

	s.Services.MatchingEngine = matching_engine.NewMatchingEngine(&matching_engine.MatchingEngineOptions{
		Logger:    s.Logger.GetLoggerWithProfile("matching_engine"),
		Utils:     s.Services.Utils,
		Tarantool: s.Services.Tarantool,
	})

	return s
}
