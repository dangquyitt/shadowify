package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/middleware"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type SentenceHandler struct {
	sentenceService *service.SentenceService
}

func NewSentenceHandler(sentenceService *service.SentenceService) *SentenceHandler {
	return &SentenceHandler{sentenceService: sentenceService}
}

func (h *SentenceHandler) RegisterRoutes(e *echo.Echo, device *middleware.Device) {
	sentences := e.Group("/sentences")
	sentences.GET("/segments/:segmentId", h.GetBySegmentId, device.Authenticate)
	sentences.DELETE("/segments/:segmentId", h.DeleteBySegmentId, device.Authenticate)
	sentences.POST("", h.Create, device.Authenticate)
	sentences.GET("", h.List, device.Authenticate)
}

func (h *SentenceHandler) GetBySegmentId(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}

	segmentId := c.Param("segmentId")
	sentence, err := h.sentenceService.GetByUserIdAndSegmentId(ctx, user.Id, segmentId)
	if err != nil {
		return response.WriteError(c, err)
	}
	if sentence == nil {
		return response.WriteError(c, apperr.NewAppErr("not_found", "Sentence not found"))
	}

	return response.Success(c, sentence)
}

func (h *SentenceHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}

	var req model.SentenceCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Invalid request format"))
	}

	sentence := &model.Sentence{
		UserId:    user.Id,
		SegmentId: req.SegmentId,
		MeaningEN: req.MeaningEN,
	}

	if err := h.sentenceService.Create(sentence); err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, sentence.Id)
}

func (h *SentenceHandler) DeleteBySegmentId(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}

	segmentId := c.Param("segmentId")
	if err := h.sentenceService.DeleteByUserIdAndSegmentId(ctx, user.Id, segmentId); err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, nil)
}

func (h *SentenceHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	var filter model.SentenceFilter
	user, ok := model.FromContext(ctx)
	if !ok {
		return response.WriteError(c, apperr.NewAppErr("unauthorized", "User not authenticated"))
	}
	if err := c.Bind(&filter); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid filter parameters"))
	}
	filter.UserId = user.Id

	sentences, total, err := h.sentenceService.List(ctx, &filter)
	if err != nil {
		return response.WriteError(c, err)
	}
	return response.SuccessWithPagination(c, sentences, filter.Pagination.WithTotal(total))
}
