package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	handlers := getHandlers()
	router := gin.Default()
	router.GET("/health_check", handlers.GetHello)

	req, _ := http.NewRequest("GET", "/health_check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Hello"}`, w.Body.String())
}
