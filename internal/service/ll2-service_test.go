package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/vamosdalian/launchdate-backend/internal/config"
)

func TestLoadLaunches(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Read sample data from file
		sampleData, err := os.ReadFile(filepath.Join("testdata", "sample.json"))
		if err != nil {
			t.Fatalf("Failed to read sample.json: %v", err)
		}
		// Send response to be tested
		rw.Write(sampleData)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use server.URL as the base URL for the service
	s := NewLL2Service(&config.Config{LL2URLPrefix: server.URL}, nil)

	launches, err := s.LoadLaunches(1, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(launches.Results) != 1 {
		t.Fatalf("Expected 1 launch, got %d", len(launches.Results))
	}

	expectedID := "eed1132a-d5aa-4c9c-bc38-c8ccb98829b6"
	if launches.Results[0].ID != expectedID {
		t.Fatalf("Expected launch ID %s, got %s", expectedID, launches.Results[0].ID)
	}
}
