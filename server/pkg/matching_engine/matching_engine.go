package matching_engine

import (
	"github.com/tarantool/go-tarantool/v2"
	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
	"github.com/wen-git-acc/orderbook/pkg/utils"
)

type MatchingEngineInterface interface {
}

type MatchingEngineOptions struct {
	Logger    logger.LoggerClientInterface
	Conn      *tarantool.Connection
	Utils     utils.UtilsClientInterface
	Tarantool tarantool_pkg.TarantoolClientInterface
}

type MatchingEngine struct {
	//Placeholder for logger client
	logger    logger.LoggerClientInterface
	utils     utils.UtilsClientInterface
	tarantool tarantool_pkg.TarantoolClientInterface
}

func NewMatchingEngine(opt *MatchingEngineOptions) MatchingEngineInterface {
	return &MatchingEngine{
		logger:    opt.Logger,
		tarantool: opt.Tarantool,
		utils:     opt.Utils,
	}
}
