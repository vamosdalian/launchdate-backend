package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (h *Handler) CreateNews(c *gin.Context) {
	var req models.CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news, err := h.newsService.CreateNews(c.Request.Context(), &req)
	if err != nil {
		h.logger.Printf("failed to create news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create news"})
		return
	}

	c.JSON(http.StatusCreated, news)
}

func (h *Handler) GetNews(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	news, err := h.newsService.GetNews(c.Request.Context(), id)
	if err != nil {
		h.logger.Printf("failed to get news: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}

	c.JSON(http.StatusOK, news)
}

func (h *Handler) ListNews(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	newsList, err := h.newsService.ListNews(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Printf("failed to list news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list news"})
		return
	}

	c.JSON(http.StatusOK, newsList)
}

func (h *Handler) UpdateNews(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news := &models.News{
		Title:    req.Title,
		Summary:  req.Summary,
		Content:  req.Content,
		NewsDate: req.NewsDate,
		URL:      req.URL,
		ImageURL: req.ImageURL,
	}

	if err := h.newsService.UpdateNews(c.Request.Context(), id, news); err != nil {
		h.logger.Printf("failed to update news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "news updated successfully"})
}

func (h *Handler) DeleteNews(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.newsService.DeleteNews(c.Request.Context(), id); err != nil {
		h.logger.Printf("failed to delete news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete news"})
		return
	}

	c.Status(http.StatusNoContent)
}
