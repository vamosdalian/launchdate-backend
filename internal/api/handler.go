package api

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/database"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
	"github.com/vamosdalian/launchdate-backend/internal/service"
)

// Handler holds all API handlers
type Handler struct {
	launchService    *service.LaunchService
	milestoneService *service.MilestoneService
	taskService      *service.TaskService
	cache            *service.CacheService
	db               *database.DB
	logger           *logrus.Logger
}

// NewHandler creates a new handler
func NewHandler(db *database.DB, cache *service.CacheService, logger *logrus.Logger) *Handler {
	// Initialize repositories
	launchRepo := repository.NewLaunchRepository(db)
	milestoneRepo := repository.NewMilestoneRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize services
	launchService := service.NewLaunchService(launchRepo, cache)
	milestoneService := service.NewMilestoneService(milestoneRepo, cache)
	taskService := service.NewTaskService(taskRepo, cache)

	return &Handler{
		launchService:    launchService,
		milestoneService: milestoneService,
		taskService:      taskService,
		cache:            cache,
		db:               db,
		logger:           logger,
	}
}
