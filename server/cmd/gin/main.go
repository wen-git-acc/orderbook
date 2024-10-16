package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/wen-git-acc/orderbook/api/controllers"
	"github.com/wen-git-acc/orderbook/api/handlers"
	"github.com/wen-git-acc/orderbook/api/middlewares"
	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load your config here from env or file
	appCfg := config.Load()

	// Initialize Tarantool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  fmt.Sprintf("%s:%s", appCfg.Tarantool.Host, appCfg.Tarantool.Port),
		User:     appCfg.Tarantool.User,
		Password: appCfg.Tarantool.Password,
	}

	// Temporily disable timeout
	opts := tarantool.Opts{
		// Timeout: time.Second,
	}

	tarantoolConn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		slog.Error("Connection refused:", err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Listen for interrupt (Ctrl + C) and termination signals

	go func() {
		<-sigs
		slog.Info("Got interrupt signal. Close connection")
		tarantoolConn.CloseGraceful()
		cancel()
	}()
	defer func() {
		tarantoolConn.CloseGraceful()
		slog.Info("Connection is closed")
	}()

	// Initialize your service client, all dependencies will be registered here
	services := service.NewServiceClient(
		&service.ServiceClientConfig{
			Mode:   appCfg.Mode,
			Ctx:    context.Background(),
			AppCfg: appCfg,
		}).
		RegisterDbPackage(appCfg.PgSql.Master).
		RegisterUtilsPackage().
		RegisterTarantoolPackage(tarantoolConn).
		RegisterMatchingEnginePackage()

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
