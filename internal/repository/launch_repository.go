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

// LaunchRepository handles database operations for launches
type LaunchRepository struct {
	db *database.DB
}

// NewLaunchRepository creates a new launch repository
func NewLaunchRepository(db *database.DB) *LaunchRepository {
	return &LaunchRepository{db: db}
}

// Create creates a new launch
func (r *LaunchRepository) Create(ctx context.Context, launch *models.Launch) error {
	query := `
		INSERT INTO launches (title, description, launch_date, status, priority, owner_id, team_id, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		launch.Title,
		launch.Description,
		launch.LaunchDate,
		launch.Status,
		launch.Priority,
		launch.OwnerID,
		launch.TeamID,
		launch.ImageURL,
		now,
		now,
	).Scan(&launch.ID, &launch.CreatedAt, &launch.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create launch: %w", err)
	}

	// Insert tags if any
	if len(launch.Tags) > 0 {
		if err := r.AddTags(ctx, launch.ID, launch.Tags); err != nil {
			return fmt.Errorf("failed to add tags: %w", err)
		}
	}

	return nil
}

// GetByID retrieves a launch by ID
func (r *LaunchRepository) GetByID(ctx context.Context, id int64) (*models.Launch, error) {
	query := `
		SELECT id, title, description, launch_date, status, priority, owner_id, team_id, image_url, created_at, updated_at, deleted_at
		FROM launches
		WHERE id = $1 AND deleted_at IS NULL`

	var launch models.Launch
	err := r.db.GetContext(ctx, &launch, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("launch not found")
		}
		return nil, fmt.Errorf("failed to get launch: %w", err)
	}

	// Load tags
	tags, err := r.GetTags(ctx, id)
	if err != nil {
		return nil, err
	}
	launch.Tags = tags

	return &launch, nil
}

// List retrieves all launches with optional filters
func (r *LaunchRepository) List(ctx context.Context, status, priority string, teamID *int64, limit, offset int) ([]*models.Launch, error) {
	query := `
		SELECT id, title, description, launch_date, status, priority, owner_id, team_id, image_url, created_at, updated_at
		FROM launches
		WHERE deleted_at IS NULL`

	args := []interface{}{}
	argPos := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, status)
		argPos++
	}

	if priority != "" {
		query += fmt.Sprintf(" AND priority = $%d", argPos)
		args = append(args, priority)
		argPos++
	}

	if teamID != nil {
		query += fmt.Sprintf(" AND team_id = $%d", argPos)
		args = append(args, *teamID)
		argPos++
	}

	query += " ORDER BY launch_date DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	var launches []*models.Launch
	err := r.db.SelectContext(ctx, &launches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list launches: %w", err)
	}

	// Load tags for each launch
	for _, launch := range launches {
		tags, err := r.GetTags(ctx, launch.ID)
		if err != nil {
			return nil, err
		}
		launch.Tags = tags
	}

	return launches, nil
}

// Update updates a launch
func (r *LaunchRepository) Update(ctx context.Context, id int64, req *models.UpdateLaunchRequest) error {
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

	if req.LaunchDate != nil {
		updates = append(updates, fmt.Sprintf("launch_date = $%d", argPos))
		args = append(args, *req.LaunchDate)
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

	if req.TeamID != nil {
		updates = append(updates, fmt.Sprintf("team_id = $%d", argPos))
		args = append(args, *req.TeamID)
		argPos++
	}

	if req.ImageURL != nil {
		updates = append(updates, fmt.Sprintf("image_url = $%d", argPos))
		args = append(args, *req.ImageURL)
		argPos++
	}

	if len(updates) == 0 && len(req.Tags) == 0 {
		return nil // Nothing to update
	}

	if len(updates) > 0 {
		query := fmt.Sprintf("UPDATE launches SET %s WHERE id = $%d AND deleted_at IS NULL", strings.Join(updates, ", "), argPos)
		args = append(args, id)

		result, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update launch: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("launch not found")
		}
	}

	// Update tags if provided
	if req.Tags != nil {
		if err := r.RemoveTags(ctx, id); err != nil {
			return err
		}
		if len(req.Tags) > 0 {
			if err := r.AddTags(ctx, id, req.Tags); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete soft deletes a launch
func (r *LaunchRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE launches SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete launch: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("launch not found")
	}

	return nil
}

// AddTags adds tags to a launch
func (r *LaunchRepository) AddTags(ctx context.Context, launchID int64, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	query := `INSERT INTO launch_tags (launch_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	for _, tag := range tags {
		_, err := r.db.ExecContext(ctx, query, launchID, tag)
		if err != nil {
			return fmt.Errorf("failed to add tag: %w", err)
		}
	}

	return nil
}

// RemoveTags removes all tags from a launch
func (r *LaunchRepository) RemoveTags(ctx context.Context, launchID int64) error {
	query := `DELETE FROM launch_tags WHERE launch_id = $1`
	_, err := r.db.ExecContext(ctx, query, launchID)
	if err != nil {
		return fmt.Errorf("failed to remove tags: %w", err)
	}
	return nil
}

// GetTags retrieves tags for a launch
func (r *LaunchRepository) GetTags(ctx context.Context, launchID int64) ([]string, error) {
	query := `SELECT tag FROM launch_tags WHERE launch_id = $1`

	var tags []string
	err := r.db.SelectContext(ctx, &tags, query, launchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}
