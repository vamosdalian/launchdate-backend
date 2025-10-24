package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testdb")
	os.Setenv("DB_PORT", "5433")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Database.Host != "testdb" {
		t.Errorf("Expected DB host testdb, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Expected DB port 5433, got %d", cfg.Database.Port)
	}
}

func TestDatabaseDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	dsn := cfg.DSN()

	if dsn != expected {
		t.Errorf("Expected DSN %s, got %s", expected, dsn)
	}
}
