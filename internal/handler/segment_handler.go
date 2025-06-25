package handler

import (
	"shadowify/internal/apperr"
	"shadowify/internal/response"
	"shadowify/internal/service"

	"github.com/labstack/echo/v4"
)

type SegmentHandler struct {
	segmentService *service.SegmentService
}

func NewSegmentHandler(segmentService *service.SegmentService) *SegmentHandler {
	return &SegmentHandler{
		segmentService: segmentService,
	}
}

// GetSegmentsByVideoID godoc
// @Summary Get segments by video ID
// @Description Get all segments for a specific video
// @Tags segments
// @Accept json
// @Produce json
// @Param video_id query string true "Video ID"
// @Success 200 {object} getSegmentsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/segments [get]
func (h *SegmentHandler) GetSegmentsByVideoID(c echo.Context) error {
	videoID := c.Param("video_id")
	if videoID == "" {
		return response.WriteError(c, apperr.NewAppErr("segment.invalid_video_id", "video_id is required"))
	}

	segments, err := h.segmentService.GetSegmentsByVideoID(c.Request().Context(), videoID)
	if err != nil {
		return response.WriteError(c, err)
	}

	return response.Success(c, segments)
}

// RegisterRoutes registers the segment routes to the provided Echo instance
func (h *SegmentHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/videos/:video_id/segments", h.GetSegmentsByVideoID)
	e.GET("/segments/:segment_id", h.GetSegmentByID)
}

func (h *SegmentHandler) GetSegmentByID(c echo.Context) error {
	segmentID := c.Param("segment_id")
	if segmentID == "" {
		return response.WriteError(c, apperr.NewAppErr("segment.invalid_segment_id", "segment_id is required"))
	}

	segment, err := h.segmentService.GetSegmentByID(c.Request().Context(), segmentID)
	if err != nil {
		return response.WriteError(c, err)
	}
	if segment == nil {
		return response.WriteError(c, apperr.NewAppErr("segment.not_found", "Segment not found"))
	}

	return response.Success(c, segment)
}
