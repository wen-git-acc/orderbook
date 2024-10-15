package middlewares

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/wen-git-acc/orderbook/pkg/logger"
	"github.com/wen-git-acc/orderbook/pkg/service"

	"github.com/gin-gonic/gin"
)

type MiddlewareInterface interface {
	RequestLogger() gin.HandlerFunc
}

type MiddlewaresClient struct {
	packages *service.ServiceClient
	logger   logger.LoggerClientInterface
}

const (
	request     = "request"
	requestBody = "request body"
	loggingTime = "logging time"
	requestPath = "request path"
)

// This function should take in any dependencies that your handlers require and initialize all the handlers.
func NewMiddlewaresClient(services *service.ServiceClient) MiddlewareInterface {

	return &MiddlewaresClient{
		packages: services,
		logger:   services.Logger.GetLoggerWithProfile("middleware"),
	}
}

func (h *MiddlewaresClient) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		path := c.Request.URL.Path
		body, err := c.GetRawData()

		if err != nil {
			// h.logger.Error("Failed to read request body", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}

		// Format the current time as a string
		currentTime := t.Format(time.RFC3339)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore the body for further use, the body can only be read once.

		h.logger.Info(request, slog.String(requestPath, path), slog.String(requestBody, string(body)), slog.String(loggingTime, currentTime))

		c.Next()
	}
}
