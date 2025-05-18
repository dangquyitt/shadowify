package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/transcript/model"

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
		VALUES (:video_id, :start_sec, :start_sec, :content)
	`
	_, err := r.db.NamedExecContext(ctx, query, segments)
	if err != nil {
		return apperr.NewAppErr("segment.create.error", "Failed to create segments").WithCause(err)
	}

	return nil
}
