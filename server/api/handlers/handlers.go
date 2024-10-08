package handlers

import (
	"template/go-api-server/pkg/logger"
	"template/go-api-server/pkg/service"
)

type HandlersInterface interface {
	HelloHandlerInterface
}

type HandlersClient struct {
	packages *service.ServiceClient
	logger   logger.LoggerClientInterface
}

// This function should take in any dependencies that your handlers require and initialize all the handlers.
func NewRouteHandlerImpl(services *service.ServiceClient) HandlersInterface {
	return &HandlersClient{
		packages: services,
		logger:   services.Logger.GetLoggerWithProfile("handlers"),
	}
}
