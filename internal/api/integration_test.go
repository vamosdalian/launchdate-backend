package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/database"
	"github.com/vamosdalian/launchdate-backend/internal/models"
	"github.com/vamosdalian/launchdate-backend/internal/service"
)

// testContainer holds the test infrastructure
type testContainer struct {
	postgresContainer *postgres.PostgresContainer
	redisContainer    *redis.RedisContainer
	db                *database.DB
	cache             *service.CacheService
	handler           *Handler
	router            *gin.Engine
	ctx               context.Context
}

// setupTestContainer initializes PostgreSQL and Redis containers for integration testing
func setupTestContainer(t *testing.T) *testContainer {
	ctx := context.Background()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Start PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Start Redis container
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err, "Failed to start Redis container")

	// Get PostgreSQL connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get PostgreSQL connection string")

	// Run migrations
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err, "Failed to get project root")
	migrationsPath := filepath.Join(projectRoot, "migrations")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connStr,
	)
	require.NoError(t, err, "Failed to create migration instance")

	err = m.Up()
	require.NoError(t, err, "Failed to run migrations")
	m.Close()

	// Get container connection details
	pgHost, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	redisHost, err := redisContainer.Host(ctx)
	require.NoError(t, err)
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	// Create database connection
	dbConfig := &config.DatabaseConfig{
		Host:     pgHost,
		Port:     pgPort.Int(),
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	db, err := database.New(dbConfig)
	require.NoError(t, err, "Failed to connect to database")

	// Create Redis connection
	redisConfig := &config.RedisConfig{
		Host:     redisHost,
		Port:     redisPort.Int(),
		Password: "",
		DB:       0,
	}

	cache, err := service.NewCacheService(redisConfig)
	require.NoError(t, err, "Failed to connect to Redis")

	// Create logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create handler and router
	handler := NewHandler(db, cache, logger)
	router := SetupRouter(handler)

	return &testContainer{
		postgresContainer: postgresContainer,
		redisContainer:    redisContainer,
		db:                db,
		cache:             cache,
		handler:           handler,
		router:            router,
		ctx:               ctx,
	}
}

// cleanup terminates the test containers and closes connections
func (tc *testContainer) cleanup(t *testing.T) {
	if tc.cache != nil {
		tc.cache.Close()
	}
	if tc.db != nil {
		tc.db.Close()
	}
	if tc.postgresContainer != nil {
		err := tc.postgresContainer.Terminate(tc.ctx)
		if err != nil {
			t.Logf("Failed to terminate PostgreSQL container: %v", err)
		}
	}
	if tc.redisContainer != nil {
		err := tc.redisContainer.Terminate(tc.ctx)
		if err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	}
}

// Helper function to make HTTP requests
func makeRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	w := makeRequest(t, tc.router, http.MethodGet, "/health", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "ok", response.Database)
	assert.Equal(t, "ok", response.Redis)
}

