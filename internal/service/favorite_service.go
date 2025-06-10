package service

import (
	"context"
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type FavoriteService struct {
	favoriteRepository *repository.FavoriteRepository
}

func NewFavoriteService(favoriteRepository *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		favoriteRepository: favoriteRepository,
	}
}

func (s *FavoriteService) Create(ctx context.Context, userId string, videoId string) error {
	favorite, _ := s.favoriteRepository.FindByUserIdAndVideoId(ctx, userId, videoId)
	if favorite != nil {
		return nil
	}
	newFavorite := &model.Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	return s.favoriteRepository.Create(ctx, newFavorite)
}

func (s *FavoriteService) Delete(ctx context.Context, userId string, videoId string) error {
	favorite, err := s.favoriteRepository.FindByUserIdAndVideoId(ctx, userId, videoId)
	if err != nil {
		return err
	}
	if favorite == nil {
		return nil // No favorite found to delete
	}
	return s.favoriteRepository.Delete(ctx, userId, videoId)
}
