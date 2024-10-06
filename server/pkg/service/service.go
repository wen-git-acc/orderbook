package service

import (
	"context"
	"template/go-api-server/config"
	"template/go-api-server/pkg/logger"
	"template/go-api-server/pkg/utils"
)

type ServiceClientConfig struct {
	Mode   config.Mode
	Ctx    context.Context
	AppCfg *config.AppConfig
}

type ServiceClient struct {
	Ctx      context.Context
	Mode     config.Mode
	Services *Services
	Logger   logger.LoggerClientInterface
}

type Services struct {
	Utils utils.UtilsClientInterface
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

func (s *ServiceClient) RegisterUtilsPackage() *ServiceClient {
	s.Logger.Info("Registering utils package")

	s.Services.Utils = utils.NewUtilsClient(&utils.UtilsClientOptions{
		Logger: s.Logger.GetLoggerWithProfile("utils")})

	return s
}
