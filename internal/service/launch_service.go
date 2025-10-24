package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

// LaunchService handles business logic for launches
type LaunchService struct {
	repo  *repository.LaunchRepository
	cache *CacheService
}

// NewLaunchService creates a new launch service
func NewLaunchService(repo *repository.LaunchRepository, cache *CacheService) *LaunchService {
	return &LaunchService{
		repo:  repo,
		cache: cache,
	}
}

// CreateLaunch creates a new launch
func (s *LaunchService) CreateLaunch(ctx context.Context, req *models.CreateLaunchRequest, ownerID int64) (*models.Launch, error) {
	// Set defaults
	if req.Status == "" {
		req.Status = "draft"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	launch := &models.Launch{
		Title:       req.Title,
		Description: req.Description,
		LaunchDate:  req.LaunchDate,
		Status:      req.Status,
		Priority:    req.Priority,
		OwnerID:     ownerID,
		TeamID:      req.TeamID,
		ImageURL:    req.ImageURL,
		Tags:        req.Tags,
	}

	if err := s.repo.Create(ctx, launch); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.DeletePattern(ctx, "launches:*")

	return launch, nil
}

// GetLaunch retrieves a launch by ID
func (s *LaunchService) GetLaunch(ctx context.Context, id int64) (*models.Launch, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf("launch:%d", id)
	var launch models.Launch
	err := s.cache.Get(ctx, cacheKey, &launch)
	if err == nil {
		return &launch, nil
	}

	// Get from database
	launchPtr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, launchPtr, 10*time.Minute)

	return launchPtr, nil
}

// ListLaunches retrieves launches with filters
func (s *LaunchService) ListLaunches(ctx context.Context, status, priority string, teamID *int64, limit, offset int) ([]*models.Launch, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf("launches:status:%s:priority:%s:team:%v:limit:%d:offset:%d", status, priority, teamID, limit, offset)
	var launches []*models.Launch
	err := s.cache.Get(ctx, cacheKey, &launches)
	if err == nil {
		return launches, nil
	}

	// Get from database
	launches, err = s.repo.List(ctx, status, priority, teamID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Store in cache
	_ = s.cache.Set(ctx, cacheKey, launches, 5*time.Minute)

	return launches, nil
}

// UpdateLaunch updates a launch
func (s *LaunchService) UpdateLaunch(ctx context.Context, id int64, req *models.UpdateLaunchRequest) error {
	if err := s.repo.Update(ctx, id, req); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "launches:*")

	return nil
}

// DeleteLaunch deletes a launch
func (s *LaunchService) DeleteLaunch(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "launches:*")

	return nil
}
