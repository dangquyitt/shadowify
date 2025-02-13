package main

import "github.com/labstack/echo/v4"

type VideoBusiness interface {
	FindAll() ([]Video, error)
}

type videoHandler struct {
	business VideoBusiness
}

func NewVideoHandler(business VideoBusiness) *videoHandler {
	return &videoHandler{
		business: business,
	}
}

func (h *videoHandler) FindAll(c echo.Context) error {
	videos, _ := h.business.FindAll()
	return c.JSON(200, videos)
}
