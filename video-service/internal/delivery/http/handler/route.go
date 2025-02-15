package handler

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, hdl *videoHandler) {
	e.GET("/videos", hdl.FindAll)
	e.GET("/videos/:videoURL/transcript", hdl.GetTranscript)
}
