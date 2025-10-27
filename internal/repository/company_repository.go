package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type CompanyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *models.Company) error {
	query := `
		INSERT INTO companies (external_id, name, description, founded, founder, headquarters, employees, website, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		company.ExternalID,
		company.Name,
		company.Description,
		company.Founded,
		company.Founder,
		company.Headquarters,
		company.Employees,
		company.Website,
		company.ImageURL,
	).Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
}

func (r *CompanyRepository) GetByID(id int64) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&company, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("company not found")
	}
	return &company, err
}

func (r *CompanyRepository) GetByExternalID(externalID int64) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE external_id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&company, query, externalID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("company not found")
	}
	return &company, err
}

func (r *CompanyRepository) List(limit, offset int) ([]models.Company, error) {
	companies := []models.Company{}
	query := `
		SELECT * FROM companies
		WHERE deleted_at IS NULL
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&companies, query, limit, offset)
	return companies, err
}

func (r *CompanyRepository) Update(id int64, company *models.Company) error {
	query := `
		UPDATE companies
		SET external_id = $1, name = $2, description = $3, founded = $4, founder = $5,
		    headquarters = $6, employees = $7, website = $8, image_url = $9
		WHERE id = $10 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		company.ExternalID,
		company.Name,
		company.Description,
		company.Founded,
		company.Founder,
		company.Headquarters,
		company.Employees,
		company.Website,
		company.ImageURL,
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
		return fmt.Errorf("company not found")
	}

	return nil
}

func (r *CompanyRepository) Delete(id int64) error {
	query := `UPDATE companies SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}
