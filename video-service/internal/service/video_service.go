package service

import (
	"context"

	"github.com/dangquyitt/shadowify/video-service/internal/model"
	"github.com/dangquyitt/shadowify/video-service/internal/repository"
)

const URLS = "https://www.youtube.com/watch?v="

type VideoService interface {
	FindAll(ctx context.Context) ([]model.Video, error)
	GetTranscript(ctx context.Context, externalId string) ([]model.Transcript, error)
}

type videoService struct {
	repository           repository.VideoRepository
	transcriptionService TranscriptionService
}

func NewVideoService(repository repository.VideoRepository, transcriptionService TranscriptionService) *videoService {
	return &videoService{repository: repository, transcriptionService: transcriptionService}
}

func (s *videoService) FindAll(ctx context.Context) ([]model.Video, error) {
	return s.repository.FindAll(ctx)
}

func (s *videoService) GetTranscript(ctx context.Context, externalId string) ([]model.Transcript, error) {
	return s.transcriptionService.GetTranscript(ctx, URLS+externalId)
}
