package controllers

import (
	"fmt"

	"github.com/wen-git-acc/orderbook/api/handlers"

	"github.com/gin-gonic/gin"
)

const (
	basePath = "orderbook"
)

func registerOrderBookController(router *gin.Engine, handler handlers.HandlersInterface) {
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "orders", "insert"), handler.InsertOrderHandler) // place order
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "orders", "cancel"), handler.CancelOrderHandler) // delete order

	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "user", "deposit"), handler.UserDepositHandler)         // user deposit
	router.GET(fmt.Sprintf("/%s/%s/%s/:userId", basePath, "user", "wallet"), handler.GetUserWalletHandler) // check balance
	router.GET(fmt.Sprintf("/%s/%s", basePath, "user/:userId/positions"), handler.GetUserPositionHandler)  // get current user position

	router.GET(fmt.Sprintf("/%s/:market", basePath), handler.GetOrderBookHandler)               // get orderbook
	router.GET(fmt.Sprintf("/%s/%s", basePath, "market-price/:market"), handler.GetMarketPrice) // get market price

	router.GET(fmt.Sprintf("/%s/%s", basePath, "view/positions"), handler.GetAllPositionsHandler) // view all opening positions
}
