package service

import (
	"testing"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// TestSyncCompanyValidation tests the syncCompany method with various inputs
func TestSyncCompanyValidation(t *testing.T) {
	service := &RocketLaunchService{}

	// Test with nil provider
	err := service.syncCompany(nil)
	if err != nil {
		t.Errorf("Expected no error for nil provider, got: %v", err)
	}

	// Test with provider with ID 0
	provider := &models.RocketLaunchProvider{ID: 0, Name: "Test"}
	err = service.syncCompany(provider)
	if err != nil {
		t.Errorf("Expected no error for provider with ID 0, got: %v", err)
	}
}

// TestSyncLocationValidation tests the syncLocation method with various inputs
func TestSyncLocationValidation(t *testing.T) {
	service := &RocketLaunchService{}

	// Test with nil location
	err := service.syncLocation(nil)
	if err != nil {
		t.Errorf("Expected no error for nil location, got: %v", err)
	}

	// Test with location with ID 0
	location := &models.RocketLaunchPadLocation{ID: 0, Name: "Test"}
	err = service.syncLocation(location)
	if err != nil {
		t.Errorf("Expected no error for location with ID 0, got: %v", err)
	}
}

// TestConvertExternalLaunchWithProvider tests conversion with provider
func TestConvertExternalLaunchWithProvider(t *testing.T) {
	service := &RocketLaunchService{}

	external := &models.ExternalRocketLaunch{
		ID:       123,
		Name:     "Test Launch",
		Provider: &models.RocketLaunchProvider{ID: 1, Name: "SpaceX"},
	}

	result := service.convertExternalLaunch(external)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Name != "Test Launch" {
		t.Errorf("Expected name 'Test Launch', got '%s'", result.Name)
	}

	if result.ExternalID == nil || *result.ExternalID != 123 {
		t.Error("Expected external_id to be set to 123")
	}
}

// TestConvertExternalLaunchWithMissions tests conversion with missions
func TestConvertExternalLaunchWithMissions(t *testing.T) {
	service := &RocketLaunchService{}

	external := &models.ExternalRocketLaunch{
		ID:   456,
		Name: "Test Launch with Missions",
		Missions: []models.RocketLaunchMission{
			{ID: 1, Name: "Mission 1", Description: "First mission"},
			{ID: 2, Name: "Mission 2", Description: "Second mission"},
		},
	}

	result := service.convertExternalLaunch(external)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Note: missions are not converted in convertExternalLaunch
	// They are handled separately in SyncLaunchesFromAPI
	if result.Name != "Test Launch with Missions" {
		t.Errorf("Expected name 'Test Launch with Missions', got '%s'", result.Name)
	}
}
