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
		INSERT INTO rocket_launches (
			external_id, cospar_id, sort_date, name, launch_date, provider_id, rocket_id, launch_base_id,
			mission_description, launch_description, window_open, t0, window_close,
			date_str, slug, weather_summary, weather_temp, weather_condition,
			weather_wind_mph, weather_icon, weather_updated, quicktext, suborbital,
			modified, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		rocketLaunch.ExternalID,
		rocketLaunch.CosparID,
		rocketLaunch.SortDate,
		rocketLaunch.Name,
		rocketLaunch.LaunchDate,
		rocketLaunch.ProviderID,
		rocketLaunch.RocketID,
		rocketLaunch.LaunchBaseID,
		rocketLaunch.MissionDescription,
		rocketLaunch.LaunchDescription,
		rocketLaunch.WindowOpen,
		rocketLaunch.T0,
		rocketLaunch.WindowClose,
		rocketLaunch.DateStr,
		rocketLaunch.Slug,
		rocketLaunch.WeatherSummary,
		rocketLaunch.WeatherTemp,
		rocketLaunch.WeatherCondition,
		rocketLaunch.WeatherWindMPH,
		rocketLaunch.WeatherIcon,
		rocketLaunch.WeatherUpdated,
		rocketLaunch.QuickText,
		rocketLaunch.Suborbital,
		rocketLaunch.Modified,
		rocketLaunch.Status,
	).Scan(&rocketLaunch.ID, &rocketLaunch.CreatedAt, &rocketLaunch.UpdatedAt)
}

func (r *RocketLaunchRepository) GetByID(id int64) (*models.RocketLaunch, error) {
	var rocketLaunch models.RocketLaunch
	query := `
		SELECT rl.*
		FROM rocket_launches rl
		WHERE rl.id = $1 AND rl.deleted_at IS NULL
	`
	err := r.db.Get(&rocketLaunch, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rocket launch not found")
	}
	if err != nil {
		return nil, err
	}

	// Load related entities
	if err := r.loadRelatedEntities(&rocketLaunch); err != nil {
		return nil, err
	}

	return &rocketLaunch, nil
}

func (r *RocketLaunchRepository) GetBySlug(slug string) (*models.RocketLaunch, error) {
	var rocketLaunch models.RocketLaunch
	query := `
		SELECT rl.*
		FROM rocket_launches rl
		WHERE rl.slug = $1 AND rl.deleted_at IS NULL
	`
	err := r.db.Get(&rocketLaunch, query, slug)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rocket launch not found")
	}
	if err != nil {
		return nil, err
	}

	// Load related entities
	if err := r.loadRelatedEntities(&rocketLaunch); err != nil {
		return nil, err
	}

	return &rocketLaunch, nil
}

func (r *RocketLaunchRepository) GetByExternalID(externalID int64) (*models.RocketLaunch, error) {
	var rocketLaunch models.RocketLaunch
	query := `
		SELECT rl.*
		FROM rocket_launches rl
		WHERE rl.external_id = $1 AND rl.deleted_at IS NULL
	`
	err := r.db.Get(&rocketLaunch, query, externalID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rocket launch not found")
	}
	if err != nil {
		return nil, err
	}

	// Load related entities
	if err := r.loadRelatedEntities(&rocketLaunch); err != nil {
		return nil, err
	}

	return &rocketLaunch, nil
}

