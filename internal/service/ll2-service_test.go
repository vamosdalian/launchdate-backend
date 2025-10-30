package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoContainer(t *testing.T) (*db.MongoDB, func()) {
	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %s", err)
	}

	mongoDB := &db.MongoDB{
		Client:   mongoClient,
		Database: "testdb",
	}

	cleanup := func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return mongoDB, cleanup
}

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

func TestLoadAgency(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Read sample data from file
		sampleData, err := os.ReadFile(filepath.Join("testdata", "agencies.json"))
		if err != nil {
			t.Fatalf("Failed to read agencies.json: %v", err)
		}
		// Send response to be tested
		rw.Write(sampleData)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use server.URL as the base URL for the service
	s := NewLL2Service(&config.Config{LL2URLPrefix: server.URL}, nil)

	agency, err := s.LoadAngecyFromAPI(1, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedID := 225
	if agency.Results[0].ID != expectedID {
		t.Fatalf("Expected agency ID %d, got %d", expectedID, agency.Results[0].ID)
	}
}

func TestUpdateLaunches(t *testing.T) {
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		sampleData, err := os.ReadFile(filepath.Join("testdata", "sample.json"))
		if err != nil {
			t.Fatalf("Failed to read sample.json: %v", err)
		}
		rw.Write(sampleData)
	}))
	defer server.Close()

	s := NewLL2Service(&config.Config{LL2URLPrefix: server.URL, LL2RequestInterval: 1}, mongoDB)

	err := s.UpdateLaunches(false)
	assert.NoError(t, err)

	var launch models.LL2LaunchNormal
	err = mongoDB.Collection(LL2COLLECTION).FindOne(context.Background(), map[string]any{"id": "eed1132a-d5aa-4c9c-bc38-c8ccb98829b6"}).Decode(&launch)
	assert.NoError(t, err)
	assert.Equal(t, "eed1132a-d5aa-4c9c-bc38-c8ccb98829b6", launch.ID)
}

func TestUpdateAngecy(t *testing.T) {
	mongoDB, cleanup := setupMongoContainer(t)
	defer cleanup()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		sampleData, err := os.ReadFile(filepath.Join("testdata", "agencies.json"))
		if err != nil {
			t.Fatalf("Failed to read agencies.json: %v", err)
		}
		rw.Write(sampleData)
	}))
	defer server.Close()

	s := NewLL2Service(&config.Config{LL2URLPrefix: server.URL, LL2RequestInterval: 1}, mongoDB)

	err := s.UpdateAngecy(false)
	assert.NoError(t, err)

	var agency models.LL2AgencyDetailed
	err = mongoDB.Collection("ll2_agency").FindOne(context.Background(), map[string]any{"id": 225}).Decode(&agency)
	assert.NoError(t, err)
	assert.Equal(t, 225, agency.ID)
}
