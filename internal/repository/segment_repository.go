package repository

import (
	"context"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type SegmentRepository struct {
	db *gorm.DB
}

func NewSegmentRepository(db *gorm.DB) *SegmentRepository {
	return &SegmentRepository{db: db}
}

func (r *SegmentRepository) Create(ctx context.Context, segments []*model.Segment) error {
	return r.db.WithContext(ctx).Create(segments).Error
}

func (r *SegmentRepository) FindByVideoID(ctx context.Context, videoID string) ([]*model.Segment, error) {
	var segments []*model.Segment
	err := r.db.WithContext(ctx).Where("video_id = ?", videoID).Order("start_sec ASC").Find(&segments).Error
	if err != nil {
		return nil, err
	}
	return segments, nil
}
