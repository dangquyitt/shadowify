package middleware

import (
	"context"

	"shadowify/internal/config"
	"shadowify/internal/model"

	"github.com/Nerzal/gocloak/v13"
	"github.com/labstack/echo/v4"
)

type KeycloakMiddleware struct {
	client *gocloak.GoCloak
	cfg    config.KeycloakConfig
}

func NewKeycloakMiddleware(ctx context.Context, cfg config.KeycloakConfig) (*KeycloakMiddleware, error) {
	return &KeycloakMiddleware{
		client: gocloak.NewClient(cfg.Host),
		cfg:    cfg,
	}, nil
}

func (m *KeycloakMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := model.NewContext(c.Request().Context(), &model.User{
			Id: "12345", // This should be replaced with actual user ID from Keycloak
		})
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
