package models

import (
	"testing"
	"time"
)

func TestRocketLaunchModel(t *testing.T) {
	now := time.Now()
	rocketLaunch := &RocketLaunch{
		ID:                1,
		Name:              "Test Rocket Launch",
		Status:            "scheduled",
		LaunchDescription: "Test launch description",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if rocketLaunch.Name != "Test Rocket Launch" {
		t.Errorf("Expected name 'Test Rocket Launch', got '%s'", rocketLaunch.Name)
	}

	if rocketLaunch.Status != "scheduled" {
		t.Errorf("Expected status 'scheduled', got '%s'", rocketLaunch.Status)
	}
}

func TestCreateRocketLaunchRequest(t *testing.T) {
	req := &CreateRocketLaunchRequest{
		Name:              "New Rocket Launch",
		DateStr:           "2024-01-01",
		LaunchDescription: "New launch description",
		Status:            "scheduled",
	}

	if req.Name != "New Rocket Launch" {
		t.Errorf("Expected name 'New Rocket Launch', got '%s'", req.Name)
	}

	if req.Status != "scheduled" {
		t.Errorf("Expected status 'scheduled', got '%s'", req.Status)
	}
}
