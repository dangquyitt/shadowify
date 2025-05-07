package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"shadowify/pkg/config"
	"shadowify/pkg/database"
	"shadowify/pkg/logger"
	"time"

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

	_, err = database.NewDatabase(cfg)
	if err != nil {
		stdlog.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.JSON(http.StatusOK, "OK")
	})

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
