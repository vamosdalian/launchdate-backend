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

	// Health check
	router.GET("/health", handler.Health)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Launches
		launches := v1.Group("/launches")
		{
			launches.GET("", handler.ListLaunches)
			launches.POST("", handler.CreateLaunch)
			launches.GET("/:id", handler.GetLaunch)
			launches.PUT("/:id", handler.UpdateLaunch)
			launches.DELETE("/:id", handler.DeleteLaunch)
		}

		// Milestones
		milestones := v1.Group("/milestones")
		{
			milestones.POST("", handler.CreateMilestone)
			milestones.GET("/:id", handler.GetMilestone)
			milestones.PUT("/:id", handler.UpdateMilestone)
			milestones.DELETE("/:id", handler.DeleteMilestone)
		}

		// Launch milestones
		v1.GET("/launches/:launch_id/milestones", handler.ListLaunchMilestones)

		// Tasks
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", handler.CreateTask)
			tasks.GET("/:id", handler.GetTask)
			tasks.PUT("/:id", handler.UpdateTask)
			tasks.DELETE("/:id", handler.DeleteTask)
		}

		// Launch tasks
		v1.GET("/launches/:launch_id/tasks", handler.ListLaunchTasks)
	}

	return router
}
