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
		// Launches (Product/Project Launches)
		launches := v1.Group("/launches")
		{
			launches.GET("", handler.ListLaunches)
			launches.POST("", handler.CreateLaunch)
			launches.GET("/:id", handler.GetLaunch)
			launches.PUT("/:id", handler.UpdateLaunch)
			launches.DELETE("/:id", handler.DeleteLaunch)
			// Nested routes for a specific launch
			launches.GET("/:id/milestones", handler.ListLaunchMilestones)
			launches.GET("/:id/tasks", handler.ListLaunchTasks)
		}

		// Milestones
		milestones := v1.Group("/milestones")
		{
			milestones.POST("", handler.CreateMilestone)
			milestones.GET("/:id", handler.GetMilestone)
			milestones.PUT("/:id", handler.UpdateMilestone)
			milestones.DELETE("/:id", handler.DeleteMilestone)
		}

		// Tasks
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", handler.CreateTask)
			tasks.GET("/:id", handler.GetTask)
			tasks.PUT("/:id", handler.UpdateTask)
			tasks.DELETE("/:id", handler.DeleteTask)
		}

		// Rockets
		rockets := v1.Group("/rockets")
		{
			rockets.GET("", handler.ListRockets)
			rockets.POST("", handler.CreateRocket)
			rockets.GET("/:id", handler.GetRocket)
			rockets.PUT("/:id", handler.UpdateRocket)
			rockets.DELETE("/:id", handler.DeleteRocket)
		}

		// Companies
		companies := v1.Group("/companies")
		{
			companies.GET("", handler.ListCompanies)
			companies.POST("", handler.CreateCompany)
			companies.GET("/:id", handler.GetCompany)
			companies.PUT("/:id", handler.UpdateCompany)
			companies.DELETE("/:id", handler.DeleteCompany)
		}

		// Launch Bases
		launchBases := v1.Group("/launch-bases")
		{
			launchBases.GET("", handler.ListLaunchBases)
			launchBases.POST("", handler.CreateLaunchBase)
			launchBases.GET("/:id", handler.GetLaunchBase)
			launchBases.PUT("/:id", handler.UpdateLaunchBase)
			launchBases.DELETE("/:id", handler.DeleteLaunchBase)
		}

		// Rocket Launches
		rocketLaunches := v1.Group("/rocket-launches")
		{
			rocketLaunches.GET("", handler.ListRocketLaunches)
			rocketLaunches.POST("", handler.CreateRocketLaunch)
			rocketLaunches.GET("/:id", handler.GetRocketLaunch)
			rocketLaunches.PUT("/:id", handler.UpdateRocketLaunch)
			rocketLaunches.DELETE("/:id", handler.DeleteRocketLaunch)
			rocketLaunches.POST("/sync", handler.SyncRocketLaunches)
		}

		// News
		news := v1.Group("/news")
		{
			news.GET("", handler.ListNews)
			news.POST("", handler.CreateNews)
			news.GET("/:id", handler.GetNews)
			news.PUT("/:id", handler.UpdateNews)
			news.DELETE("/:id", handler.DeleteNews)
		}
	}

	return router
}
