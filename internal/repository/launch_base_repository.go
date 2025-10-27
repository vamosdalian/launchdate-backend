package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type LaunchBaseRepository struct {
	db *sqlx.DB
}

func NewLaunchBaseRepository(db *sqlx.DB) *LaunchBaseRepository {
	return &LaunchBaseRepository{db: db}
}

func (r *LaunchBaseRepository) Create(launchBase *models.LaunchBase) error {
	query := `
		INSERT INTO launch_bases (external_id, name, location, country, description, image_url, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		launchBase.ExternalID,
		launchBase.Name,
		launchBase.Location,
		launchBase.Country,
		launchBase.Description,
		launchBase.ImageURL,
		launchBase.Latitude,
		launchBase.Longitude,
	).Scan(&launchBase.ID, &launchBase.CreatedAt, &launchBase.UpdatedAt)
}

func (r *LaunchBaseRepository) GetByID(id int64) (*models.LaunchBase, error) {
	var launchBase models.LaunchBase
	query := `SELECT * FROM launch_bases WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&launchBase, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("launch base not found")
	}
	return &launchBase, err
}

func (r *LaunchBaseRepository) GetByExternalID(externalID int64) (*models.LaunchBase, error) {
	var launchBase models.LaunchBase
	query := `SELECT * FROM launch_bases WHERE external_id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&launchBase, query, externalID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("launch base not found")
	}
	return &launchBase, err
}

func (r *LaunchBaseRepository) List(limit, offset int) ([]models.LaunchBase, error) {
	launchBases := []models.LaunchBase{}
	query := `
		SELECT * FROM launch_bases
		WHERE deleted_at IS NULL
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&launchBases, query, limit, offset)
	return launchBases, err
}

func (r *LaunchBaseRepository) Update(id int64, launchBase *models.LaunchBase) error {
	query := `
		UPDATE launch_bases
		SET external_id = $1, name = $2, location = $3, country = $4, description = $5,
		    image_url = $6, latitude = $7, longitude = $8
		WHERE id = $9 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		launchBase.ExternalID,
		launchBase.Name,
		launchBase.Location,
		launchBase.Country,
		launchBase.Description,
		launchBase.ImageURL,
		launchBase.Latitude,
		launchBase.Longitude,
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
		return fmt.Errorf("launch base not found")
	}

	return nil
}

func (r *LaunchBaseRepository) Delete(id int64) error {
	query := `UPDATE launch_bases SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("launch base not found")
	}

	return nil
}
