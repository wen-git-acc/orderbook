package handlers

import (
	"strings"

	"github.com/wen-git-acc/orderbook/api/dto"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"

	"github.com/gin-gonic/gin"
)

type OrderBookHandlerInterface interface {
	InsertOrderHandler(context *gin.Context)
	CancelOrderHandler(context *gin.Context)
	UserDepositHandler(context *gin.Context)
	GetUserWalletHandler(context *gin.Context)
	GetOrderBookHandler(context *gin.Context)
	GetUserPositionHandler(context *gin.Context)
	GetAllPositionsHandler(context *gin.Context)
	GetMarketPrice(context *gin.Context)
}

func (client *HandlersClient) InsertOrderHandler(context *gin.Context) {
	var insertOrderRequest dto.InsertOrderRequest

	if err := context.ShouldBindJSON(&insertOrderRequest); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := strings.ToLower(insertOrderRequest.UserId)
	tarantoolClient := client.packages.Services.Tarantool
	if !tarantoolClient.IsUserRegistered(userId) {
		context.JSON(400, gin.H{"error": "User not registered"})
		return
	}

	order := &tarantool_pkg.OrderStruct{
		UserId:       strings.ToLower(userId),
		Price:        insertOrderRequest.Price,
		Market:       strings.ToLower(insertOrderRequest.Market),
		Side:         strings.ToLower(insertOrderRequest.Side),
		PositionSize: insertOrderRequest.PositionSize,
	}
	// To check if there is a active opposite order
	_, err := tarantoolClient.GetNetPositionSizeByValidatingPosition(order)

	if err != nil {
		client.logger.Error("opposite order logic error", err)
		context.JSON(400, gin.H{"error": "Opposite order logic error"})
		return
	}

	if order.PositionSize == 0 {
		context.JSON(200, &dto.InsertOrderResponse{
			IsSuccess: true,
		})
		return
	}

	userWalletBalance := tarantoolClient.GetUserWalletBalance(userId)

	if (insertOrderRequest.Price * insertOrderRequest.PositionSize) >= userWalletBalance {
		context.JSON(400, gin.H{"error": "Insufficient balance"})
		return
	}

	tarantoolClient.OrderMatcher(order)

	context.JSON(200, &dto.InsertOrderResponse{
		IsSuccess: true,
	})
}

func (client *HandlersClient) CancelOrderHandler(context *gin.Context) {
	var deleteOrderRequest dto.DeleteOrderRequest

	if err := context.ShouldBindJSON(&deleteOrderRequest); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tarantoolClient := client.packages.Services.Tarantool

	order := tarantoolClient.GetOrderByPrimaryKey(deleteOrderRequest.UserId, deleteOrderRequest.Price, deleteOrderRequest.Side, deleteOrderRequest.Market)

	if order == nil {
		client.logger.Error("order not found, wrong order pass in")
		context.JSON(400, gin.H{"error": "Order not found"})
		return
	}

	err := tarantoolClient.DeleteOrderByPrimaryKey(deleteOrderRequest.UserId, deleteOrderRequest.Price, deleteOrderRequest.Side, deleteOrderRequest.Market)

	if err != nil {
		client.logger.Error("failed to delete order", err)
		context.JSON(500, gin.H{"error": "Internal system problem"})
		return
	}

	walletBalance := tarantoolClient.GetUserWalletBalance(deleteOrderRequest.UserId)
	refundedAmount := order.Price * order.PositionSize
	tarantoolClient.UpdateUserWalletBalance(deleteOrderRequest.UserId, walletBalance+refundedAmount)

	context.JSON(200, &dto.DeleteOrderResponse{
		IsSuccess: true,
	})
}

func (client *HandlersClient) UserDepositHandler(context *gin.Context) {

	var depositRequest dto.UserDepositRequest
	if err := context.ShouldBindJSON(&depositRequest); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tarantoolClient := client.packages.Services.Tarantool
	userId := strings.ToLower(depositRequest.UserID)
	isUserRegisterd := tarantoolClient.IsUserRegistered(userId)

	if isUserRegisterd {
		currentBalance := tarantoolClient.GetUserWalletBalance(userId)
		newBalance := currentBalance + depositRequest.DepositAmount
		err := tarantoolClient.UpdateUserWalletBalance(userId, newBalance)
		if err != nil {
			client.logger.Error("failed to update user wallet balance", err)
			context.JSON(500, gin.H{"error": "Internal system problem"})
			return
		}
	} else {
		err := tarantoolClient.CreateUserWalletBalance(userId, depositRequest.DepositAmount)
		if err != nil {
			client.logger.Error("failed to create user wallet balance", err)
			context.JSON(500, gin.H{"error": "Internal system problem"})
			return
		}
	}

	context.JSON(200, dto.UserDepositResponse{
		UserID:       userId,
		WalletAmount: tarantoolClient.GetUserWalletBalance(userId),
	})
}

func (client *HandlersClient) GetUserWalletHandler(context *gin.Context) {
	userId := context.Param("userId")
	tarantoolClient := client.packages.Services.Tarantool
	lowerUserId := strings.ToLower(userId)
	walletAmount := tarantoolClient.GetUserWalletBalance(lowerUserId)

	context.JSON(200, dto.UserDepositResponse{
		UserID:       lowerUserId,
		WalletAmount: walletAmount,
	})
}

func (client *HandlersClient) GetOrderBookHandler(context *gin.Context) {
	market := context.Param("market")

	tarantoolClient := client.packages.Services.Tarantool
	orders := tarantoolClient.GetOrderBook(strings.ToLower(market))

	context.JSON(200, orders)
}

func (client *HandlersClient) GetUserPositionHandler(context *gin.Context) {
	userId := context.Param("userId")

	tarantoolClient := client.packages.Services.Tarantool
	positions, err := tarantoolClient.GetUserPositions(strings.ToLower(userId))

	if err != nil {
		client.logger.Error("failed to get user position", err)
		context.JSON(500, gin.H{"error": "Internal system problem"})
		return
	}

	context.JSON(200, &dto.GetUserPositionResponse{
		Positions: positions,
	})
}

func (client *HandlersClient) GetAllPositionsHandler(context *gin.Context) {
	tarantoolClient := client.packages.Services.Tarantool
	positions, err := tarantoolClient.GetAllPositions()

	if err != nil {
		context.JSON(500, gin.H{"error": "Internal system problem"})
		return
	}

	context.JSON(200, &dto.GetAllPositionsResposne{
		Positions: positions,
	})
}

func (client *HandlersClient) GetMarketPrice(context *gin.Context) {
	market := context.Param("market")
	market = strings.ToLower(market)
	tarantoolClient := client.packages.Services.Tarantool
	marketPrice := tarantoolClient.GetMarketPriceByMarket(market)

	context.JSON(200, &dto.GetMarketPriceResponse{
		MarketPrice: marketPrice,
	})
}
