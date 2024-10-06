package handlers

import (
	"template/go-api-server/api/dto"

	"github.com/gin-gonic/gin"
)

type HelloHandlerInterface interface {
	GetHello(context *gin.Context)
}

func (client *HandlersClient) GetHello(context *gin.Context) {
	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}
	context.JSON(200, resp)
}
