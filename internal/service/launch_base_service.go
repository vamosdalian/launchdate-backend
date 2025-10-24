package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type LaunchBaseService struct {
	repo  *repository.LaunchBaseRepository
	cache *CacheService
}

func NewLaunchBaseService(repo *repository.LaunchBaseRepository, cache *CacheService) *LaunchBaseService {
	return &LaunchBaseService{
		repo:  repo,
		cache: cache,
	}
}

func (s *LaunchBaseService) CreateLaunchBase(ctx context.Context, req *models.CreateLaunchBaseRequest) (*models.LaunchBase, error) {
	launchBase := &models.LaunchBase{
		Name:        req.Name,
		Location:    req.Location,
		Country:     req.Country,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	err := s.repo.Create(launchBase)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "launch_bases:*")

	return launchBase, nil
}

func (s *LaunchBaseService) GetLaunchBase(ctx context.Context, id int64) (*models.LaunchBase, error) {
	cacheKey := fmt.Sprintf("launch_base:%d", id)

	// Try to get from cache
	var launchBase models.LaunchBase
	err := s.cache.Get(ctx, cacheKey, &launchBase)
	if err == nil {
		return &launchBase, nil
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

func (s *LaunchBaseService) ListLaunchBases(ctx context.Context, limit, offset int) ([]models.LaunchBase, error) {
	cacheKey := fmt.Sprintf("launch_bases:list:%d:%d", limit, offset)

	// Try to get from cache
	var launchBases []models.LaunchBase
	err := s.cache.Get(ctx, cacheKey, &launchBases)
	if err == nil {
		return launchBases, nil
	}

	// Get from database
	launchBases, err = s.repo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, launchBases, 5*time.Minute)

	return launchBases, nil
}

func (s *LaunchBaseService) UpdateLaunchBase(ctx context.Context, id int64, launchBase *models.LaunchBase) error {
	err := s.repo.Update(id, launchBase)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("launch_base:%d", id))
	_ = s.cache.DeletePattern(ctx, "launch_bases:*")

	return nil
}

func (s *LaunchBaseService) DeleteLaunchBase(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("launch_base:%d", id))
	_ = s.cache.DeletePattern(ctx, "launch_bases:*")

	return nil
}
