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
	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "orders", "cancel"), handler.CancelOrderHandler) // delete order TODO

	router.POST(fmt.Sprintf("/%s/%s/%s", basePath, "user", "deposit"), handler.UserDepositHandler)         // user deposit done
	router.GET(fmt.Sprintf("/%s/%s/%s/:userId", basePath, "user", "wallet"), handler.GetUserWalletHandler) // user deposit done
	router.GET(fmt.Sprintf("/%s/%s", basePath, "user/:userId/positions"), handler.GetUserPositionHandler)  // get current user position done

	router.GET(fmt.Sprintf("/%s/:market", basePath), handler.GetOrderBookHandler)               // get orderbook done
	router.GET(fmt.Sprintf("/%s/%s", basePath, "market-price/:market"), handler.GetMarketPrice) // get market price done

	router.POST(fmt.Sprintf("/%s/%s", basePath, "insert/position"), handler.InsertPositionHandler) // Insert positions TODO: remove
	router.GET(fmt.Sprintf("/%s/%s", basePath, "view/positions"), handler.GetAllPositionsHandler)  // View all opening positions
}

// should i do delete position?
