package controllers_test

import (
	"context"
	"template/go-api-server/api/handlers"
	"template/go-api-server/config"
	"template/go-api-server/pkg/service"
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
