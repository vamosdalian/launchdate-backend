package api

import (
	"testing"
)

// TestSetupRouter_NoPanic verifies that SetupRouter doesn't panic due to route conflicts
func TestSetupRouter_NoPanic(t *testing.T) {
	// Create a handler with nil services for testing router setup
	handler := &Handler{}

	// This should not panic - if it does, the test will fail
	router := SetupRouter(handler)

	if router == nil {
		t.Fatal("Expected router to be created, got nil")
	}

	// Verify all expected routes are registered
	routes := router.Routes()

	// Check that we have routes registered
	if len(routes) == 0 {
		t.Fatal("Expected routes to be registered, got 0")
	}

	// Verify the specific routes we care about
	foundGetLaunch := false
	foundListLaunchTasks := false

	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/v1/launches/:id" {
			foundGetLaunch = true
		}
		if route.Method == "GET" && route.Path == "/api/v1/launches/:id/tasks" {
			foundListLaunchTasks = true
		}
	}

	if !foundGetLaunch {
		t.Error("Expected GET /api/v1/launches/:id route to be registered")
	}

	if !foundListLaunchTasks {
		t.Error("Expected GET /api/v1/launches/:id/tasks route to be registered")
	}
}
