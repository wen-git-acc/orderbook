package controllers

import (
	"template/go-api-server/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.Engine, handler handlers.HandlersInterface) {
	registerHealthCheckController(router, handler)
	registerOrderBookController(router, handler)
}
