package service

import (
	"shadowify/internal/video/domain"

	"google.golang.org/api/youtube/v3"
)

type videoService struct {
	ytService *youtube.Service
}

func NewVideoService(ytService *youtube.Service) *videoService {
	return &videoService{ytService: ytService}
}

func (s *videoService) GetVideos() ([]domain.Video, error) {
	var videos []domain.Video
	ytVideos, err := s.ytService.Videos.List([]string{"id"}).Do()
	if err != nil {
		return nil, err
	}
	for _, item := range ytVideos.Items {
		videos = append(videos, domain.Video{
			ID:          item.Id,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			URL:         "https://www.youtube.com/watch?v=" + item.Id,
		})
	}
	return videos, nil
}
