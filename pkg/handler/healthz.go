package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthzHandler returns a Gin handler function
func HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
