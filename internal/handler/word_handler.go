package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/middleware"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type WordHandler struct {
	wordService *service.WordService
}

func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}

func (h *WordHandler) RegisterRoutes(e *echo.Echo, auth *middleware.KeycloakMiddleware) {
	words := e.Group("/words")
	words.POST("", h.Create, auth.Authenticate)
	words.GET("", h.List, auth.Authenticate)
	words.DELETE("/:word", h.Delete, auth.Authenticate)
}

func (h *WordHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	var req model.WordCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Invalid request format"))
	}
	word := &model.Word{
		MeaningEN: req.MeaningEN,
		UserId:    user.Id,
		SegmentId: req.SegmentId,
	}
	if err := h.wordService.Create(ctx, word); err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, word.Id)
}

func (h *WordHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	var filter model.WordFilter
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	if err := c.Bind(&filter); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid filter parameters"))
	}
	filter.UserId = user.Id

	videos, total, err := h.wordService.List(ctx, &filter)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.SuccessWithPagination(c, videos, filter.Pagination.WithTotal(total))
}

func (h *WordHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	word := c.Param("word")

	if err := h.wordService.DeleteByWord(ctx, word, user.Id); err != nil {
		return response.WriteError(c, err)
	}
	return response.Success(c, nil)
}
