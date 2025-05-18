package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/model"

	"github.com/jmoiron/sqlx"
)

type VideoRepository struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(ctx context.Context, video *model.Video) error {
	query := `
			INSERT INTO videos (id, title, full_title, description, youtube_id, duration, duration_string, thumbnail, tags, categories, created_at, updated_at)
			VALUES (:id, :title, :full_title, :description, :youtube_id, :duration, :duration_string, :thumbnail, :tags, :categories, :created_at, :updated_at)
		`
	_, err := r.db.NamedExecContext(ctx, query, video)
	if err != nil {
		return apperr.NewAppErr("video.create.error", "Failed to create video").WithCause(err)
	}
	return nil
}

func (r *VideoRepository) GetById(ctx context.Context, id string) (*model.Video, error) {
	var video model.Video
	err := r.db.GetContext(ctx, &video, "SELECT * FROM videos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *VideoRepository) List(ctx context.Context) ([]*model.Video, error) {
	var videos []*model.Video
	err := r.db.SelectContext(ctx, &videos, "SELECT * FROM videos ORDER BY created_at DESC")
	return videos, err
}

func (r *VideoRepository) Update(ctx context.Context, video *model.Video) error {
	query := `UPDATE video SET title=:title, description=:description, tags=:tags, updated_at=:updated_at WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, video)
	return err
}

func (r *VideoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM videos WHERE id = $1", id)
	return err
}
