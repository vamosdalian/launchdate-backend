package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func (h *Handler) CreateCompany(c *gin.Context) {
	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company, err := h.companyService.CreateCompany(c.Request.Context(), &req)
	if err != nil {
		h.logger.Printf("failed to create company: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *Handler) GetCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	company, err := h.companyService.GetCompany(c.Request.Context(), id)
	if err != nil {
		h.logger.Printf("failed to get company: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (h *Handler) ListCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	companies, err := h.companyService.ListCompanies(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Printf("failed to list companies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list companies"})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (h *Handler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := &models.Company{
		Name:         req.Name,
		Description:  req.Description,
		Founded:      req.Founded,
		Founder:      req.Founder,
		Headquarters: req.Headquarters,
		Employees:    req.Employees,
		Website:      req.Website,
		ImageURL:     req.ImageURL,
	}

	if err := h.companyService.UpdateCompany(c.Request.Context(), id, company); err != nil {
		h.logger.Printf("failed to update company: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "company updated successfully"})
}

func (h *Handler) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.companyService.DeleteCompany(c.Request.Context(), id); err != nil {
		h.logger.Printf("failed to delete company: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete company"})
		return
	}

	c.Status(http.StatusNoContent)
}
