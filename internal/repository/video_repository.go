package repository

import (
	"shadowify/internal/model"

	"github.com/jmoiron/sqlx"
)

type VideoRepository struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepository {
	return &VideoRepository{
		db: db,
	}
}

func (r *VideoRepository) FindById(id string) (*model.Video, error) {
	var video model.Video
	err := r.db.Get(&video, "SELECT * FROM videos WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *VideoRepository) DeleteById(id string) error {
	_, err := r.db.Exec("DELETE FROM videos WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
