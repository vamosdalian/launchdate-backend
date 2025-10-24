package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type CompanyService struct {
	repo  *repository.CompanyRepository
	cache *CacheService
}

func NewCompanyService(repo *repository.CompanyRepository, cache *CacheService) *CompanyService {
	return &CompanyService{
		repo:  repo,
		cache: cache,
	}
}

func (s *CompanyService) CreateCompany(ctx context.Context, req *models.CreateCompanyRequest) (*models.Company, error) {
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

	err := s.repo.Create(company)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "companies:*")

	return company, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, id int64) (*models.Company, error) {
	cacheKey := fmt.Sprintf("company:%d", id)

	// Try to get from cache
	var company models.Company
	err := s.cache.Get(ctx, cacheKey, &company)
	if err == nil {
		return &company, nil
	}

	// Get from database
	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, result, 10*time.Minute)

	return result, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, limit, offset int) ([]models.Company, error) {
	cacheKey := fmt.Sprintf("companies:list:%d:%d", limit, offset)

	// Try to get from cache
	var companies []models.Company
	err := s.cache.Get(ctx, cacheKey, &companies)
	if err == nil {
		return companies, nil
	}

	// Get from database
	companies, err = s.repo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, companies, 5*time.Minute)

	return companies, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id int64, company *models.Company) error {
	err := s.repo.Update(id, company)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("company:%d", id))
	_ = s.cache.DeletePattern(ctx, "companies:*")

	return nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("company:%d", id))
	_ = s.cache.DeletePattern(ctx, "companies:*")

	return nil
}
