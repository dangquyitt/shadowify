package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dangquyitt/shadowify/video-service/internal/delivery/http/handler"
	"github.com/dangquyitt/shadowify/video-service/internal/repository"
	"github.com/dangquyitt/shadowify/video-service/internal/service"
	pb "github.com/dangquyitt/shadowify/video-service/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const TRANSCRIPTION_SERVICE_ADDRESS = ":50051"

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

	// Start server
	conn, err := grpc.NewClient(TRANSCRIPTION_SERVICE_ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	transcriptionServiceClient := pb.NewTranscriptionServiceClient(conn)
	transcriptionService := service.NewTranscriptionService(transcriptionServiceClient)

	videoRepo := repository.NewVideoRepository()
	videoService := service.NewVideoService(videoRepo, transcriptionService)
	videoHandler := handler.NewVideoHandler(videoService)
	handler.RegisterRoutes(e.Group("/v1"), videoHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
