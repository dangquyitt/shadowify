package service

import "shadowify/internal/repository"

type VideoService struct {
	repository *repository.VideoRepository
}

func NewVideoService(repository *repository.VideoRepository) *VideoService {
	return &VideoService{
		repository: repository,
	}
}
