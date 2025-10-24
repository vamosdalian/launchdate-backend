package models

import (
	"testing"
	"time"
)

func TestLaunchModel(t *testing.T) {
	now := time.Now()
	launch := &Launch{
		ID:          1,
		Title:       "Test Launch",
		Description: "Test Description",
		LaunchDate:  now,
		Status:      "draft",
		Priority:    "high",
		OwnerID:     1,
		Tags:        []string{"test", "launch"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if launch.Title != "Test Launch" {
		t.Errorf("Expected title 'Test Launch', got '%s'", launch.Title)
	}

	if launch.Status != "draft" {
		t.Errorf("Expected status 'draft', got '%s'", launch.Status)
	}

	if len(launch.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(launch.Tags))
	}
}

func TestCreateLaunchRequest(t *testing.T) {
	now := time.Now()
	req := &CreateLaunchRequest{
		Title:       "New Launch",
		Description: "New Description",
		LaunchDate:  now,
		Status:      "planned",
		Priority:    "medium",
		Tags:        []string{"new"},
	}

	if req.Title != "New Launch" {
		t.Errorf("Expected title 'New Launch', got '%s'", req.Title)
	}

	if req.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", req.Priority)
	}
}
