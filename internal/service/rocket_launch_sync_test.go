package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vamosdalian/launchdate-backend/internal/models"
)

// TestConvertExternalLaunchSetsExternalID verifies that external_id is properly set
func TestConvertExternalLaunchSetsExternalID(t *testing.T) {
	service := &RocketLaunchService{}

	externalID := int64(12345)
	external := &models.ExternalRocketLaunch{
		ID:   externalID,
		Name: "Test Launch",
		Slug: "test-launch",
		T0:   &models.FlexibleTime{Time: time.Now()},
	}

	result := service.convertExternalLaunch(external)

	assert.NotNil(t, result.ExternalID, "ExternalID should be set")
	assert.Equal(t, externalID, *result.ExternalID, "ExternalID should match the external API ID")
}

// TestConvertExternalLaunchWithDifferentIDs tests multiple external IDs are preserved
func TestConvertExternalLaunchWithDifferentIDs(t *testing.T) {
	service := &RocketLaunchService{}

	testCases := []struct {
		name       string
		externalID int64
	}{
		{"Launch 1", 1},
		{"Launch 2", 2},
		{"Launch 12345", 12345},
		{"Launch 99999", 99999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			external := &models.ExternalRocketLaunch{
				ID:   tc.externalID,
				Name: tc.name,
				Slug: "test-launch",
				T0:   &models.FlexibleTime{Time: time.Now()},
			}

			result := service.convertExternalLaunch(external)

			assert.NotNil(t, result.ExternalID, "ExternalID should be set for %s", tc.name)
			assert.Equal(t, tc.externalID, *result.ExternalID, "ExternalID should match for %s", tc.name)
		})
	}
}

// TestExternalIDFieldPresenceInModel verifies the field exists in the model
func TestExternalIDFieldPresenceInModel(t *testing.T) {
	externalID := int64(12345)
	launch := &models.RocketLaunch{
		ExternalID: &externalID,
		Name:       "Test Launch",
	}

	assert.NotNil(t, launch.ExternalID, "ExternalID field should exist in RocketLaunch model")
	assert.Equal(t, externalID, *launch.ExternalID, "ExternalID value should be preserved")
}

