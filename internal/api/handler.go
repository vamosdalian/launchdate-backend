package api

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/database"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
	"github.com/vamosdalian/launchdate-backend/internal/service"
)

// Handler holds all API handlers
type Handler struct {
	launchService       *service.LaunchService
	taskService         *service.TaskService
	rocketService       *service.RocketService
	companyService      *service.CompanyService
	launchBaseService   *service.LaunchBaseService
	rocketLaunchService *service.RocketLaunchService
	newsService         *service.NewsService
	cache               *service.CacheService
	db                  *database.DB
	logger              *logrus.Logger
}

// NewHandler creates a new handler
func NewHandler(db *database.DB, cache *service.CacheService, logger *logrus.Logger) *Handler {
	// Initialize repositories
	launchRepo := repository.NewLaunchRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	rocketRepo := repository.NewRocketRepository(db.DB)
	companyRepo := repository.NewCompanyRepository(db.DB)
	launchBaseRepo := repository.NewLaunchBaseRepository(db.DB)
	rocketLaunchRepo := repository.NewRocketLaunchRepository(db.DB)
	newsRepo := repository.NewNewsRepository(db.DB)

	// Initialize services
	launchService := service.NewLaunchService(launchRepo, cache)
	taskService := service.NewTaskService(taskRepo, cache)
	rocketService := service.NewRocketService(rocketRepo, cache)
	companyService := service.NewCompanyService(companyRepo, cache)
	launchBaseService := service.NewLaunchBaseService(launchBaseRepo, cache)
	rocketLaunchService := service.NewRocketLaunchService(rocketLaunchRepo, cache)
	newsService := service.NewNewsService(newsRepo, cache)

	return &Handler{
		launchService:       launchService,
		taskService:         taskService,
		rocketService:       rocketService,
		companyService:      companyService,
		launchBaseService:   launchBaseService,
		rocketLaunchService: rocketLaunchService,
		newsService:         newsService,
		cache:               cache,
		db:                  db,
		logger:              logger,
	}
}
