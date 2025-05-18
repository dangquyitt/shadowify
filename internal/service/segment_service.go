package service

import (
	"context"
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type SegmentService struct {
	repo *repository.SegmentRepository
}

func NewSegmentService(repo *repository.SegmentRepository) *SegmentService {
	return &SegmentService{repo: repo}
}

// GetSegmentsByVideoID retrieves all segments for a given video ID
func (s *SegmentService) GetSegmentsByVideoID(ctx context.Context, videoID string) ([]*model.Segment, error) {
	return s.repo.FindByVideoID(ctx, videoID)
}
