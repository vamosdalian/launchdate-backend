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
	"github.com/vamosdalian/launchdate-backend/internal/db"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	db, cleandb, err := db.NewMongoDB(cfg.MongodbURL, cfg.MongodbDatabase)
	if err != nil {
		logger.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer cleandb()
	logger.Infof("create mongodb database: %s", cfg.MongodbDatabase)

	handler := api.NewHandler(logger, cfg, db)
	router := api.SetupRouter(handler)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
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
