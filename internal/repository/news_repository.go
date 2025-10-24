package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

type NewsRepository struct {
	db *sqlx.DB
}

func NewNewsRepository(db *sqlx.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) Create(news *models.News) error {
	query := `
		INSERT INTO news (title, summary, content, news_date, url, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		news.Title,
		news.Summary,
		news.Content,
		news.NewsDate,
		news.URL,
		news.ImageURL,
	).Scan(&news.ID, &news.CreatedAt, &news.UpdatedAt)
}

func (r *NewsRepository) GetByID(id int64) (*models.News, error) {
	var news models.News
	query := `SELECT * FROM news WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&news, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("news not found")
	}
	return &news, err
}

func (r *NewsRepository) List(limit, offset int) ([]models.News, error) {
	newsList := []models.News{}
	query := `
		SELECT * FROM news
		WHERE deleted_at IS NULL
		ORDER BY news_date DESC
		LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&newsList, query, limit, offset)
	return newsList, err
}

func (r *NewsRepository) Update(id int64, news *models.News) error {
	query := `
		UPDATE news
		SET title = $1, summary = $2, content = $3, news_date = $4,
		    url = $5, image_url = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		news.Title,
		news.Summary,
		news.Content,
		news.NewsDate,
		news.URL,
		news.ImageURL,
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
		return fmt.Errorf("news not found")
	}

	return nil
}

func (r *NewsRepository) Delete(id int64) error {
	query := `UPDATE news SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("news not found")
	}

	return nil
}
