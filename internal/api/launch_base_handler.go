package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (h *Handler) CreateLaunchBase(c *gin.Context) {
	var req models.CreateLaunchBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	launchBase, err := h.launchBaseService.CreateLaunchBase(c.Request.Context(), &req)
	if err != nil {
		h.logger.Printf("failed to create launch base: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create launch base"})
		return
	}

	c.JSON(http.StatusCreated, launchBase)
}

func (h *Handler) GetLaunchBase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	launchBase, err := h.launchBaseService.GetLaunchBase(c.Request.Context(), id)
	if err != nil {
		h.logger.Printf("failed to get launch base: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "launch base not found"})
		return
	}

	c.JSON(http.StatusOK, launchBase)
}

func (h *Handler) ListLaunchBases(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	launchBases, err := h.launchBaseService.ListLaunchBases(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Printf("failed to list launch bases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list launch bases"})
		return
	}

	c.JSON(http.StatusOK, launchBases)
}

func (h *Handler) UpdateLaunchBase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CreateLaunchBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	launchBase := &models.LaunchBase{
		Name:        req.Name,
		Location:    req.Location,
		Country:     req.Country,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := h.launchBaseService.UpdateLaunchBase(c.Request.Context(), id, launchBase); err != nil {
		h.logger.Printf("failed to update launch base: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update launch base"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "launch base updated successfully"})
}

func (h *Handler) DeleteLaunchBase(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.launchBaseService.DeleteLaunchBase(c.Request.Context(), id); err != nil {
		h.logger.Printf("failed to delete launch base: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete launch base"})
		return
	}

	c.Status(http.StatusNoContent)
}
