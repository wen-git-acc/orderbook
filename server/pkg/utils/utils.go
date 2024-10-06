package utils

import (
	"template/go-api-server/pkg/logger"
)

type UtilsClientInterface interface {
	GetHello() string
}

type UtilsClientOptions struct {
	Logger logger.LoggerClientInterface
}

type UtilsClient struct {
	//Placeholder for logger client
	logger logger.LoggerClientInterface
}

func NewUtilsClient(opt *UtilsClientOptions) UtilsClientInterface {
	return &UtilsClient{
		logger: opt.Logger,
	}
}

// Example..
func (u *UtilsClient) GetHello() string {
	return "Hello"
}
