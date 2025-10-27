package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type RocketLaunchService struct {
	repo  *repository.RocketLaunchRepository
	cache *CacheService
}

func NewRocketLaunchService(repo *repository.RocketLaunchRepository, cache *CacheService) *RocketLaunchService {
	return &RocketLaunchService{
		repo:  repo,
		cache: cache,
	}
}

func (s *RocketLaunchService) CreateRocketLaunch(ctx context.Context, req *models.CreateRocketLaunchRequest) (*models.RocketLaunch, error) {
	status := req.Status
	if status == "" {
		status = "scheduled"
	}

	rocketLaunch := &models.RocketLaunch{
		CosparID:           req.CosparID,
		SortDate:           req.SortDate,
		Name:               req.Name,
		LaunchDate:         req.LaunchDate,
		ProviderID:         req.ProviderID,
		RocketID:           req.RocketID,
		LaunchBaseID:       req.LaunchBaseID,
		MissionDescription: req.MissionDescription,
		LaunchDescription:  req.LaunchDescription,
		WindowOpen:         req.WindowOpen,
		T0:                 req.T0,
		WindowClose:        req.WindowClose,
		DateStr:            req.DateStr,
		Slug:               req.Slug,
		WeatherSummary:     req.WeatherSummary,
		WeatherTemp:        req.WeatherTemp,
		WeatherCondition:   req.WeatherCondition,
		WeatherWindMPH:     req.WeatherWindMPH,
		WeatherIcon:        req.WeatherIcon,
		WeatherUpdated:     req.WeatherUpdated,
		QuickText:          req.QuickText,
		Suborbital:         req.Suborbital,
		Modified:           req.Modified,
		Status:             status,
	}

	err := s.repo.Create(rocketLaunch)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return rocketLaunch, nil
}

func (s *RocketLaunchService) GetRocketLaunch(ctx context.Context, id int64) (*models.RocketLaunch, error) {
	cacheKey := fmt.Sprintf("rocket_launch:%d", id)

	// Try to get from cache
	var rocketLaunch models.RocketLaunch
	err := s.cache.Get(ctx, cacheKey, &rocketLaunch)
	if err == nil {
		return &rocketLaunch, nil
	}

	// Get from database
	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, result, 10*time.Minute)

	return result, nil
}

func (s *RocketLaunchService) ListRocketLaunches(ctx context.Context, status *string, limit, offset int) ([]models.RocketLaunch, error) {
	cacheKey := fmt.Sprintf("rocket_launches:list:%v:%d:%d", status, limit, offset)

	// Try to get from cache
	var rocketLaunches []models.RocketLaunch
	err := s.cache.Get(ctx, cacheKey, &rocketLaunches)
	if err == nil {
		return rocketLaunches, nil
	}

	// Get from database
	rocketLaunches, err = s.repo.List(status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, rocketLaunches, 5*time.Minute)

	return rocketLaunches, nil
}

func (s *RocketLaunchService) UpdateRocketLaunch(ctx context.Context, id int64, rocketLaunch *models.RocketLaunch) error {
	err := s.repo.Update(id, rocketLaunch)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket_launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return nil
}

func (s *RocketLaunchService) DeleteRocketLaunch(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = s.cache.Delete(ctx, fmt.Sprintf("rocket_launch:%d", id))
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return nil
}

// SyncLaunchesFromAPI fetches the latest rocket launches from RocketLaunch.Live API and saves them to the database
func (s *RocketLaunchService) SyncLaunchesFromAPI(ctx context.Context) (int, error) {
	// Fetch data from the external API
	resp, err := http.Get("https://fdo.rocketlaunch.live/json/launches/next/5")
	if err != nil {
		return 0, fmt.Errorf("failed to fetch launches from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	// Parse the response
	var apiResponse models.ExternalRocketLaunchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return 0, fmt.Errorf("failed to decode API response: %w", err)
	}

	// Save launches to the database
	savedCount := 0
	for _, extLaunch := range apiResponse.Result {
		// Convert external launch to our model
		rocketLaunch := s.convertExternalLaunch(&extLaunch)

		// Try to create the launch
		err := s.repo.Create(rocketLaunch)
		if err != nil {
			// If creation fails (e.g., duplicate), try to update by slug
			if rocketLaunch.Slug != "" {
				existingLaunch, findErr := s.repo.GetBySlug(rocketLaunch.Slug)
				if findErr == nil && existingLaunch != nil {
					// Update existing launch
					if updateErr := s.repo.Update(existingLaunch.ID, rocketLaunch); updateErr == nil {
						savedCount++
					}
				}
			}
			continue
		}

		savedCount++
	}

	// Invalidate list cache
	_ = s.cache.DeletePattern(ctx, "rocket_launches:*")

	return savedCount, nil
}

// convertExternalLaunch converts an external API launch to our internal model
func (s *RocketLaunchService) convertExternalLaunch(ext *models.ExternalRocketLaunch) *models.RocketLaunch {
	// Use T0 as launch_date if available, otherwise use current time as placeholder
	launchDate := time.Now()
	if ext.T0 != nil {
		t0 := convertFlexibleTime(ext.T0)
		if t0 != nil {
			launchDate = *t0
		}
	}

	launch := &models.RocketLaunch{
		CosparID:           ext.CosparID,
		SortDate:           ext.SortDate,
		Name:               ext.Name,
		LaunchDate:         launchDate,
		MissionDescription: ext.MissionDescription,
		LaunchDescription:  ext.LaunchDescription,
		WindowOpen:         convertFlexibleTime(ext.WindowOpen),
		T0:                 convertFlexibleTime(ext.T0),
		WindowClose:        convertFlexibleTime(ext.WindowClose),
		DateStr:            ext.DateStr,
		Slug:               ext.Slug,
		WeatherSummary:     ext.WeatherSummary,
		WeatherTemp:        ext.WeatherTemp,
		WeatherCondition:   ext.WeatherCondition,
		WeatherWindMPH:     ext.WeatherWindMPH,
		WeatherIcon:        ext.WeatherIcon,
		WeatherUpdated:     convertFlexibleTime(ext.WeatherUpdated),
		QuickText:          ext.QuickText,
		Suborbital:         ext.Suborbital,
		Modified:           convertFlexibleTime(ext.Modified),
		Status:             "scheduled", // Default status
	}

	// Note: We don't set ProviderID, RocketID, or LaunchBaseID because
	// the external API IDs don't match our database IDs
	// These should be linked manually or through a separate matching process

	return launch
}

// convertFlexibleTime converts *FlexibleTime to *time.Time
func convertFlexibleTime(ft *models.FlexibleTime) *time.Time {
	if ft == nil {
		return nil
	}
	return ft.ToTimePtr()
}
