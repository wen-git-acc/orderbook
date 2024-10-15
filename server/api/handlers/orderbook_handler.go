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
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
}

func (client *HandlersClient) GetUserWalletHandler(context *gin.Context) {
	//Will auto handle if userId is not in the path
	userId := context.Param("userId")

	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello() + " for user " + userId,
	}

	context.JSON(200, resp)
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