func (r *RocketLaunchRepository) List(status *string, limit, offset int) ([]models.RocketLaunch, error) {
	rocketLaunches := []models.RocketLaunch{}
	query := `
		SELECT rl.*
		FROM rocket_launches rl
		WHERE rl.deleted_at IS NULL
	`
	args := []interface{}{}
	argIndex := 1

	if status != nil {
		query += fmt.Sprintf(" AND rl.status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY COALESCE(rl.t0, rl.window_open, rl.created_at) DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	err := r.db.Select(&rocketLaunches, query, args...)
	if err != nil {
		return nil, err
	}

	// Load related entities for each launch
	for i := range rocketLaunches {
		if err := r.loadRelatedEntities(&rocketLaunches[i]); err != nil {
			return nil, err
		}
	}

	return rocketLaunches, nil
}

// loadRelatedEntities loads provider, vehicle, pad, missions, and tags for a rocket launch
func (r *RocketLaunchRepository) loadRelatedEntities(rl *models.RocketLaunch) error {
	// Load provider
	if rl.ProviderID != nil {
		var provider models.RocketLaunchProvider
		providerQuery := `SELECT id, name, '' as slug FROM companies WHERE id = $1`
		if err := r.db.Get(&provider, providerQuery, *rl.ProviderID); err != nil && err != sql.ErrNoRows {
			return err
		}
		if provider.ID > 0 {
			rl.Provider = &provider
		}
	}

	// Load vehicle (rocket)
	if rl.RocketID != nil {
		var vehicle models.RocketLaunchVehicle
		vehicleQuery := `SELECT id, name, company_id, '' as slug FROM rockets WHERE id = $1`
		if err := r.db.Get(&vehicle, vehicleQuery, *rl.RocketID); err != nil && err != sql.ErrNoRows {
			return err
		}
		if vehicle.ID > 0 {
			rl.Vehicle = &vehicle
		}
	}

	// Load pad (launch base)
	if rl.LaunchBaseID != nil {
		var pad models.RocketLaunchPad
		padQuery := `SELECT id, name FROM launch_bases WHERE id = $1`
		if err := r.db.Get(&pad, padQuery, *rl.LaunchBaseID); err != nil && err != sql.ErrNoRows {
			return err
		}
		if pad.ID > 0 {
			// Load location details for the pad
			var location models.RocketLaunchPadLocation
			locationQuery := `SELECT id, name, country, '' as state, '' as statename, '' as slug FROM launch_bases WHERE id = $1`
			if err := r.db.Get(&location, locationQuery, *rl.LaunchBaseID); err == nil {
				pad.Location = &location
			}
			rl.Pad = &pad
		}
	}

	// Load missions
	missionsQuery := `SELECT id, name, description FROM rocket_launch_missions WHERE rocket_launch_id = $1`
	missions := []models.RocketLaunchMission{}
	if err := r.db.Select(&missions, missionsQuery, rl.ID); err != nil && err != sql.ErrNoRows {
		return err
	}
	rl.Missions = missions

	// Load tags
	tagsQuery := `SELECT id, text FROM rocket_launch_tags WHERE rocket_launch_id = $1`
	tags := []models.RocketLaunchTag{}
	if err := r.db.Select(&tags, tagsQuery, rl.ID); err != nil && err != sql.ErrNoRows {
		return err
	}
	rl.Tags = tags

	return nil
}

func (r *RocketLaunchRepository) Update(id int64, rocketLaunch *models.RocketLaunch) error {
	query := `
		UPDATE rocket_launches
		SET external_id = $1, cospar_id = $2, sort_date = $3, name = $4, launch_date = $5, provider_id = $6, rocket_id = $7, 
		    launch_base_id = $8, mission_description = $9, launch_description = $10,
		    window_open = $11, t0 = $12, window_close = $13, date_str = $14, slug = $15,
		    weather_summary = $16, weather_temp = $17, weather_condition = $18,
		    weather_wind_mph = $19, weather_icon = $20, weather_updated = $21,
		    quicktext = $22, suborbital = $23, modified = $24, status = $25
		WHERE id = $26 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		rocketLaunch.ExternalID,
		rocketLaunch.CosparID,
		rocketLaunch.SortDate,
		rocketLaunch.Name,
		rocketLaunch.LaunchDate,
		rocketLaunch.ProviderID,
		rocketLaunch.RocketID,
		rocketLaunch.LaunchBaseID,
		rocketLaunch.MissionDescription,
		rocketLaunch.LaunchDescription,
		rocketLaunch.WindowOpen,
		rocketLaunch.T0,
		rocketLaunch.WindowClose,
		rocketLaunch.DateStr,
		rocketLaunch.Slug,
		rocketLaunch.WeatherSummary,
		rocketLaunch.WeatherTemp,
		rocketLaunch.WeatherCondition,
		rocketLaunch.WeatherWindMPH,
		rocketLaunch.WeatherIcon,
		rocketLaunch.WeatherUpdated,
		rocketLaunch.QuickText,
		rocketLaunch.Suborbital,
		rocketLaunch.Modified,
		rocketLaunch.Status,
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
