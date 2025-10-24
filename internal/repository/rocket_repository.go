package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type RocketRepository struct {
	db *sqlx.DB
}

func NewRocketRepository(db *sqlx.DB) *RocketRepository {
	return &RocketRepository{db: db}
}

func (r *RocketRepository) Create(rocket *models.Rocket) error {
	query := `
		INSERT INTO rockets (name, description, height, diameter, mass, company_id, image_url, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		rocket.Name,
		rocket.Description,
		rocket.Height,
		rocket.Diameter,
		rocket.Mass,
		rocket.CompanyID,
		rocket.ImageURL,
		rocket.Active,
	).Scan(&rocket.ID, &rocket.CreatedAt, &rocket.UpdatedAt)
}

func (r *RocketRepository) GetByID(id int64) (*models.Rocket, error) {
	var rocket models.Rocket
	query := `
		SELECT r.*, c.name as company
		FROM rockets r
		LEFT JOIN companies c ON r.company_id = c.id
		WHERE r.id = $1 AND r.deleted_at IS NULL
	`
	err := r.db.Get(&rocket, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rocket not found")
	}
	return &rocket, err
}

func (r *RocketRepository) List(active *bool, limit, offset int) ([]models.Rocket, error) {
	rockets := []models.Rocket{}
	query := `
		SELECT r.*, c.name as company
		FROM rockets r
		LEFT JOIN companies c ON r.company_id = c.id
		WHERE r.deleted_at IS NULL
	`
	args := []interface{}{}
	argIndex := 1

	if active != nil {
		query += fmt.Sprintf(" AND r.active = $%d", argIndex)
		args = append(args, *active)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY r.name LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	err := r.db.Select(&rockets, query, args...)
	return rockets, err
}

func (r *RocketRepository) Update(id int64, rocket *models.Rocket) error {
	query := `
		UPDATE rockets
		SET name = $1, description = $2, height = $3, diameter = $4,
		    mass = $5, company_id = $6, image_url = $7, active = $8
		WHERE id = $9 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		rocket.Name,
		rocket.Description,
		rocket.Height,
		rocket.Diameter,
		rocket.Mass,
		rocket.CompanyID,
		rocket.ImageURL,
		rocket.Active,
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
		return fmt.Errorf("rocket not found")
	}

	return nil
}

func (r *RocketRepository) Delete(id int64) error {
	query := `UPDATE rockets SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("rocket not found")
	}

	return nil
}
