package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// CreateLaunch creates a new launch
// @Summary Create a new launch
// @Tags launches
// @Accept json
// @Produce json
// @Param launch body models.CreateLaunchRequest true "Launch data"
// @Success 201 {object} models.Launch
// @Router /api/v1/launches [post]
func (h *Handler) CreateLaunch(c *gin.Context) {
	var req models.CreateLaunchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, use a default owner ID (would come from auth in production)
	ownerID := int64(1)

	launch, err := h.launchService.CreateLaunch(c.Request.Context(), &req, ownerID)
	if err != nil {
		h.logger.WithError(err).Error("failed to create launch")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create launch"})
		return
	}

	c.JSON(http.StatusCreated, launch)
}

// GetLaunch retrieves a launch by ID
// @Summary Get a launch
// @Tags launches
// @Produce json
// @Param id path int true "Launch ID"
// @Success 200 {object} models.Launch
// @Router /api/v1/launches/{id} [get]
func (h *Handler) GetLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	launch, err := h.launchService.GetLaunch(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("failed to get launch")
		c.JSON(http.StatusNotFound, gin.H{"error": "launch not found"})
		return
	}

	c.JSON(http.StatusOK, launch)
}

// ListLaunches retrieves all launches
// @Summary List launches
// @Tags launches
// @Produce json
// @Param status query string false "Status filter"
// @Param priority query string false "Priority filter"
// @Param team_id query int false "Team ID filter"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Launch
// @Router /api/v1/launches [get]
func (h *Handler) ListLaunches(c *gin.Context) {
	status := c.Query("status")
	priority := c.Query("priority")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var teamID *int64
	if teamIDStr := c.Query("team_id"); teamIDStr != "" {
		id, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err == nil {
			teamID = &id
		}
	}

	launches, err := h.launchService.ListLaunches(c.Request.Context(), status, priority, teamID, limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("failed to list launches")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list launches"})
		return
	}

	c.JSON(http.StatusOK, launches)
}

// UpdateLaunch updates a launch
// @Summary Update a launch
// @Tags launches
// @Accept json
// @Produce json
// @Param id path int true "Launch ID"
// @Param launch body models.UpdateLaunchRequest true "Launch data"
// @Success 200 {object} gin.H
// @Router /api/v1/launches/{id} [put]
func (h *Handler) UpdateLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateLaunchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.launchService.UpdateLaunch(c.Request.Context(), id, &req); err != nil {
		h.logger.WithError(err).Error("failed to update launch")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update launch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "launch updated successfully"})
}

// DeleteLaunch deletes a launch
// @Summary Delete a launch
// @Tags launches
// @Param id path int true "Launch ID"
// @Success 204
// @Router /api/v1/launches/{id} [delete]
func (h *Handler) DeleteLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.launchService.DeleteLaunch(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("failed to delete launch")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete launch"})
		return
	}

	c.Status(http.StatusNoContent)
}
