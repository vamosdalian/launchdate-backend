package api

import (
	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/service"
)

// Handler holds all API handlers
type Handler struct {
	logger    *logrus.Logger
	ll2Server *service.LL2Service
}

// NewHandler creates a new handler
func NewHandler(logger *logrus.Logger, cfg *config.Config, db *db.MongoDB) *Handler {
	ll2server := service.NewLL2Service(db)
	return &Handler{
		logger:    logger,
		ll2Server: ll2server,
	}
}
