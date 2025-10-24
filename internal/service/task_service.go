package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

// TaskService handles business logic for tasks
type TaskService struct {
	repo  *repository.TaskRepository
	cache *CacheService
}

// NewTaskService creates a new task service
func NewTaskService(repo *repository.TaskRepository, cache *CacheService) *TaskService {
	return &TaskService{
		repo:  repo,
		cache: cache,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error) {
	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	task := &models.Task{
		LaunchID:    req.LaunchID,
		MilestoneID: req.MilestoneID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Invalidate cache
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("tasks:launch:%d:*", req.LaunchID))

	return task, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	cacheKey := fmt.Sprintf("task:%d", id)
	var task models.Task
	err := s.cache.Get(ctx, cacheKey, &task)
	if err == nil {
		return &task, nil
	}

	taskPtr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, taskPtr, 10*time.Minute)

	return taskPtr, nil
}

// ListTasksByLaunchID retrieves all tasks for a launch
func (s *TaskService) ListTasksByLaunchID(ctx context.Context, launchID int64, milestoneID *int64) ([]*models.Task, error) {
	cacheKey := fmt.Sprintf("tasks:launch:%d:milestone:%v", launchID, milestoneID)
	var tasks []*models.Task
	err := s.cache.Get(ctx, cacheKey, &tasks)
	if err == nil {
		return tasks, nil
	}

	tasks, err = s.repo.ListByLaunchID(ctx, launchID, milestoneID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, tasks, 5*time.Minute)

	return tasks, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(ctx context.Context, id int64, req *models.UpdateTaskRequest) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("task:%d", id))
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("tasks:launch:%d:*", task.LaunchID))

	return nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	_ = s.cache.Delete(ctx, fmt.Sprintf("task:%d", id))
	_ = s.cache.DeletePattern(ctx, fmt.Sprintf("tasks:launch:%d:*", task.LaunchID))

	return nil
}
