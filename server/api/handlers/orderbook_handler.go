package handlers

import (
	"fmt"
	"strings"

	"github.com/wen-git-acc/orderbook/api/dto"
	"github.com/wen-git-acc/orderbook/pkg/tarantool_pkg"

	"github.com/gin-gonic/gin"
)

type OrderBookHandlerInterface interface {
	InsertOrderHandler(context *gin.Context)
	DeleteOrderHandler(context *gin.Context)
	UserDepositHandler(context *gin.Context)
	GetUserWalletHandler(context *gin.Context)
	GetOrderBookHandler(context *gin.Context)
	GetMatchHistoryHandler(context *gin.Context)
	GetUserPositionHandler(context *gin.Context)
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

	total_price := insertOrderRequest.Price * insertOrderRequest.PositionSize
	a := tarantoolClient.GetUserWalletBalance(userId)
	fmt.Println("total_price", total_price, "a", a)
	fmt.Println("Smaller?", total_price < tarantoolClient.GetUserWalletBalance(userId))
	if (insertOrderRequest.Price * insertOrderRequest.PositionSize) > tarantoolClient.GetUserWalletBalance(userId) {
		context.JSON(400, gin.H{"error": "Insufficient balance"})
		return
	}

	err := tarantoolClient.InsertNewOrder(&tarantool_pkg.OrderStruct{
		UserId:       strings.ToLower(userId),
		Price:        insertOrderRequest.Price,
		Market:       strings.ToLower(insertOrderRequest.Market),
		Side:         strings.ToLower(insertOrderRequest.Side),
		PositionSize: insertOrderRequest.PositionSize,
	})

	if err != nil {
		context.JSON(500, gin.H{"error": "Internal system problem"})
		return
	}

	context.JSON(200, &dto.InsertOrderResponse{
		IsSuccess: true,
	})
}

func (client *HandlersClient) DeleteOrderHandler(context *gin.Context) {
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
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
			context.JSON(500, gin.H{"error": "Internal system problem"})
			return
		}
	} else {
		err := tarantoolClient.CreateUserWalletBalance(userId, depositRequest.DepositAmount)
		if err != nil {
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
	//Will auto handle if userId is not in the path
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

	tarantoolClient := client.packages.Services.Tarantool

	orders := tarantoolClient.GetAllOrders()

	context.JSON(200, &dto.GetAllOrderResponse{
		Orders: orders,
	})
}
func (client *HandlersClient) GetMatchHistoryHandler(context *gin.Context) {
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
}
func (client *HandlersClient) GetUserPositionHandler(context *gin.Context) {
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
}
