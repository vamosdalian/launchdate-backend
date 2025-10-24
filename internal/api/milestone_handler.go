package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// CreateMilestone creates a new milestone
// @Summary Create a new milestone
// @Tags milestones
// @Accept json
// @Produce json
// @Param milestone body models.CreateMilestoneRequest true "Milestone data"
// @Success 201 {object} models.Milestone
// @Router /api/v1/milestones [post]
func (h *Handler) CreateMilestone(c *gin.Context) {
	var req models.CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	milestone, err := h.milestoneService.CreateMilestone(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("failed to create milestone")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create milestone"})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// GetMilestone retrieves a milestone by ID
// @Summary Get a milestone
// @Tags milestones
// @Produce json
// @Param id path int true "Milestone ID"
// @Success 200 {object} models.Milestone
// @Router /api/v1/milestones/{id} [get]
func (h *Handler) GetMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	milestone, err := h.milestoneService.GetMilestone(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("failed to get milestone")
		c.JSON(http.StatusNotFound, gin.H{"error": "milestone not found"})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// ListLaunchMilestones retrieves all milestones for a launch
// @Summary List milestones for a launch
// @Tags milestones
// @Produce json
// @Param id path int true "Launch ID"
// @Success 200 {array} models.Milestone
// @Router /api/v1/launches/{id}/milestones [get]
func (h *Handler) ListLaunchMilestones(c *gin.Context) {
	launchID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid launch_id"})
		return
	}

	milestones, err := h.milestoneService.ListMilestonesByLaunchID(c.Request.Context(), launchID)
	if err != nil {
		h.logger.WithError(err).Error("failed to list milestones")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

// UpdateMilestone updates a milestone
// @Summary Update a milestone
// @Tags milestones
// @Accept json
// @Produce json
// @Param id path int true "Milestone ID"
// @Param milestone body models.UpdateMilestoneRequest true "Milestone data"
// @Success 200 {object} gin.H
// @Router /api/v1/milestones/{id} [put]
func (h *Handler) UpdateMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.milestoneService.UpdateMilestone(c.Request.Context(), id, &req); err != nil {
		h.logger.WithError(err).Error("failed to update milestone")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update milestone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "milestone updated successfully"})
}

// DeleteMilestone deletes a milestone
// @Summary Delete a milestone
// @Tags milestones
// @Param id path int true "Milestone ID"
// @Success 204
// @Router /api/v1/milestones/{id} [delete]
func (h *Handler) DeleteMilestone(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.milestoneService.DeleteMilestone(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("failed to delete milestone")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete milestone"})
		return
	}

	c.Status(http.StatusNoContent)
}
