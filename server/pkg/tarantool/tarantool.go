package tarantool

import (
	"github.com/wen-git-acc/orderbook/pkg/logger"
)

type TarantoolClientInterface interface {
	GetHello() string
}

type TarantoolClientOptions struct {
	Logger logger.LoggerClientInterface
}

type TarantoolClient struct {
	//Placeholder for logger client
	logger logger.LoggerClientInterface
}

func NewTarantoolClient(opt *TarantoolClientOptions) TarantoolClientInterface {
	return &TarantoolClient{
		logger: opt.Logger.GetLoggerWithProfile("tarantool_pkg"),
	}
}

// Example..
func (u *TarantoolClient) GetHello() string {
	return "Hello"
}
