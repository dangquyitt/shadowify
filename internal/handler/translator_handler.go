package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type TranslatorHandler struct {
	translatorService *service.TranslatorService
}

func NewTranslatorHandler(translatorService *service.TranslatorService) *TranslatorHandler {
	return &TranslatorHandler{
		translatorService: translatorService,
	}
}

func (h *TranslatorHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/translate", h.Translate)
}

func (h *TranslatorHandler) Translate(c echo.Context) error {
	var input model.TranslateInput

	if err := c.Bind(&input); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "invalid input parameters").WithCause(err))
	}

	output, err := h.translatorService.Translate(c.Request().Context(), &input)
	if err != nil {
		return response.WriteError(c, apperr.NewAppErr("internal_error", "failed to translate text").WithCause(err))
	}

	return response.Success(c, output)
}
