package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (h *Handler) CreateRocket(c *gin.Context) {
	var req models.CreateRocketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocket, err := h.rocketService.CreateRocket(c.Request.Context(), &req)
	if err != nil {
		h.logger.Printf("failed to create rocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rocket"})
		return
	}

	c.JSON(http.StatusCreated, rocket)
}

func (h *Handler) GetRocket(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rocket, err := h.rocketService.GetRocket(c.Request.Context(), id)
	if err != nil {
		h.logger.Printf("failed to get rocket: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "rocket not found"})
		return
	}

	c.JSON(http.StatusOK, rocket)
}

func (h *Handler) ListRockets(c *gin.Context) {
	var active *bool
	if activeParam := c.Query("active"); activeParam != "" {
		activeBool, err := strconv.ParseBool(activeParam)
		if err == nil {
			active = &activeBool
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	rockets, err := h.rocketService.ListRockets(c.Request.Context(), active, limit, offset)
	if err != nil {
		h.logger.Printf("failed to list rockets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rockets"})
		return
	}

	c.JSON(http.StatusOK, rockets)
}

func (h *Handler) UpdateRocket(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CreateRocketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocket := &models.Rocket{
		Name:        req.Name,
		Description: req.Description,
		Height:      req.Height,
		Diameter:    req.Diameter,
		Mass:        req.Mass,
		CompanyID:   req.CompanyID,
		ImageURL:    req.ImageURL,
		Active:      req.Active,
	}

	if err := h.rocketService.UpdateRocket(c.Request.Context(), id, rocket); err != nil {
		h.logger.Printf("failed to update rocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rocket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rocket updated successfully"})
}

func (h *Handler) DeleteRocket(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.rocketService.DeleteRocket(c.Request.Context(), id); err != nil {
		h.logger.Printf("failed to delete rocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete rocket"})
		return
	}

	c.Status(http.StatusNoContent)
}
