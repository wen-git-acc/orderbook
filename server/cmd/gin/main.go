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

	// Initialize your service client, all dependencies will be registered here
	services := service.NewServiceClient(
		&service.ServiceClientConfig{
			Mode:   appCfg.Mode,
			Ctx:    context.Background(),
			AppCfg: appCfg,
		}).
		RegisterDbPackage(appCfg.PgSql.Master).
		RegisterUtilsPackage()

	// Initialize your handlers and middleware, while injecting the dependencies here.
	handlers := handlers.NewRouteHandlerImpl(services)
	middlewares := middlewares.NewMiddlewaresClient(services)

	routes := gin.Default()
	routes.Use(gin.Recovery())
	routes.Use(middlewares.RequestLogger())

	// Register your controllers here
	controllers.RegisterControllers(routes, handlers)

	routes.Run() // listen and serve on 0.0.0.0:8080
}
