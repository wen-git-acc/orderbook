package handlers

import (
	"github.com/wen-git-acc/orderbook/api/dto"

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
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
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

	isUserRegisterd := tarantoolClient.IsUserRegistered(depositRequest.UserID)

	if isUserRegisterd {
		currentBalance := tarantoolClient.GetUserWalletBalance(depositRequest.UserID)
		newBalance := currentBalance + depositRequest.DepositAmount
		err := tarantoolClient.UpdateUserWalletBalance(depositRequest.UserID, newBalance)
		if err != nil {
			context.JSON(500, gin.H{"error": "Internal system problem"})
			return
		}
	} else {
		err := tarantoolClient.CreateUserWalletBalance(depositRequest.UserID, depositRequest.DepositAmount)
		if err != nil {
			context.JSON(500, gin.H{"error": "Internal system problem"})
			return
		}
	}

	context.JSON(200, dto.UserDepositResponse{
		UserID:       depositRequest.UserID,
		WalletAmount: tarantoolClient.GetUserWalletBalance(depositRequest.UserID),
	})
}

func (client *HandlersClient) GetUserWalletHandler(context *gin.Context) {
	//Will auto handle if userId is not in the path
	userId := context.Param("userId")
	tarantoolClient := client.packages.Services.Tarantool

	walletAmount := tarantoolClient.GetUserWalletBalance(userId)

	context.JSON(200, dto.UserDepositResponse{
		UserID:       userId,
		WalletAmount: walletAmount,
	})
}

func (client *HandlersClient) GetOrderBookHandler(context *gin.Context) {
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
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
