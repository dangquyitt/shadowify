package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type STTHandler struct {
	sttService *service.STTService
}

func NewSTTHandler(sttService *service.STTService) *STTHandler {
	return &STTHandler{
		sttService: sttService,
	}
}

func (h *STTHandler) RegisterRoutes(e *echo.Echo) {
	stt := e.Group("/stt")
	stt.POST("/transcribe", h.TranscribeAudio)
}

func (h *STTHandler) TranscribeAudio(c echo.Context) error {
	var request model.TranscribeInput
	if err := c.Bind(&request); err != nil {
		response.WriteError(c, apperr.NewAppErr("bad_request", "invalid filter parameters"))
	}

	output, err := h.sttService.Transcribe(c.Request().Context(), &request)
	if err != nil {
		return response.WriteError(c, apperr.NewAppErr("internal_error", "failed to transcribe audio").WithCause(err))
	}

	return response.Success(c, output)
}
