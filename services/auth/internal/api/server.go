package api

import (
	"gitlab.com/ajesicus/super_services/pkg/handler"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()

	// Public endpoints
	r.GET("/healthz", handler.HealthzHandler)

	return r
}
