package repository

import (
	"context"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

func (r *FavoriteRepository) FindByUserIdAndVideoId(ctx context.Context, userId string, videoId string) (*model.Favorite, error) {
	var favorite model.Favorite
	err := r.db.WithContext(ctx).Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	err := r.db.WithContext(ctx).Create(favorite).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRepository) Delete(ctx context.Context, userId string, videoId string) error {
	err := r.db.WithContext(ctx).Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&model.Favorite{}).Error
	if err != nil {
		return err
	}
	return nil
}
