package service

import (
	"encoding/json"
	"testing"

	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// TestExternalRocketLaunchJSON tests that we can unmarshal JSON from the API
// with the problematic time format (missing seconds)
func TestExternalRocketLaunchJSON(t *testing.T) {
	// This is an example JSON that mimics what the RocketLaunch.Live API returns
	// Note the time format "2025-10-26T14:05Z" which is missing seconds
	jsonData := `{
		"id": 12345,
		"cospar_id": "2025-001A",
		"sort_date": "20251026",
		"name": "Test Launch",
		"mission_description": "Test mission",
		"launch_description": "Test launch",
		"win_open": "2025-10-26T14:00Z",
		"t0": "2025-10-26T14:05Z",
		"win_close": "2025-10-26T14:10Z",
		"date_str": "Oct 26, 2025",
		"slug": "test-launch-2025",
		"weather_summary": "Clear skies",
		"weather_temp": 22.5,
		"weather_condition": "Sunny",
		"weather_wind_mph": 10.0,
		"weather_icon": "sun",
		"weather_updated": "2025-10-26T12:00Z",
		"quicktext": "Quick test",
		"suborbital": false,
		"modified": "2025-10-26T13:30Z"
	}`

	var launch models.ExternalRocketLaunch
	err := json.Unmarshal([]byte(jsonData), &launch)

	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the launch was parsed correctly
	if launch.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", launch.ID)
	}

	if launch.Name != "Test Launch" {
		t.Errorf("Expected Name 'Test Launch', got %s", launch.Name)
	}

	// Verify that the timestamps were parsed (not nil)
	if launch.T0 == nil {
		t.Error("Expected T0 to be parsed, but it's nil")
	}

	if launch.WindowOpen == nil {
		t.Error("Expected WindowOpen to be parsed, but it's nil")
	}

	if launch.WindowClose == nil {
		t.Error("Expected WindowClose to be parsed, but it's nil")
	}

	if launch.Modified == nil {
		t.Error("Expected Modified to be parsed, but it's nil")
	}

	if launch.WeatherUpdated == nil {
		t.Error("Expected WeatherUpdated to be parsed, but it's nil")
	}

	// Verify we can convert to internal model
	service := &RocketLaunchService{}
	internalLaunch := service.convertExternalLaunch(&launch)

	if internalLaunch == nil {
		t.Fatal("Expected non-nil internal launch")
	}

	if internalLaunch.Name != launch.Name {
		t.Errorf("Expected Name %s, got %s", launch.Name, internalLaunch.Name)
	}

	if internalLaunch.T0 == nil {
		t.Error("Expected T0 to be converted, but it's nil")
	}
}

// TestExternalRocketLaunchResponseJSON tests the full response structure
func TestExternalRocketLaunchResponseJSON(t *testing.T) {
	// Simulates the response from https://fdo.rocketlaunch.live/json/launches/next/5
	jsonData := `{
		"result": [
			{
				"id": 1,
				"name": "Launch 1",
				"t0": "2025-10-26T14:05Z",
				"slug": "launch-1"
			},
			{
				"id": 2,
				"name": "Launch 2",
				"t0": "2025-10-27T10:30Z",
				"slug": "launch-2"
			}
		]
	}`

	var response models.ExternalRocketLaunchResponse
	err := json.Unmarshal([]byte(jsonData), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if len(response.Result) != 2 {
		t.Errorf("Expected 2 launches, got %d", len(response.Result))
	}

	// Verify first launch
	if response.Result[0].ID != 1 {
		t.Errorf("Expected first launch ID 1, got %d", response.Result[0].ID)
	}

	if response.Result[0].T0 == nil {
		t.Error("Expected first launch T0 to be parsed, but it's nil")
	}

	// Verify second launch
	if response.Result[1].ID != 2 {
		t.Errorf("Expected second launch ID 2, got %d", response.Result[1].ID)
	}

	if response.Result[1].T0 == nil {
		t.Error("Expected second launch T0 to be parsed, but it's nil")
	}
}
