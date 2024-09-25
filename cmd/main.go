package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "music-service/docs" // Import your generated docs
	"music-service/internal/config"
	"music-service/internal/delivery/handler"
	"music-service/internal/delivery/router"
	"music-service/internal/repository"
	"music-service/internal/service"
	"music-service/pkg/database"
	"music-service/pkg/logger"
	"music-service/pkg/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Music Service API
// @version 1.0
// @description This is a music service API.

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.MustLoadConfig()

	loggers, err := logger.SetupLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("Could not set up logger: %v", err)
	}
	loggers.InfoLogger.Info("Starting the service")

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := database.NewDatabase(dbConnStr)
	if err != nil {
		loggers.ErrorLogger.Error("Failed to connect to the database", utils.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := database.CloseDatabase(db); err != nil {
			loggers.ErrorLogger.Error("Failed to close the database", utils.Err(err))
		}
	}()

	if err := database.RunMigrations(db, "./migrations"); err != nil {
		loggers.ErrorLogger.Error("Failed to run migrations", utils.Err(err))
		os.Exit(1)
	}
	loggers.InfoLogger.Info("Migrations applied successfully")

	songRepo := repository.NewSongRepository(db, loggers)
	songService := service.NewSongService(songRepo, loggers)
	songHandler := handler.NewSongHandler(songService, loggers)

	r := router.NewRouter(songHandler)

	// Serve Swagger API documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:         ":" + fmt.Sprint(cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout * time.Second,
		WriteTimeout: cfg.HTTP.Timeout * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		loggers.InfoLogger.Info("HTTP server is starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggers.ErrorLogger.Error("Failed to start HTTP server", utils.Err(err))
			os.Exit(1)
		}
	}()

	gracefulShutdown(srv, loggers)
}

func gracefulShutdown(srv *http.Server, loggers *logger.Loggers) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	loggers.InfoLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		loggers.ErrorLogger.Error("Server forced to shutdown", utils.Err(err))
	} else {
		loggers.InfoLogger.Info("Server gracefully stopped")
	}
}
