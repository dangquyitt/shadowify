package handler

import (
	"github.com/dangquyitt/shadowify/video-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type videoHandler struct {
	service service.VideoService
}

func NewVideoHandler(service service.VideoService) *videoHandler {
	return &videoHandler{
		service: service,
	}
}

func (h *videoHandler) GetTranscript(c echo.Context) error {
	videoURL := c.Param("videoURL")
	transcript, err := h.service.GetTranscript(c.Request().Context(), videoURL)
	if err != nil {
		log.Error(err)
		return c.JSON(500, "Error while getting transcript")
	}
	return c.JSON(200, transcript)
}

func (h *videoHandler) FindAll(c echo.Context) error {
	videos, _ := h.service.FindAll(c.Request().Context())
	return c.JSON(200, videos)
}
