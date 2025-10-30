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

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", handler.Health)
	}
	ll2 := apiV1.Group("/ll2")
	{
		ll2.GET("/launches", handler.GetLL2Launches)
		ll2.POST("/launches/update", handler.StartLL2LaunchUpdate)
		ll2.GET("/angecies", handler.GetLL2Angecy)
		ll2.POST("/angecies/update", handler.StartLL2AngecyUpdate)
	}

	return router
}
