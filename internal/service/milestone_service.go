package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

// MilestoneService handles business logic for milestones
type MilestoneService struct {
	repo  *repository.MilestoneRepository
	cache *CacheService
}

// NewMilestoneService creates a new milestone service
func NewMilestoneService(repo *repository.MilestoneRepository, cache *CacheService) *MilestoneService {
	return &MilestoneService{
		repo:  repo,
		cache: cache,
	}
}

// CreateMilestone creates a new milestone
func (s *MilestoneService) CreateMilestone(ctx context.Context, req *models.CreateMilestoneRequest) (*models.Milestone, error) {
	if req.Status == "" {
		req.Status = "pending"
	}

	milestone := &models.Milestone{
		LaunchID:    req.LaunchID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      req.Status,
		Order:       req.Order,
	}

	if err := s.repo.Create(ctx, milestone); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("milestones:launch:%d:*", req.LaunchID))

	return milestone, nil
}

// GetMilestone retrieves a milestone by ID
func (s *MilestoneService) GetMilestone(ctx context.Context, id int64) (*models.Milestone, error) {
	cacheKey := fmt.Sprintf("milestone:%d", id)
	var milestone models.Milestone
	err := s.cache.Get(ctx, cacheKey, &milestone)
	if err == nil {
		return &milestone, nil
	}

	milestonePtr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, milestonePtr, 10*time.Minute)

	return milestonePtr, nil
}

// ListMilestonesByLaunchID retrieves all milestones for a launch
func (s *MilestoneService) ListMilestonesByLaunchID(ctx context.Context, launchID int64) ([]*models.Milestone, error) {
	cacheKey := fmt.Sprintf("milestones:launch:%d", launchID)
	var milestones []*models.Milestone
	err := s.cache.Get(ctx, cacheKey, &milestones)
	if err == nil {
		return milestones, nil
	}

	milestones, err = s.repo.ListByLaunchID(ctx, launchID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, milestones, 5*time.Minute)

	return milestones, nil
}

// UpdateMilestone updates a milestone
func (s *MilestoneService) UpdateMilestone(ctx context.Context, id int64, req *models.UpdateMilestoneRequest) error {
	milestone, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("milestone:%d", id))
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("milestones:launch:%d:*", milestone.LaunchID))

	return nil
}

// DeleteMilestone deletes a milestone
func (s *MilestoneService) DeleteMilestone(ctx context.Context, id int64) error {
	milestone, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("milestone:%d", id))
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("milestones:launch:%d:*", milestone.LaunchID))

	return nil
}
