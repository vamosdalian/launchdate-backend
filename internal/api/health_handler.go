package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// Health checks the health of the application
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	response := models.HealthResponse{
		Status:   "ok",
		Database: "ok",
		Redis:    "ok",
	}

	// Check database
	if err := h.db.Health(); err != nil {
		response.Database = "error"
		response.Status = "degraded"
	}

	// Check Redis
	if err := h.cache.Health(c.Request.Context()); err != nil {
		response.Redis = "error"
		response.Status = "degraded"
	}

	statusCode := http.StatusOK
	if response.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
