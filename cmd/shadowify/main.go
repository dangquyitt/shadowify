package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"shadowify/internal/config"
	"shadowify/internal/database"
	"shadowify/internal/logger"
	"shadowify/internal/server"
	"shadowify/internal/transcript/repository"
	videohttp "shadowify/internal/video/delivery/http"
	vdrepo "shadowify/internal/video/repository"
	"shadowify/internal/video/service"
	extractor "shadowify/proto"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Errorf("Failed to connect to gRPC server: %v", err)
		return
	}
	client := extractor.NewExtractorServiceClient(conn)

	segmentRepo := repository.NewSegmentRepository(db)
	extractorRepo := vdrepo.NewExtractorRepository(client)

	// Setup service dependencies (use nil for repository and grpc client for now)
	videoRepository := vdrepo.NewVideoRepository(db)
	videoService := service.NewVideoService(videoRepository, segmentRepo, extractorRepo)
	videoHandler := videohttp.NewVideoHandler(videoService)

	// Prepare route registrars
	videoRegistrar := func(e *echo.Echo) {
		videoHandler.RegisterRoutes(e)
		// Optionally add root or health check endpoint
		e.GET("/", func(c echo.Context) error {
			time.Sleep(5 * time.Second)
			return c.JSON(http.StatusOK, "OK")
		})
	}

	// Start HTTP server using the new server package
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	server := server.New(":8080", videoRegistrar)
	if err := server.Start(ctx); err != nil {
		logger.Errorf("Server error: %v", err)
	}
}
