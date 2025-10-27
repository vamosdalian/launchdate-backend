package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/repository"
)

type RocketLaunchService struct {
	repo           *repository.RocketLaunchRepository
	companyRepo    *repository.CompanyRepository
	launchBaseRepo *repository.LaunchBaseRepository
	cache          *CacheService
	apiConfig      *config.RocketLaunchAPIConfig
}

func NewRocketLaunchService(
	repo *repository.RocketLaunchRepository,
	companyRepo *repository.CompanyRepository,
	launchBaseRepo *repository.LaunchBaseRepository,
	cache *CacheService,
	apiConfig *config.RocketLaunchAPIConfig,
) *RocketLaunchService {
	return &RocketLaunchService{
		repo:           repo,
		companyRepo:    companyRepo,
		launchBaseRepo: launchBaseRepo,
		cache:          cache,
		apiConfig:      apiConfig,
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

// SyncLaunchesFromAPI fetches rocket launches from RocketLaunch.Live API and saves them to the database
func (s *RocketLaunchService) SyncLaunchesFromAPI(ctx context.Context, limit int) (int, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Build API URL
	apiURL := fmt.Sprintf("%s/launches/next/%d", s.apiConfig.BaseURL, limit)

	// Create HTTP client and request
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if available
	if s.apiConfig.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiConfig.APIKey))
	}

	// Execute request
	resp, err := client.Do(req)
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
		// Sync provider (company) if available
		if extLaunch.Provider != nil && extLaunch.Provider.ID > 0 {
			if err := s.syncCompany(extLaunch.Provider); err != nil {
				// Non-critical error: log and continue processing other launches
				// TODO: Use structured logging (e.g., logrus) when available in service layer
				fmt.Printf("Failed to sync company %d: %v\n", extLaunch.Provider.ID, err)
			}
		}

		// Sync pad location if available
		if extLaunch.Pad != nil && extLaunch.Pad.Location != nil && extLaunch.Pad.Location.ID > 0 {
			if err := s.syncLocation(extLaunch.Pad.Location); err != nil {
				// Non-critical error: log and continue processing other launches
				// TODO: Use structured logging (e.g., logrus) when available in service layer
				fmt.Printf("Failed to sync location %d: %v\n", extLaunch.Pad.Location.ID, err)
			}
		}

		// Convert external launch to our model
		rocketLaunch := s.convertExternalLaunch(&extLaunch)

		// Check if launch with this external_id already exists
		var launchID int64
		if rocketLaunch.ExternalID != nil {
			existingLaunch, findErr := s.repo.GetByExternalID(*rocketLaunch.ExternalID)
			if findErr == nil && existingLaunch != nil {
				// Update existing launch
				if updateErr := s.repo.Update(existingLaunch.ID, rocketLaunch); updateErr == nil {
					launchID = existingLaunch.ID
					savedCount++
				} else {
					continue
				}
			} else {
				// Try to create the launch
				err := s.repo.Create(rocketLaunch)
				if err != nil {
					// If creation fails due to duplicate external_id, skip it
					continue
				}
				launchID = rocketLaunch.ID
				savedCount++
			}
		} else {
			// No external ID, try to create
			err := s.repo.Create(rocketLaunch)
			if err != nil {
				continue
			}
			launchID = rocketLaunch.ID
			savedCount++
		}

		// Sync missions for this launch
		if len(extLaunch.Missions) > 0 {
			missions := make([]models.RocketLaunchMission, 0, len(extLaunch.Missions))
			for _, m := range extLaunch.Missions {
				externalMissionID := int64(m.ID)
				missions = append(missions, models.RocketLaunchMission{
					ExternalID:  &externalMissionID,
					Name:        m.Name,
					Description: m.Description,
				})
			}
			if err := s.repo.SyncMissions(launchID, missions); err != nil {
				// Non-critical error: log and continue
				// TODO: Use structured logging (e.g., logrus) when available in service layer
				fmt.Printf("Failed to sync missions for launch %d: %v\n", launchID, err)
			}
		}
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
		ExternalID:         &ext.ID,
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

// syncCompany syncs a company from the external API to the database
func (s *RocketLaunchService) syncCompany(provider *models.RocketLaunchProvider) error {
	if provider == nil || provider.ID == 0 {
		return nil
	}

	externalID := int64(provider.ID)

	// Check if company already exists with this external_id
	existingCompany, err := s.companyRepo.GetByExternalID(externalID)
	if err == nil && existingCompany != nil {
		// Company exists, update if needed
		existingCompany.Name = provider.Name
		return s.companyRepo.Update(existingCompany.ID, existingCompany)
	}

	// Create new company
	company := &models.Company{
		ExternalID:  &externalID,
		Name:        provider.Name,
		Description: "",
	}

	return s.companyRepo.Create(company)
}

// syncLocation syncs a location (launch base) from the external API to the database
func (s *RocketLaunchService) syncLocation(location *models.RocketLaunchPadLocation) error {
	if location == nil || location.ID == 0 {
		return nil
	}

	externalID := int64(location.ID)

	// Check if launch base already exists with this external_id
	existingBase, err := s.launchBaseRepo.GetByExternalID(externalID)
	if err == nil && existingBase != nil {
		// Launch base exists, update if needed
		existingBase.Name = location.Name
		existingBase.Country = location.Country
		existingBase.Location = fmt.Sprintf("%s, %s", location.State, location.StateName)
		return s.launchBaseRepo.Update(existingBase.ID, existingBase)
	}

	// Create new launch base
	launchBase := &models.LaunchBase{
		ExternalID: &externalID,
		Name:       location.Name,
		Location:   fmt.Sprintf("%s, %s", location.State, location.StateName),
		Country:    location.Country,
	}

	return s.launchBaseRepo.Create(launchBase)
}
