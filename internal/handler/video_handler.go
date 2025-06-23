package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/dto"
	"shadowify/internal/middleware"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	service *service.VideoService
}

func NewVideoHandler(s *service.VideoService) *VideoHandler {
	return &VideoHandler{service: s}
}

func (h *VideoHandler) RegisterRoutes(e *echo.Echo, device *middleware.Device) {
	v := e.Group("/videos")
	v.POST("", h.Create)
	v.GET("/:id", h.GetByID, device.Authenticate)
	v.GET("", h.List)
	v.GET("/categories", h.Categories)
	v.GET("/favorites", h.GetFavoriteVideos, device.Authenticate)
}

func (h *VideoHandler) GetFavoriteVideos(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	var filter model.FavoriteVideoFilter
	if err := c.Bind(&filter); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid filter parameters"))
	}

	videos, total, err := h.service.GetFavoriteVideos(ctx, user.Id, &filter)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.SuccessWithPagination(c, videos, filter.Pagination.WithTotal(total))
}

func (h *VideoHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req dto.CreateVideoRequest
	if err := c.Bind(&req); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid request"))
	}
	video, err := h.service.Create(ctx, &req)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, video.Id)
}

func (h *VideoHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	video, err := h.service.GetById(ctx, id, user.Id)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, video)
}

func (h *VideoHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	var filter model.VideoFilter
	if err := c.Bind(&filter); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid filter parameters"))
	}

	videos, total, err := h.service.List(ctx, &filter)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.SuccessWithPagination(c, videos, filter.Pagination.WithTotal(total))
}

func (h *VideoHandler) Categories(c echo.Context) error {
	ctx := c.Request().Context()
	categories, err := h.service.Categories(ctx)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, categories)
}
