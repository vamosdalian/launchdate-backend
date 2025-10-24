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

// MilestoneRepository handles database operations for milestones
type MilestoneRepository struct {
	db *database.DB
}

// NewMilestoneRepository creates a new milestone repository
func NewMilestoneRepository(db *database.DB) *MilestoneRepository {
	return &MilestoneRepository{db: db}
}

// Create creates a new milestone
func (r *MilestoneRepository) Create(ctx context.Context, milestone *models.Milestone) error {
	query := `
		INSERT INTO milestones (launch_id, title, description, due_date, status, order_num, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		milestone.LaunchID,
		milestone.Title,
		milestone.Description,
		milestone.DueDate,
		milestone.Status,
		milestone.Order,
		now,
		now,
	).Scan(&milestone.ID, &milestone.CreatedAt, &milestone.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	return nil
}

// GetByID retrieves a milestone by ID
func (r *MilestoneRepository) GetByID(ctx context.Context, id int64) (*models.Milestone, error) {
	query := `
		SELECT id, launch_id, title, description, due_date, status, order_num, created_at, updated_at, deleted_at
		FROM milestones
		WHERE id = $1 AND deleted_at IS NULL`

	var milestone models.Milestone
	err := r.db.GetContext(ctx, &milestone, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("milestone not found")
		}
		return nil, fmt.Errorf("failed to get milestone: %w", err)
	}

	return &milestone, nil
}

// ListByLaunchID retrieves all milestones for a launch
func (r *MilestoneRepository) ListByLaunchID(ctx context.Context, launchID int64) ([]*models.Milestone, error) {
	query := `
		SELECT id, launch_id, title, description, due_date, status, order_num, created_at, updated_at
		FROM milestones
		WHERE launch_id = $1 AND deleted_at IS NULL
		ORDER BY order_num ASC, due_date ASC`

	var milestones []*models.Milestone
	err := r.db.SelectContext(ctx, &milestones, query, launchID)
	if err != nil {
		return nil, fmt.Errorf("failed to list milestones: %w", err)
	}

	return milestones, nil
}

// Update updates a milestone
func (r *MilestoneRepository) Update(ctx context.Context, id int64, req *models.UpdateMilestoneRequest) error {
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

	if req.DueDate != nil {
		updates = append(updates, fmt.Sprintf("due_date = $%d", argPos))
		args = append(args, *req.DueDate)
		argPos++
	}

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	if req.Order != nil {
		updates = append(updates, fmt.Sprintf("order_num = $%d", argPos))
		args = append(args, *req.Order)
		argPos++
	}

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE milestones SET %s WHERE id = $%d AND deleted_at IS NULL", strings.Join(updates, ", "), argPos)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}

// Delete soft deletes a milestone
func (r *MilestoneRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE milestones SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("milestone not found")
	}

	return nil
}
