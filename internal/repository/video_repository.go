package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/ftsearch"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *model.Video) error {
	err := r.db.WithContext(ctx).Create(video).Error
	if err != nil {
		return apperr.NewAppErr("video.create.error", "Failed to create video").WithCause(err)
	}
	return nil
}

func (r *VideoRepository) GetById(ctx context.Context, id, userId string) (*model.VideoDetail, error) {
	var video model.VideoDetail
	if err := r.db.WithContext(ctx).Model(&model.Video{}).Select("*, (SELECT 1 FROM favorites WHERE user_id = ? AND video_id = videos.id) AS is_favorite", userId).Where("id = ?", id).Scan(&video).Error; err != nil {
		return nil, apperr.NewAppErr("video.get.error", "Failed to get video by ID").WithCause(err)
	}
	return &video, nil
}

func (r *VideoRepository) List(ctx context.Context, filter *model.VideoFilter) ([]*model.Video, int64, error) {
	var videos []*model.Video
	var total int64

	// Base queries for count and data selection
	countQuery := r.db.WithContext(ctx).Model(&model.Video{})
	query := r.db.WithContext(ctx).Model(&model.Video{})

	// Apply full-text search filter if provided
	if filter.Q != nil && *filter.Q != "" {
		tsquery := gorm.Expr("to_tsquery('simple', ?)", ftsearch.BuildTsqueryExpression(*filter.Q, ftsearch.WithPrefixMatching()))
		// Search across title, full_title, and description
		ftSearch := "to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(full_title,'') || ' ' || coalesce(description,'')) @@ ?"
		countQuery = countQuery.Where(ftSearch, tsquery)
		query = query.Where(ftSearch, tsquery)
	}

	// Count total with filter
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.list.error", "Failed to count videos").WithCause(err)
	}

	// Fetch paginated data with filter
	err = query.Order("created_at DESC").
		Offset(filter.Pagination.Offset()).
		Limit(filter.Pagination.Limit()).
		Find(&videos).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.list.error", "Failed to list videos").WithCause(err)
	}

	return videos, total, nil
}

func (r *VideoRepository) Update(ctx context.Context, video *model.Video) error {
	return r.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", video.Id).Updates(video).Error
}

func (r *VideoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Video{}).Error
}

func (r *VideoRepository) FindFavoriteVideos(ctx context.Context, userId string, filter *model.FavoriteVideoFilter) ([]*model.Video, int64, error) {
	var videos []*model.Video
	query := r.db.WithContext(ctx).
		Model(&model.Video{}).
		Joins("JOIN favorites ON favorites.video_id = videos.id").
		Where("favorites.user_id = ?", userId).
		Order("favorites.created_at DESC")

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.find_favorite.error", "Failed to count favorite videos").WithCause(err)
	}

	err = query.Offset(filter.Pagination.Offset()).
		Limit(filter.Pagination.Limit()).
		Find(&videos).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.find_favorite.error", "Failed to list favorite videos").WithCause(err)
	}
	return videos, total, nil
}
