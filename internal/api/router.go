package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/middleware"
)

// SetupRouter sets up the API routes
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(handler.logger))

	return router
}
