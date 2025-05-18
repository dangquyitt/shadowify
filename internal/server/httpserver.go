package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

// RouteRegistrar is a function that registers routes to the Echo instance.
type RouteRegistrar func(e *echo.Echo)

// HTTPServer wraps the Echo instance and handles graceful shutdown.
type HTTPServer struct {
	echo      *echo.Echo
	registrars []RouteRegistrar
	addr      string
}

// New creates a new HTTPServer with the given address and route registrars.
func New(addr string, registrars ...RouteRegistrar) *HTTPServer {
	e := echo.New()
	return &HTTPServer{
		echo:      e,
		registrars: registrars,
		addr:      addr,
	}
}

// Start runs the HTTP server, registers routes, and handles graceful shutdown.
func (s *HTTPServer) Start(ctx context.Context) error {
	// Register all routes
	for _, reg := range s.registrars {
		reg(s.echo)
	}

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		err := s.echo.Start(s.addr)
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Listen for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case <-ctx.Done():
	case <-sigCh:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.echo.Shutdown(shutdownCtx)
}

// Echo returns the underlying Echo instance for advanced configuration.
func (s *HTTPServer) Echo() *echo.Echo {
	return s.echo
}
