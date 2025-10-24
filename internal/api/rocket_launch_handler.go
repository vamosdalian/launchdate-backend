package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (h *Handler) CreateRocketLaunch(c *gin.Context) {
	var req models.CreateRocketLaunchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocketLaunch, err := h.rocketLaunchService.CreateRocketLaunch(c.Request.Context(), &req)
	if err != nil {
		h.logger.Printf("failed to create rocket launch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rocket launch"})
		return
	}

	c.JSON(http.StatusCreated, rocketLaunch)
}

func (h *Handler) GetRocketLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rocketLaunch, err := h.rocketLaunchService.GetRocketLaunch(c.Request.Context(), id)
	if err != nil {
		h.logger.Printf("failed to get rocket launch: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rocket launch not found"})
		return
	}

	c.JSON(http.StatusOK, rocketLaunch)
}

func (h *Handler) ListRocketLaunches(c *gin.Context) {
	var status *string
	if statusParam := c.Query("status"); statusParam != "" {
		status = &statusParam
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	rocketLaunches, err := h.rocketLaunchService.ListRocketLaunches(c.Request.Context(), status, limit, offset)
	if err != nil {
		h.logger.Printf("failed to list rocket launches: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rocket launches"})
		return
	}

	c.JSON(http.StatusOK, rocketLaunches)
}

func (h *Handler) UpdateRocketLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CreateRocketLaunchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocketLaunch := &models.RocketLaunch{
		Name:         req.Name,
		LaunchDate:   req.LaunchDate,
		RocketID:     req.RocketID,
		LaunchBaseID: req.LaunchBaseID,
		Status:       req.Status,
		Description:  req.Description,
	}

	if err := h.rocketLaunchService.UpdateRocketLaunch(c.Request.Context(), id, rocketLaunch); err != nil {
		h.logger.Printf("failed to update rocket launch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rocket launch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rocket launch updated successfully"})
}

func (h *Handler) DeleteRocketLaunch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.rocketLaunchService.DeleteRocketLaunch(c.Request.Context(), id); err != nil {
		h.logger.Printf("failed to delete rocket launch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete rocket launch"})
		return
	}

	c.Status(http.StatusNoContent)
}
