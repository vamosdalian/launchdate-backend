package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type RocketService struct {
	repo  *repository.RocketRepository
	cache *CacheService
}

func NewRocketService(repo *repository.RocketRepository, cache *CacheService) *RocketService {
	return &RocketService{
		repo:  repo,
		cache: cache,
	}
}

func (s *RocketService) CreateRocket(ctx context.Context, req *models.CreateRocketRequest) (*models.Rocket, error) {
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

	err := s.repo.Create(rocket)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "rockets:*")

	return rocket, nil
}

func (s *RocketService) GetRocket(ctx context.Context, id int64) (*models.Rocket, error) {
	cacheKey := fmt.Sprintf("rocket:%d", id)

	// Try to get from cache
	var rocket models.Rocket
	err := s.cache.Get(ctx, cacheKey, &rocket)
	if err == nil {
		return &rocket, nil
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

func (s *RocketService) ListRockets(ctx context.Context, active *bool, limit, offset int) ([]models.Rocket, error) {
	cacheKey := fmt.Sprintf("rockets:list:%v:%d:%d", active, limit, offset)

	// Try to get from cache
	var rockets []models.Rocket
	err := s.cache.Get(ctx, cacheKey, &rockets)
	if err == nil {
		return rockets, nil
	}

	// Get from database
	rockets, err = s.repo.List(active, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, rockets, 5*time.Minute)

	return rockets, nil
}

func (s *RocketService) UpdateRocket(ctx context.Context, id int64, rocket *models.Rocket) error {
	err := s.repo.Update(id, rocket)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket:%d", id))
	_ = s.cache.DeletePattern(ctx, "rockets:*")

	return nil
}

func (s *RocketService) DeleteRocket(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket:%d", id))
	_ = s.cache.DeletePattern(ctx, "rockets:*")

	return nil
}
