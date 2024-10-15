package controllers

import (
	"github.com/wen-git-acc/orderbook/api/handlers"

	"github.com/gin-gonic/gin"
)

func registerHealthCheckController(router *gin.Engine, handler handlers.HandlersInterface) {
	router.GET("/health_check", handler.GetHello)
}
