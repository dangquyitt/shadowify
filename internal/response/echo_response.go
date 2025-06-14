// Package response provides functions to write HTTP responses.
//
// It provides functions to write success and error responses.

package response

import (
	"net/http"

	"shadowify/internal/apperr"
	"shadowify/internal/logger"
	"shadowify/internal/pagination"

	"github.com/labstack/echo/v4"
)

// Success writes a success response with 200 status.
//
// It takes an echo.Context and any data as input and returns an error.
func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, NewSuccessResponse(data))
}

func SuccessWithPagination(c echo.Context, data any, pagination *pagination.Pagination) error {
	return c.JSON(http.StatusOK, NewSuccessResponse(data).WithPagination(pagination))
}

// WriteError writes an error response based on error type.
func WriteError(c echo.Context, err error) error {
	if err == nil {
		return c.NoContent(http.StatusNoContent)
	}
	if appErr, ok := err.(*apperr.AppErr); ok {
		logger.Errorf("error: %s, cause: %s", appErr.Error(), appErr.Unwrap().Error())
		return c.JSON(AppErrCodeToStatus(appErr.Code), NewErrorResponse(appErr))
	}
	logger.Errorf("error: %s", err.Error())
	return c.JSON(AppErrCodeToStatus("bad_request"), NewErrorResponse(
		apperr.NewAppErr("unexpected_error", "An unexpected error occurred"),
	))
}

// AppErrCodeToStatus maps AppErr code to HTTP status code.
func AppErrCodeToStatus(code string) int {
	switch code {
	case "unauthorized":
		return http.StatusUnauthorized
	case "forbidden":
		return http.StatusForbidden
	case "bad_request":
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
