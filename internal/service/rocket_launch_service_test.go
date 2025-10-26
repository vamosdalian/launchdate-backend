package service

import (
	"testing"
	"time"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

func TestConvertExternalLaunch(t *testing.T) {
	service := &RocketLaunchService{}

	now := time.Now()
	flexTime := &models.FlexibleTime{Time: now}
	temp := float32(75.5)
	wind := float32(10.0)

	external := &models.ExternalRocketLaunch{
		CosparID:           "2025-001A",
		SortDate:           "20251024",
		Name:               "Test Launch",
		MissionDescription: "Test Mission",
		LaunchDescription:  "Test Description",
		WindowOpen:         flexTime,
		T0:                 flexTime,
		WindowClose:        flexTime,
		DateStr:            "Oct 24, 2025",
		Slug:               "test-launch",
		WeatherSummary:     "Clear",
		WeatherTemp:        &temp,
		WeatherCondition:   "Sunny",
		WeatherWindMPH:     &wind,
		WeatherIcon:        "sun",
		WeatherUpdated:     flexTime,
		QuickText:          "Quick text",
		Suborbital:         false,
		Modified:           flexTime,
	}

	result := service.convertExternalLaunch(external)

	// Verify all fields are copied correctly
	if result.CosparID != external.CosparID {
		t.Errorf("Expected CosparID %s, got %s", external.CosparID, result.CosparID)
	}
	if result.SortDate != external.SortDate {
		t.Errorf("Expected SortDate %s, got %s", external.SortDate, result.SortDate)
	}
	if result.Name != external.Name {
		t.Errorf("Expected Name %s, got %s", external.Name, result.Name)
	}
	if result.Slug != external.Slug {
		t.Errorf("Expected Slug %s, got %s", external.Slug, result.Slug)
	}
	if result.Status != "scheduled" {
		t.Errorf("Expected Status 'scheduled', got %s", result.Status)
	}
	if result.Suborbital != external.Suborbital {
		t.Errorf("Expected Suborbital %v, got %v", external.Suborbital, result.Suborbital)
	}

	// Verify that foreign key IDs are nil (not copied from external)
	if result.ProviderID != nil {
		t.Errorf("Expected ProviderID to be nil, got %v", *result.ProviderID)
	}
	if result.RocketID != nil {
		t.Errorf("Expected RocketID to be nil, got %v", *result.RocketID)
	}
	if result.LaunchBaseID != nil {
		t.Errorf("Expected LaunchBaseID to be nil, got %v", *result.LaunchBaseID)
	}
}
