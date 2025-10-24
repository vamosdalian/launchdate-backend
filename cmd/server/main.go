package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/api"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/database"
	"github.com/vamosdalian/launchdate-backend/internal/service"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize cache service
	cacheService, err := service.NewCacheService(&cfg.Redis)
	if err != nil {
		logger.Fatalf("failed to connect to redis: %v", err)
	}
	defer cacheService.Close()

	// Initialize handlers
	handler := api.NewHandler(db, cacheService, logger)

	// Setup router
	router := api.SetupRouter(handler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("server starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Info("server stopped")
}
