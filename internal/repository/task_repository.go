package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/database"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *database.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *database.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (launch_id, title, description, assignee_id, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		task.LaunchID,
		task.Title,
		task.Description,
		task.AssigneeID,
		task.Status,
		task.Priority,
		task.DueDate,
		now,
		now,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `
		SELECT id, launch_id, title, description, assignee_id, status, priority, due_date, created_at, updated_at, deleted_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL`

	var task models.Task
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// ListByLaunchID retrieves all tasks for a launch
func (r *TaskRepository) ListByLaunchID(ctx context.Context, launchID int64) ([]*models.Task, error) {
	query := `
		SELECT id, launch_id, title, description, assignee_id, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE launch_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	var tasks []*models.Task
	err := r.db.SelectContext(ctx, &tasks, query, launchID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(ctx context.Context, id int64, req *models.UpdateTaskRequest) error {
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", argPos))
		args = append(args, *req.Title)
		argPos++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *req.Description)
		argPos++
	}

	if req.AssigneeID != nil {
		updates = append(updates, fmt.Sprintf("assignee_id = $%d", argPos))
		args = append(args, *req.AssigneeID)
		argPos++
	}

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	if req.Priority != nil {
		updates = append(updates, fmt.Sprintf("priority = $%d", argPos))
		args = append(args, *req.Priority)
		argPos++
	}

	if req.DueDate != nil {
		updates = append(updates, fmt.Sprintf("due_date = $%d", argPos))
		args = append(args, *req.DueDate)
		argPos++
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d AND deleted_at IS NULL", strings.Join(updates, ", "), argPos)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// Delete soft deletes a task
func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE tasks SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
