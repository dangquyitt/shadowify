package middleware

import (
	"context"
	"net/http"
	"strings"

	"shadowify/internal/apperr"
	"shadowify/internal/config"
	"shadowify/internal/logger"
	"shadowify/internal/response"

	"github.com/Nerzal/gocloak/v13"
	"github.com/labstack/echo/v4"
)

type KeycloakMiddleware struct {
	client *gocloak.GoCloak
	cfg    *config.KeycloakConfig
}

func NewKeycloakMiddleware(ctx context.Context, cfg *config.KeycloakConfig) (*KeycloakMiddleware, error) {
	return &KeycloakMiddleware{
		client: gocloak.NewClient(cfg.Host),
		cfg:    cfg,
	}, nil
}

func (m *KeycloakMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// Extract Bearer token
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, response.NewErrorResponse(apperr.NewAppErr("unauthorized", "missing or invalid token")))
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		_, err := m.client.RetrospectToken(ctx, token, m.cfg.ClientID, m.cfg.ClientSecret, m.cfg.Realm)
		if err != nil {
			logger.Errorf("failed to validate token: %v", err)
			return c.JSON(http.StatusUnauthorized, response.NewErrorResponse(apperr.NewAppErr("unauthorized", "invalid token")))
		}
		return next(c)
	}
}
