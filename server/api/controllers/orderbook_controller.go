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
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "orders", "insert"), handler.InsertOrderHandler)        // place order
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "orders", "delete"), handler.DeleteOrderHandler)        // delete order
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "user", "deposit"), handler.UserDepositHandler)         // user deposit
	router.GET(fmt.Sprintf("/%s/%s/%s/:userId", basePath, "user", "wallet"), handler.GetUserWalletHandler) // user deposit
	router.GET(fmt.Sprintf("/%s", basePath), handler.GetOrderBookHandler)                                  // get orderbook
	router.POST(fmt.Sprintf("/%s/%s", basePath, "match-history"), handler.GetMatchHistoryHandler)          // get match hisotry
	router.GET(fmt.Sprintf("/%s/%s", basePath, "market-price/:market"), handler.GetMarketPrice)            // get market price
	router.POST(fmt.Sprintf("/%s/%s", basePath, "position"), handler.GetUserPositionHandler)               // get current user position
	router.POST(fmt.Sprintf("/%s/%s", basePath, "insert/position"), handler.InsertPositionHandler)         // Insert positions
	router.POST(fmt.Sprintf("/%s/%s", basePath, "view/positions"), handler.ViewPositionsHandler)           // View all opening positions
}
