package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/dto"
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

func (h *VideoHandler) RegisterRoutes(e *echo.Echo) {
	v := e.Group("/videos")
	v.POST("", h.Create)
	v.GET("/:id", h.GetByID)
	v.GET("", h.List)
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
	video, err := h.service.GetById(ctx, id)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, video)
}

func (h *VideoHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	videos, err := h.service.List(ctx)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, videos)
}
