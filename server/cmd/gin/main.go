package main

import (
	"context"
	"template/go-api-server/api/controllers"
	"template/go-api-server/api/handlers"
	"template/go-api-server/api/middlewares"
	"template/go-api-server/config"
	"template/go-api-server/pkg/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load your config here from env or file
	appCfg := config.Load()

	services := service.NewServiceClient(
		&service.ServiceClientConfig{
			Mode:   appCfg.Mode,
			Ctx:    context.Background(),
			AppCfg: appCfg,
		}).
		RegisterUtilsPackage()

	handlers := handlers.NewRouteHandlerImpl(services)
	middlewares := middlewares.NewMiddlewaresClient(services)

	routes := gin.Default()
	routes.Use(gin.Recovery())
	routes.Use(middlewares.RequestLogger())
	controllers.RegisterControllers(routes, handlers)

	routes.Run() // listen and serve on 0.0.0.0:8080
}
