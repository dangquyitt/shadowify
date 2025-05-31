package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/model"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type LanguageHandler struct {
	languageService *service.LanguageService
}

func NewLanguageHandler(languageService *service.LanguageService) *LanguageHandler {
	return &LanguageHandler{
		languageService: languageService,
	}
}

// RegisterRoutes registers routes for language API
func (h *LanguageHandler) RegisterRoutes(e *echo.Echo) {
	languages := e.Group("/languages")
	languages.POST("", h.CreateLanguage)
	languages.GET("", h.GetAllLanguages)
	languages.GET("/:id", h.GetLanguageByID)
	languages.PUT("/:id", h.UpdateLanguage)
	languages.DELETE("/:id", h.DeleteLanguage)
}

// CreateLanguage creates a new language
// @Summary Create a new language
// @Description Create a new language with the provided details
// @Tags Languages
// @Accept json
// @Produce json
// @Param language body service.CreateLanguageRequest true "Language details"
// @Success 201 {object} model.Language
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /languages [post]
func (h *LanguageHandler) CreateLanguage(c echo.Context) error {
	ctx := c.Request().Context()

	var req model.CreateLanguageRequest
	if err := c.Bind(&req); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Invalid request format"))
	}

	// Validate request
	if req.Code == "" {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Code is required").WithField("code"))
	}

	if req.Name == "" {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Name is required").WithField("name"))
	}

	language, err := h.languageService.Create(ctx, &req)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, language)
}

// GetAllLanguages gets all languages
// @Summary Get all languages
// @Description Get a list of all languages
// @Tags Languages
// @Produce json
// @Success 200 {array} model.Language
// @Failure 500 {object} response.Response
// @Router /languages [get]
func (h *LanguageHandler) GetAllLanguages(c echo.Context) error {
	ctx := c.Request().Context()

	languages, err := h.languageService.GetAll(ctx)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, languages)
}

// GetLanguageByID gets a language by ID
// @Summary Get a language by ID
// @Description Get a language by its ID
// @Tags Languages
// @Produce json
// @Param id path string true "Language ID"
// @Success 200 {object} model.Language
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /languages/{id} [get]
func (h *LanguageHandler) GetLanguageByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	language, err := h.languageService.GetByID(ctx, id)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, language)
}

// UpdateLanguage updates a language
// @Summary Update a language
// @Description Update a language with the provided details
// @Tags Languages
// @Accept json
// @Produce json
// @Param id path string true "Language ID"
// @Param language body service.UpdateLanguageRequest true "Language details"
// @Success 200 {object} model.Language
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /languages/{id} [put]
func (h *LanguageHandler) UpdateLanguage(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var req model.UpdateLanguageRequest
	if err := c.Bind(&req); err != nil {
		return response.WriteError(c, apperr.NewAppErr("bad_request", "Invalid request format"))
	}

	language, err := h.languageService.Update(ctx, id, &req)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, language)
}

// DeleteLanguage deletes a language
// @Summary Delete a language
// @Description Delete a language by its ID
// @Tags Languages
// @Produce json
// @Param id path string true "Language ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /languages/{id} [delete]
func (h *LanguageHandler) DeleteLanguage(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	err := h.languageService.Delete(ctx, id)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, map[string]string{"message": "Language deleted successfully"})
}
