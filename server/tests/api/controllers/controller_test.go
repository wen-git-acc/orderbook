package controllers_test

import (
	"context"

	"github.com/wen-git-acc/orderbook/api/handlers"
	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/service"
)

func getHandlers() handlers.HandlersInterface {
	appCfg := config.Load()

	services := service.NewServiceClient(
		&service.ServiceClientConfig{
			Mode:   appCfg.Mode,
			Ctx:    context.Background(),
			AppCfg: appCfg,
		}).
		RegisterUtilsPackage()

	return handlers.NewRouteHandlerImpl(services)
}
