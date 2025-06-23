package middleware

import (
	"shadowify/internal/logger"
	"shadowify/internal/model"

	"github.com/labstack/echo/v4"
)

type Device struct {
}

func NewDevice() *Device {
	return &Device{}
}

func (d *Device) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Get("X-Device-ID") // Example header to identify the device
		// Here you would typically check for a device token or similar
		// For now, we will just set a dummy device ID in the context
		ctx := model.NewContext(c.Request().Context(), &model.User{
			Id: c.Request().Header.Get("X-Device-ID"), // Replace with actual device ID logic
		})
		logger.Infof("Device ID: %s", c.Request().Header.Get("X-Device-ID"))
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
