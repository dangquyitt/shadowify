package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/model"

	"github.com/jmoiron/sqlx"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func NewSegmentRepository(db *sqlx.DB) *SegmentRepository {
	return &SegmentRepository{db: db}
}

func (r *SegmentRepository) Create(ctx context.Context, segments []*model.Segment) error {
	query := `
		INSERT INTO segments (video_id, start_sec, end_sec, content)
		VALUES (:video_id, :start_sec, :end_sec, :content)
	`
	_, err := r.db.NamedExecContext(ctx, query, segments)
	if err != nil {
		return apperr.NewAppErr("segment.create.error", "Failed to create segments").WithCause(err)
	}

	return nil
}

func (r *SegmentRepository) FindByVideoID(ctx context.Context, videoID string) ([]*model.Segment, error) {
	var segments []*model.Segment
	query := `
		SELECT id, video_id, start_sec, end_sec, content
		FROM segments
		WHERE video_id = $1
		ORDER BY start_sec ASC
	`

	err := r.db.SelectContext(ctx, &segments, query, videoID)
	if err != nil {
		return nil, apperr.NewAppErr("segment.find_by_video_id.error", "Failed to find segments by video ID").WithCause(err)
	}

	return segments, nil
}
