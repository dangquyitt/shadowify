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
	"github.com/labstack/echo/v4/middleware"
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

	db := database.NewGromDatabase(cfg.Database)
	if err != nil {
		stdlog.Fatalf("Failed to connect to database: %v", err)
	}

	segmentRepo := repository.NewSegmentRepository(db)

	// Setup service dependencies (use nil for repository and grpc client for now)
	whisperService := service.NewWhisperService()
	ytDLPService := service.NewYTDLPService()
	videoRepository := repository.NewVideoRepository(db)
	segmentRepository := repository.NewSegmentRepository(db)
	languageRepository := repository.NewLanguageRepository(db)

	// Initialize services
	videoService := service.NewVideoService(videoRepository, segmentRepo, whisperService, ytDLPService)
	segmentService := service.NewSegmentService(segmentRepository)
	sttService := service.NewSTTService(whisperService)
	translatorService := service.NewTranslatorService(cfg.Azure.Translator)

	// Setup handlers
	videoHandler := handler.NewVideoHandler(videoService)
	segmentHandler := handler.NewSegmentHandler(segmentService)
	languageService := service.NewLanguageService(languageRepository)
	languageHandler := handler.NewLanguageHandler(languageService)
	sttHandler := handler.NewSTTHandler(sttService)
	translatorHandler := handler.NewTranslatorHandler(translatorService)

	e := echo.New()
	e.Use(middleware.CORS())
	videoHandler.RegisterRoutes(e)
	segmentHandler.RegisterRoutes(e)
	languageHandler.RegisterRoutes(e)
	sttHandler.RegisterRoutes(e)
	translatorHandler.RegisterRoutes(e)

	e.Start(":" + cfg.HTTP.Port)

	// Start HTTP server using the new server package
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	logger.Info("Shutting down server...")

}
