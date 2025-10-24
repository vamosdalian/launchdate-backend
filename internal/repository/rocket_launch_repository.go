package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type RocketLaunchRepository struct {
	db *sqlx.DB
}

func NewRocketLaunchRepository(db *sqlx.DB) *RocketLaunchRepository {
	return &RocketLaunchRepository{db: db}
}

func (r *RocketLaunchRepository) Create(rocketLaunch *models.RocketLaunch) error {
	query := `
		INSERT INTO rocket_launches (name, launch_date, rocket_id, launch_base_id, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		rocketLaunch.Name,
		rocketLaunch.LaunchDate,
		rocketLaunch.RocketID,
		rocketLaunch.LaunchBaseID,
		rocketLaunch.Status,
		rocketLaunch.Description,
	).Scan(&rocketLaunch.ID, &rocketLaunch.CreatedAt, &rocketLaunch.UpdatedAt)
}

func (r *RocketLaunchRepository) GetByID(id int64) (*models.RocketLaunch, error) {
	var rocketLaunch models.RocketLaunch
	query := `
		SELECT rl.*, r.name as rocket, lb.name as launchBase
		FROM rocket_launches rl
		LEFT JOIN rockets r ON rl.rocket_id = r.id
		LEFT JOIN launch_bases lb ON rl.launch_base_id = lb.id
		WHERE rl.id = $1 AND rl.deleted_at IS NULL
	`
	err := r.db.Get(&rocketLaunch, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rocket launch not found")
	}
	return &rocketLaunch, err
}

func (r *RocketLaunchRepository) List(status *string, limit, offset int) ([]models.RocketLaunch, error) {
	rocketLaunches := []models.RocketLaunch{}
	query := `
		SELECT rl.*, r.name as rocket, lb.name as launchBase
		FROM rocket_launches rl
		LEFT JOIN rockets r ON rl.rocket_id = r.id
		LEFT JOIN launch_bases lb ON rl.launch_base_id = lb.id
		WHERE rl.deleted_at IS NULL
	`
	args := []interface{}{}
	argIndex := 1

	if status != nil {
		query += fmt.Sprintf(" AND rl.status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY rl.launch_date DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	err := r.db.Select(&rocketLaunches, query, args...)
	return rocketLaunches, err
}

func (r *RocketLaunchRepository) Update(id int64, rocketLaunch *models.RocketLaunch) error {
	query := `
		UPDATE rocket_launches
		SET name = $1, launch_date = $2, rocket_id = $3, launch_base_id = $4,
		    status = $5, description = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		rocketLaunch.Name,
		rocketLaunch.LaunchDate,
		rocketLaunch.RocketID,
		rocketLaunch.LaunchBaseID,
		rocketLaunch.Status,
		rocketLaunch.Description,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("rocket launch not found")
	}

	return nil
}

func (r *RocketLaunchRepository) Delete(id int64) error {
	query := `UPDATE rocket_launches SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("rocket launch not found")
	}

	return nil
}
