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
	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/service"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"
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

	insertMockUserProfiles(services)
	insertMockOrderForBtc(services)
	insertMockOrderForEth(services)
	fillUpInitialPositionForUserA(services)

	services.Logger.Info("Done!")
}

func insertMockUserProfiles(services *service.ServiceClient) {
	tarantoolClient := services.Services.Tarantool

	users := []string{"usera", "userb", "userc", "userd", "usere", "userf", "userg", "userh", "useri", "userj"}
	for _, user := range users {
		if user == "usera" {
			tarantoolClient.CreateUserWalletBalance(user, 20000.00)
			continue
		}

		tarantoolClient.CreateUserWalletBalance(user, 1000000.0)
	}
}

func insertMockOrderForBtc(services *service.ServiceClient) {
	services.Logger.Info("insertMockOrderForBtc")
	tarantoolClient := services.Services.Tarantool

	longOrders := []*tarantool_pkg.OrderStruct{
		{
			UserId:       "usera",
			Price:        22345.12,
			Market:       "btc",
			Side:         "1",
			PositionSize: 0.5,
		},
		{
			UserId:       "userb",
			Price:        22450.50,
			Market:       "btc",
			Side:         "1",
			PositionSize: 0.3,
		},
		{
			UserId:       "userc",
			Price:        22575.75,
			Market:       "btc",
			Side:         "1",
			PositionSize: 0.8,
		},
		{
			UserId:       "usera",
			Price:        22599.12,
			Market:       "btc",
			Side:         "1",
			PositionSize: 3,
		},
		{
			UserId:       "userd",
			Price:        22600.00,
			Market:       "btc",
			Side:         "1",
			PositionSize: 1.2,
		},
		{
			UserId:       "usere",
			Price:        22725.25,
			Market:       "btc",
			Side:         "1",
			PositionSize: 1.5,
		},
	}

	marketPrice := 22800.00

	shortOrders := []*tarantool_pkg.OrderStruct{
		{
			UserId:       "userf",
			Price:        22850.50,
			Market:       "btc",
			Side:         "-1",
			PositionSize: 0.7,
		},
		{
			UserId:       "userg",
			Price:        22975.75,
			Market:       "btc",
			Side:         "-1",
			PositionSize: 0.6,
		},
		{
			UserId:       "userh",
			Price:        23000.00,
			Market:       "btc",
			Side:         "-1",
			PositionSize: 0.9,
		},
		{
			UserId:       "useri",
			Price:        23125.25,
			Market:       "btc",
			Side:         "-1",
			PositionSize: 1.1,
		},
		{
			UserId:       "userj",
			Price:        23250.50,
			Market:       "btc",
			Side:         "-1",
			PositionSize: 1.3,
		},
	}

	orders := append(longOrders, shortOrders...)

	for _, order := range orders {
		time.Sleep(500 * time.Millisecond)
		services.Logger.Info("insertin..", "order", order)
		tarantoolClient.InsertNewOrder(order)
	}

	tarantoolClient.UpdateMarketPrice("btc", marketPrice)
}

func insertMockOrderForEth(services *service.ServiceClient) {
	services.Logger.Info("insertMockOrderForEth")

	tarantoolClient := services.Services.Tarantool

	longOrders := []*tarantool_pkg.OrderStruct{
		{
			UserId:       "usera",
			Price:        1234.12,
			Market:       "eth",
			Side:         "1",
			PositionSize: 0.5,
		},
		{
			UserId:       "userb",
			Price:        1245.50,
			Market:       "eth",
			Side:         "1",
			PositionSize: 0.3,
		},
		{
			UserId:       "userc",
			Price:        1257.75,
			Market:       "eth",
			Side:         "1",
			PositionSize: 0.8,
		},
		{
			UserId:       "userd",
			Price:        1260.00,
			Market:       "eth",
			Side:         "1",
			PositionSize: 1.2,
		},
		{
			UserId:       "usere",
			Price:        1272.25,
			Market:       "eth",
			Side:         "1",
			PositionSize: 1.5,
		},
	}

	marketPrice := 1280.00

	shortOrders := []*tarantool_pkg.OrderStruct{
		{
			UserId:       "userf",
			Price:        1285.50,
			Market:       "eth",
			Side:         "-1",
			PositionSize: 0.7,
		},
		{
			UserId:       "userg",
			Price:        1297.75,
			Market:       "eth",
			Side:         "-1",
			PositionSize: 0.6,
		},
		{
			UserId:       "userh",
			Price:        1300.00,
			Market:       "eth",
			Side:         "-1",
			PositionSize: 0.9,
		},
		{
			UserId:       "useri",
			Price:        1312.25,
			Market:       "eth",
			Side:         "-1",
			PositionSize: 1.1,
		},
		{
			UserId:       "userj",
			Price:        1325.50,
			Market:       "eth",
			Side:         "-1",
			PositionSize: 1.3,
		},
	}

	orders := append(longOrders, shortOrders...)

	for _, order := range orders {
		time.Sleep(500 * time.Millisecond)
		services.Logger.Info("insertin..", "order", order)
		tarantoolClient.InsertNewOrder(order)
	}

	tarantoolClient.UpdateMarketPrice("eth", marketPrice)
}

func fillUpInitialPositionForUserA(services *service.ServiceClient) {
	services.Logger.Info("fillUpInitialPositionForUserA")
	tarantoolClient := services.Services.Tarantool

	positions := []*tarantool_pkg.PositionStruct{
		{
			UserID:       "usera",
			Market:       "btc",
			PositionSize: 2,
			AvgPrice:     30000.12,
			Side:         "1",
		},
		{
			UserID:       "usera",
			Market:       "eth",
			PositionSize: 3,
			AvgPrice:     2000.12,
			Side:         "1",
		},
	}

	for _, position := range positions {
		tarantoolClient.InsertPosition(position)
	}
}
