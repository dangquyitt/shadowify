package repository

import (
	"context"
	"encoding/json"
	"shadowify/internal/apperr"
	"shadowify/internal/ftsearch"
	"shadowify/internal/logger"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *model.Video, segments []*model.Segment) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Video{}).Create(video).Error; err != nil {
			return apperr.NewAppErr("video.create.error", "Failed to create video").WithCause(err)
		}

		for _, segment := range segments {
			segment.VideoId = video.Id
		}
		if err := tx.Model(&model.Segment{}).Create(segments).Error; err != nil {
			return apperr.NewAppErr("video.create.error", "Failed to create video segments").WithCause(err)
		}
		return nil
	})
	if err != nil {
		return apperr.NewAppErr("video.create.error", "Failed to create video with segments").WithCause(err)
	}

	return nil
}

func (r *VideoRepository) IncrementViewCount(ctx context.Context, videoId string) error {
	return r.db.WithContext(ctx).Model(&model.Video{}).Where("id = ?", videoId).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
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
	query := r.db.WithContext(ctx).Model(&model.Video{})

	// Apply full-text search filter if provided
	if filter.Q != nil && *filter.Q != "" {
		tsquery := gorm.Expr("to_tsquery('simple', ?)", ftsearch.BuildTsqueryExpression(*filter.Q, ftsearch.WithPrefixMatching()))
		// Search across title, full_title, and description
		ftSearch := "to_tsvector('simple', coalesce(title,'') || ' ' || coalesce(full_title,'') || ' ' || coalesce(description,'')) @@ ?"
		query = query.Where(ftSearch, tsquery)
	}
	if filter.Category != nil && *filter.Category != "" {
		jsonVal, err := json.Marshal([]string{*filter.Category})
		if err != nil {
			return nil, 0, apperr.NewAppErr("filter.encode.error", "Failed to encode category").WithCause(err)
		}

		categoryFilter := gorm.Expr("categories @> ?", string(jsonVal))
		query = query.Where(categoryFilter)
	}

	// Count total with filter
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.list.error", "Failed to count videos").WithCause(err)
	}

	if filter.Type != "" {
		switch filter.Type {
		case model.VideoPopular:
			query = query.Order("view_count DESC")
		default:
			logger.Warnf("Unknown video type filter: %s", filter.Type)
		}
	}
	err = query.Order("created_at DESC").
		Offset(filter.Pagination.Offset()).
		Limit(filter.Pagination.Limit()).
		Find(&videos).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("video.list.error", "Failed to list videos").WithCause(err)
	}

	return videos, total, nil
}

func (r *VideoRepository) DistinctCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).Raw(`SELECT DISTINCT jsonb_array_elements_text(categories) FROM videos WHERE jsonb_typeof(categories) = 'array'`).Pluck("jsonb_array_elements_text", &categories).Error
	if err != nil {
		return nil, apperr.NewAppErr("video.categories.error", "Failed to get distinct video categories").WithCause(err)
	}
	return categories, nil
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
