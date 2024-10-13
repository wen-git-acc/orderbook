package handlers

import (
	"template/go-api-server/api/dto"
	"template/go-api-server/storage/database"

	"github.com/gin-gonic/gin"
)

type HelloHandlerInterface interface {
	GetHello(context *gin.Context)
}

func (client *HandlersClient) GetHello(context *gin.Context) {
	example := &database.ExampleModel{}
	err := client.packages.DatabaseDaos.ExampleDao.Read("SELECT * FROM slack_message_information WHERE id = 1", example)
	if err != nil {
		client.packages.Logger.Info("error", err)
	}

	resp := dto.HelloHandlerResponse{
		Message: client.packages.Services.Utils.GetHello(),
	}

	context.JSON(200, resp)
}
