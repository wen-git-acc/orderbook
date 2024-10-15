package tarantool_pkg

import (
	"github.com/tarantool/go-tarantool/v2"
	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/utils"
)

type TarantoolClientInterface interface {
	TarantoolUserConnInterface
	TarantoolOrderBookConnInterface
}

type TarantoolClientOptions struct {
	Logger logger.LoggerClientInterface
	Conn   *tarantool.Connection
	Utils  utils.UtilsClientInterface
}

type TarantoolClient struct {
	//Placeholder for logger client
	logger logger.LoggerClientInterface
	conn   *tarantool.Connection
	utils  utils.UtilsClientInterface
}

func NewTarantoolClient(opt *TarantoolClientOptions) TarantoolClientInterface {
	return &TarantoolClient{
		logger: opt.Logger,
		conn:   opt.Conn,
		utils:  opt.Utils,
	}
}