// TestCompaniesAPI tests all company endpoints
func TestCompaniesAPI(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	// Test Create Company
	t.Run("CreateCompany", func(t *testing.T) {
		company := models.CreateCompanyRequest{
			Name:         "SpaceX",
			Description:  "Space Exploration Technologies Corp.",
			Founded:      2002,
			Founder:      "Elon Musk",
			Headquarters: "Hawthorne, California",
			Employees:    12000,
			Website:      "https://www.spacex.com",
			ImageURL:     "https://example.com/spacex.jpg",
		}

		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
		assert.Equal(t, http.StatusCreated, w.Code)

		var created models.Company
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)

		assert.Equal(t, company.Name, created.Name)
		assert.Equal(t, company.Description, created.Description)
		assert.Equal(t, company.Founded, created.Founded)
		assert.Greater(t, created.ID, int64(0))
	})

	// Test List Companies
	t.Run("ListCompanies", func(t *testing.T) {
		w := makeRequest(t, tc.router, http.MethodGet, "/api/v1/companies", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var companies []models.Company
		err := json.Unmarshal(w.Body.Bytes(), &companies)
		require.NoError(t, err)
		assert.Greater(t, len(companies), 0)
	})

	// Test Get Company
	t.Run("GetCompany", func(t *testing.T) {
		// First create a company
		company := models.CreateCompanyRequest{
			Name:        "Blue Origin",
			Description: "Aerospace manufacturer and spaceflight services company",
			Founded:     2000,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Company
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Get the company
		w := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", created.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var fetched models.Company
		err = json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, company.Name, fetched.Name)
	})

	// Test Update Company
	t.Run("UpdateCompany", func(t *testing.T) {
		// First create a company
		company := models.CreateCompanyRequest{
			Name:        "Rocket Lab",
			Description: "Small satellite launch provider",
			Founded:     2006,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Company
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Update the company
		updated := models.CreateCompanyRequest{
			Name:        "Rocket Lab USA",
			Description: "Updated description",
			Founded:     2006,
			Employees:   500,
		}

		w := makeRequest(t, tc.router, http.MethodPut, fmt.Sprintf("/api/v1/companies/%d", created.ID), updated)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the update
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", created.ID), nil)
		var fetched models.Company
		err = json.Unmarshal(wGet.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, updated.Name, fetched.Name)
		assert.Equal(t, updated.Employees, fetched.Employees)
	})

	// Test Delete Company
	t.Run("DeleteCompany", func(t *testing.T) {
		// First create a company
		company := models.CreateCompanyRequest{
			Name:        "Virgin Galactic",
			Description: "Spaceflight company",
			Founded:     2004,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Company
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Delete the company
		w := makeRequest(t, tc.router, http.MethodDelete, fmt.Sprintf("/api/v1/companies/%d", created.ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", created.ID), nil)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})
}

// TestRocketsAPI tests all rocket endpoints
func TestRocketsAPI(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	// First create a company for the rocket
	company := models.CreateCompanyRequest{
		Name:        "SpaceX",
		Description: "Space Exploration Technologies Corp.",
		Founded:     2002,
	}
	wCompany := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
	require.Equal(t, http.StatusCreated, wCompany.Code)
	var createdCompany models.Company
	err := json.Unmarshal(wCompany.Body.Bytes(), &createdCompany)
	require.NoError(t, err)

	// Test Create Rocket
	t.Run("CreateRocket", func(t *testing.T) {
		rocket := models.CreateRocketRequest{
			Name:        "Falcon 9",
			Description: "Two-stage rocket designed for the reliable and safe transport of satellites",
			Height:      70.0,
			Diameter:    3.7,
			Mass:        549054.0,
			CompanyID:   &createdCompany.ID,
			Active:      true,
			ImageURL:    "https://example.com/falcon9.jpg",
		}

		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rockets", rocket)
		assert.Equal(t, http.StatusCreated, w.Code)

		var created models.Rocket
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)

		assert.Equal(t, rocket.Name, created.Name)
		assert.Equal(t, rocket.Height, created.Height)
		assert.Equal(t, rocket.Active, created.Active)
		assert.Greater(t, created.ID, int64(0))
	})

	// Test List Rockets
	t.Run("ListRockets", func(t *testing.T) {
		w := makeRequest(t, tc.router, http.MethodGet, "/api/v1/rockets", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var rockets []models.Rocket
		err := json.Unmarshal(w.Body.Bytes(), &rockets)
		require.NoError(t, err)
		assert.Greater(t, len(rockets), 0)
	})

	// Test Get Rocket
	t.Run("GetRocket", func(t *testing.T) {
		// First create a rocket
		rocket := models.CreateRocketRequest{
			Name:        "Starship",
			Description: "Fully reusable transportation system",
			Height:      120.0,
			Diameter:    9.0,
			Active:      true,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rockets", rocket)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Rocket
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Get the rocket
		w := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rockets/%d", created.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var fetched models.Rocket
		err = json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, rocket.Name, fetched.Name)
	})

	// Test Update Rocket
	t.Run("UpdateRocket", func(t *testing.T) {
		// First create a rocket
		rocket := models.CreateRocketRequest{
			Name:        "Falcon Heavy",
			Description: "Heavy-lift launch vehicle",
			Height:      70.0,
			Active:      true,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rockets", rocket)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Rocket
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Update the rocket
		updated := models.CreateRocketRequest{
			Name:        "Falcon Heavy (Updated)",
			Description: "Updated description",
			Height:      70.5,
			Active:      true,
		}

		w := makeRequest(t, tc.router, http.MethodPut, fmt.Sprintf("/api/v1/rockets/%d", created.ID), updated)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the update
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rockets/%d", created.ID), nil)
		var fetched models.Rocket
		err = json.Unmarshal(wGet.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, updated.Name, fetched.Name)
		assert.Equal(t, updated.Height, fetched.Height)
	})

	// Test Delete Rocket
	t.Run("DeleteRocket", func(t *testing.T) {
		// First create a rocket
		rocket := models.CreateRocketRequest{
			Name:        "Test Rocket",
			Description: "To be deleted",
			Active:      false,
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rockets", rocket)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.Rocket
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Delete the rocket
		w := makeRequest(t, tc.router, http.MethodDelete, fmt.Sprintf("/api/v1/rockets/%d", created.ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rockets/%d", created.ID), nil)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})
}

// TestLaunchBasesAPI tests all launch base endpoints
func TestLaunchBasesAPI(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	// Test Create Launch Base
	t.Run("CreateLaunchBase", func(t *testing.T) {
		launchBase := models.CreateLaunchBaseRequest{
			Name:        "Cape Canaveral",
			Location:    "Florida",
			Country:     "USA",
			Description: "Primary launch site for American space missions",
			Latitude:    28.3922,
			Longitude:   -80.6077,
			ImageURL:    "https://example.com/cape-canaveral.jpg",
		}

		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/launch-bases", launchBase)
		assert.Equal(t, http.StatusCreated, w.Code)

		var created models.LaunchBase
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)

		assert.Equal(t, launchBase.Name, created.Name)
		assert.Equal(t, launchBase.Country, created.Country)
		assert.Equal(t, launchBase.Latitude, created.Latitude)
		assert.Greater(t, created.ID, int64(0))
	})

	// Test List Launch Bases
	t.Run("ListLaunchBases", func(t *testing.T) {
		w := makeRequest(t, tc.router, http.MethodGet, "/api/v1/launch-bases", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var launchBases []models.LaunchBase
		err := json.Unmarshal(w.Body.Bytes(), &launchBases)
		require.NoError(t, err)
		assert.Greater(t, len(launchBases), 0)
	})

	// Test Get Launch Base
	t.Run("GetLaunchBase", func(t *testing.T) {
		// First create a launch base
		launchBase := models.CreateLaunchBaseRequest{
			Name:        "Kennedy Space Center",
			Location:    "Florida",
			Country:     "USA",
			Description: "NASA's primary launch center",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/launch-bases", launchBase)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.LaunchBase
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Get the launch base
		w := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/launch-bases/%d", created.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var fetched models.LaunchBase
		err = json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, launchBase.Name, fetched.Name)
	})

	// Test Update Launch Base
	t.Run("UpdateLaunchBase", func(t *testing.T) {
		// First create a launch base
		launchBase := models.CreateLaunchBaseRequest{
			Name:        "Vandenberg Space Force Base",
			Location:    "California",
			Country:     "USA",
			Description: "West coast launch site",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/launch-bases", launchBase)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.LaunchBase
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Update the launch base
		updated := models.CreateLaunchBaseRequest{
			Name:        "Vandenberg SFB (Updated)",
			Location:    "California",
			Country:     "United States",
			Description: "Updated description",
			Latitude:    34.7420,
			Longitude:   -120.5724,
		}

		w := makeRequest(t, tc.router, http.MethodPut, fmt.Sprintf("/api/v1/launch-bases/%d", created.ID), updated)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the update
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/launch-bases/%d", created.ID), nil)
		var fetched models.LaunchBase
		err = json.Unmarshal(wGet.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, updated.Name, fetched.Name)
		assert.Equal(t, updated.Country, fetched.Country)
	})

	// Test Delete Launch Base
	t.Run("DeleteLaunchBase", func(t *testing.T) {
		// First create a launch base
		launchBase := models.CreateLaunchBaseRequest{
			Name:     "Test Launch Site",
			Location: "Test Location",
			Country:  "Test",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/launch-bases", launchBase)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.LaunchBase
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Delete the launch base
		w := makeRequest(t, tc.router, http.MethodDelete, fmt.Sprintf("/api/v1/launch-bases/%d", created.ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/launch-bases/%d", created.ID), nil)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})
}

// TestRocketLaunchesAPI tests all rocket launch endpoints
func TestRocketLaunchesAPI(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	// First create dependencies
	company := models.CreateCompanyRequest{
		Name:    "SpaceX",
		Founded: 2002,
	}
	wCompany := makeRequest(t, tc.router, http.MethodPost, "/api/v1/companies", company)
	require.Equal(t, http.StatusCreated, wCompany.Code)
	var createdCompany models.Company
	err := json.Unmarshal(wCompany.Body.Bytes(), &createdCompany)
	require.NoError(t, err)

	rocket := models.CreateRocketRequest{
		Name:      "Falcon 9",
		CompanyID: &createdCompany.ID,
		Active:    true,
	}
	wRocket := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rockets", rocket)
	require.Equal(t, http.StatusCreated, wRocket.Code)
	var createdRocket models.Rocket
	err = json.Unmarshal(wRocket.Body.Bytes(), &createdRocket)
	require.NoError(t, err)

	launchBase := models.CreateLaunchBaseRequest{
		Name:    "Cape Canaveral",
		Country: "USA",
	}
	wBase := makeRequest(t, tc.router, http.MethodPost, "/api/v1/launch-bases", launchBase)
	require.Equal(t, http.StatusCreated, wBase.Code)
	var createdBase models.LaunchBase
	err = json.Unmarshal(wBase.Body.Bytes(), &createdBase)
	require.NoError(t, err)

	// Test Create Rocket Launch
	t.Run("CreateRocketLaunch", func(t *testing.T) {
		launchDate := time.Now().Add(24 * time.Hour)
		t0 := time.Now().Add(24 * time.Hour)
		rocketLaunch := models.CreateRocketLaunchRequest{
			Name:               "Starlink Mission 42",
			LaunchDate:         launchDate,
			DateStr:            "2024-12-31",
			RocketID:           &createdRocket.ID,
			LaunchBaseID:       &createdBase.ID,
			ProviderID:         &createdCompany.ID,
			T0:                 &t0,
			MissionDescription: "Deploy 60 Starlink satellites",
			Status:             "scheduled",
		}

		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rocket-launches", rocketLaunch)
		assert.Equal(t, http.StatusCreated, w.Code)

		var created models.RocketLaunch
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)

		assert.Equal(t, rocketLaunch.Name, created.Name)
		assert.Equal(t, rocketLaunch.Status, created.Status)
		assert.Greater(t, created.ID, int64(0))
	})

	// Test List Rocket Launches
	t.Run("ListRocketLaunches", func(t *testing.T) {
		w := makeRequest(t, tc.router, http.MethodGet, "/api/v1/rocket-launches", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var launches []models.RocketLaunch
		err := json.Unmarshal(w.Body.Bytes(), &launches)
		require.NoError(t, err)
		assert.Greater(t, len(launches), 0)
	})

	// Test Get Rocket Launch
	t.Run("GetRocketLaunch", func(t *testing.T) {
		// First create a rocket launch
		launchDate := time.Now().Add(48 * time.Hour)
		rocketLaunch := models.CreateRocketLaunchRequest{
			Name:       "Test Mission",
			LaunchDate: launchDate,
			DateStr:    "2024-11-01",
			Status:     "scheduled",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rocket-launches", rocketLaunch)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.RocketLaunch
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Get the rocket launch
		w := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rocket-launches/%d", created.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var fetched models.RocketLaunch
		err = json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, rocketLaunch.Name, fetched.Name)
	})

	// Test Update Rocket Launch
	t.Run("UpdateRocketLaunch", func(t *testing.T) {
		// First create a rocket launch
		launchDate := time.Now().Add(72 * time.Hour)
		rocketLaunch := models.CreateRocketLaunchRequest{
			Name:       "Update Test Mission",
			LaunchDate: launchDate,
			DateStr:    "2024-11-15",
			Status:     "scheduled",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rocket-launches", rocketLaunch)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.RocketLaunch
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Update the rocket launch
		updatedLaunchDate := time.Now().Add(96 * time.Hour)
		updated := models.CreateRocketLaunchRequest{
			Name:               "Updated Mission Name",
			LaunchDate:         updatedLaunchDate,
			DateStr:            "2024-11-20",
			Status:             "successful",
			MissionDescription: "Mission completed successfully",
		}

		w := makeRequest(t, tc.router, http.MethodPut, fmt.Sprintf("/api/v1/rocket-launches/%d", created.ID), updated)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the update
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rocket-launches/%d", created.ID), nil)
		var fetched models.RocketLaunch
		err = json.Unmarshal(wGet.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, updated.Name, fetched.Name)
		assert.Equal(t, updated.Status, fetched.Status)
	})

	// Test Delete Rocket Launch
	t.Run("DeleteRocketLaunch", func(t *testing.T) {
		// First create a rocket launch
		launchDate := time.Now().Add(120 * time.Hour)
		rocketLaunch := models.CreateRocketLaunchRequest{
			Name:       "Delete Test Mission",
			LaunchDate: launchDate,
			DateStr:    "2024-12-01",
			Status:     "cancelled",
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rocket-launches", rocketLaunch)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.RocketLaunch
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Delete the rocket launch
		w := makeRequest(t, tc.router, http.MethodDelete, fmt.Sprintf("/api/v1/rocket-launches/%d", created.ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/rocket-launches/%d", created.ID), nil)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})

	// Test Sync Rocket Launches
	t.Run("SyncRocketLaunches", func(t *testing.T) {
		// Call the sync endpoint
		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/rocket-launches/sync", nil)

		// The sync endpoint should return either success or error
		// Since it calls an external API, we check for valid response codes
		assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)

		// If successful, verify the response structure
		if w.Code == http.StatusOK {
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response contains expected fields
			assert.Contains(t, response, "message")
			assert.Contains(t, response, "count")
			assert.Equal(t, "rocket launches synced successfully", response["message"])

			// Count should be a number (can be 0 or more)
			count, ok := response["count"].(float64)
			assert.True(t, ok, "count should be a number")
			assert.GreaterOrEqual(t, count, float64(0))
		}
	})
}

// TestNewsAPI tests all news endpoints
func TestNewsAPI(t *testing.T) {
	tc := setupTestContainer(t)
	defer tc.cleanup(t)

	// Test Create News
	t.Run("CreateNews", func(t *testing.T) {
		news := models.CreateNewsRequest{
			Title:    "SpaceX Successfully Launches Starship",
			Summary:  "SpaceX achieved a major milestone with successful Starship launch",
			Content:  "Full article content here...",
			NewsDate: time.Now(),
			URL:      "https://example.com/spacex-starship-launch",
			ImageURL: "https://example.com/starship.jpg",
		}

		w := makeRequest(t, tc.router, http.MethodPost, "/api/v1/news", news)
		assert.Equal(t, http.StatusCreated, w.Code)

		var created models.News
		err := json.Unmarshal(w.Body.Bytes(), &created)
		require.NoError(t, err)

		assert.Equal(t, news.Title, created.Title)
		assert.Equal(t, news.Summary, created.Summary)
		assert.Greater(t, created.ID, int64(0))
	})

	// Test List News
	t.Run("ListNews", func(t *testing.T) {
		w := makeRequest(t, tc.router, http.MethodGet, "/api/v1/news", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var newsList []models.News
		err := json.Unmarshal(w.Body.Bytes(), &newsList)
		require.NoError(t, err)
		assert.Greater(t, len(newsList), 0)
	})

	// Test Get News
	t.Run("GetNews", func(t *testing.T) {
		// First create news
		news := models.CreateNewsRequest{
			Title:    "New Mars Rover Landed",
			Summary:  "NASA's latest Mars rover successfully landed",
			NewsDate: time.Now(),
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/news", news)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.News
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Get the news
		w := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/news/%d", created.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var fetched models.News
		err = json.Unmarshal(w.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, news.Title, fetched.Title)
	})

	// Test Update News
	t.Run("UpdateNews", func(t *testing.T) {
		// First create news
		news := models.CreateNewsRequest{
			Title:    "Breaking Space News",
			Summary:  "Initial summary",
			NewsDate: time.Now(),
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/news", news)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.News
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Update the news
		updated := models.CreateNewsRequest{
			Title:    "Updated Breaking Space News",
			Summary:  "Updated summary with more details",
			Content:  "Full updated content",
			NewsDate: time.Now(),
		}

		w := makeRequest(t, tc.router, http.MethodPut, fmt.Sprintf("/api/v1/news/%d", created.ID), updated)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify the update
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/news/%d", created.ID), nil)
		var fetched models.News
		err = json.Unmarshal(wGet.Body.Bytes(), &fetched)
		require.NoError(t, err)
		assert.Equal(t, updated.Title, fetched.Title)
		assert.Equal(t, updated.Summary, fetched.Summary)
	})

	// Test Delete News
	t.Run("DeleteNews", func(t *testing.T) {
		// First create news
		news := models.CreateNewsRequest{
			Title:    "Old Space News",
			Summary:  "To be deleted",
			NewsDate: time.Now(),
		}

		wCreate := makeRequest(t, tc.router, http.MethodPost, "/api/v1/news", news)
		require.Equal(t, http.StatusCreated, wCreate.Code)

		var created models.News
		err := json.Unmarshal(wCreate.Body.Bytes(), &created)
		require.NoError(t, err)

		// Delete the news
		w := makeRequest(t, tc.router, http.MethodDelete, fmt.Sprintf("/api/v1/news/%d", created.ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		wGet := makeRequest(t, tc.router, http.MethodGet, fmt.Sprintf("/api/v1/news/%d", created.ID), nil)
		assert.Equal(t, http.StatusNotFound, wGet.Code)
	})
}
