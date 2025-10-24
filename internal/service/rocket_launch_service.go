package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type RocketLaunchService struct {
	repo  *repository.RocketLaunchRepository
	cache *CacheService
}

func NewRocketLaunchService(repo *repository.RocketLaunchRepository, cache *CacheService) *RocketLaunchService {
	return &RocketLaunchService{
		repo:  repo,
		cache: cache,
	}
}

func (s *RocketLaunchService) CreateRocketLaunch(ctx context.Context, req *models.CreateRocketLaunchRequest) (*models.RocketLaunch, error) {
	status := req.Status
	if status == "" {
		status = "scheduled"
	}

	rocketLaunch := &models.RocketLaunch{
		Name:         req.Name,
		LaunchDate:   req.LaunchDate,
		RocketID:     req.RocketID,
		LaunchBaseID: req.LaunchBaseID,
		Status:       status,
		Description:  req.Description,
	}

	err := s.repo.Create(rocketLaunch)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return rocketLaunch, nil
}

func (s *RocketLaunchService) GetRocketLaunch(ctx context.Context, id int64) (*models.RocketLaunch, error) {
	cacheKey := fmt.Sprintf("rocket_launch:%d", id)

	// Try to get from cache
	var rocketLaunch models.RocketLaunch
	err := s.cache.Get(ctx, cacheKey, &rocketLaunch)
	if err == nil {
		return &rocketLaunch, nil
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

func (s *RocketLaunchService) ListRocketLaunches(ctx context.Context, status *string, limit, offset int) ([]models.RocketLaunch, error) {
	cacheKey := fmt.Sprintf("rocket_launches:list:%v:%d:%d", status, limit, offset)

	// Try to get from cache
	var rocketLaunches []models.RocketLaunch
	err := s.cache.Get(ctx, cacheKey, &rocketLaunches)
	if err == nil {
		return rocketLaunches, nil
	}

	// Get from database
	rocketLaunches, err = s.repo.List(status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, rocketLaunches, 5*time.Minute)

	return rocketLaunches, nil
}

func (s *RocketLaunchService) UpdateRocketLaunch(ctx context.Context, id int64, rocketLaunch *models.RocketLaunch) error {
	err := s.repo.Update(id, rocketLaunch)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket_launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return nil
}

func (s *RocketLaunchService) DeleteRocketLaunch(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket_launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return nil
}
