package service

import "shadowify/internal/transcript/repository"

type SegmentService struct {
	repo *repository.SegmentRepository
}

func NewSegmentService(repo *repository.SegmentRepository) *SegmentService {
	return &SegmentService{repo: repo}
}
