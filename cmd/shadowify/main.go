package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"shadowify/internal/config"
	"shadowify/internal/database"
	"shadowify/internal/handler"
	"shadowify/internal/logger"
	"shadowify/internal/repository"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Load config
	cfg, err := config.LoadConfig(fmt.Sprintf("configs/config.%s.yml", env))
	if err != nil {
		stdlog.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger.SetDefaultLogger(logger.NewZerologAdapter(cfg.Logger))
	logger.Infof("App started in %s mode", env)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		stdlog.Fatalf("Failed to connect to database: %v", err)
	}

	segmentRepo := repository.NewSegmentRepository(db)
	extractorRepo, err := repository.NewExtractorRepository("localhost:50051")
	if err != nil {
		logger.Fatalf("Failed to create gRPC client: %v", err)
	}

	// Setup service dependencies (use nil for repository and grpc client for now)
	videoRepository := repository.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepository, segmentRepo, extractorRepo)
	videoHandler := handler.NewVideoHandler(videoService)

	e := echo.New()
	videoHandler.RegisterRoutes(e)

	e.Start(":" + cfg.HTTP.Port)

	// Start HTTP server using the new server package
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	logger.Info("Shutting down server...")

}
