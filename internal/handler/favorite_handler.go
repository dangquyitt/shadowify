package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/middleware"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

func (h *FavoriteHandler) RegisterRoutes(e *echo.Echo, device *middleware.Device) {
	favorites := e.Group("/favorites")
	favorites.Use(device.Authenticate)
	favorites.POST("/:video_id", h.Create)
	favorites.DELETE("/:video_id", h.Delete)
}

func (h *FavoriteHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	userId, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	videoId := c.Param("video_id")

	err := h.favoriteService.Create(ctx, userId.Id, videoId)
	if err != nil {
		return response.WriteError(c, apperr.NewAppErr("internal_error", "Failed to create favorite").WithCause(err))
	}
	return response.Success(c, "Favorite created successfully")
}

func (h *FavoriteHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	userId, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	videoId := c.Param("video_id")

	err := h.favoriteService.Delete(ctx, userId.Id, videoId)
	if err != nil {
		return response.WriteError(c, apperr.NewAppErr("internal_error", "Failed to delete favorite").WithCause(err))
	}
	return response.Success(c, "Favorite deleted successfully")
}
