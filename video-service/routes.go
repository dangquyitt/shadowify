package main

import "github.com/labstack/echo/v4"

func InitVideoRoutes(e *echo.Echo) {
	repository := NewVideoRepository()
	business := NewVideoBusiness(repository)
	handler := NewVideoHandler(business)
	e.GET("/videos", handler.FindAll)
}
