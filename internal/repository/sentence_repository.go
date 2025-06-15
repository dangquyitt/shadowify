package repository

import (
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type SentenceRepository struct {
	db *gorm.DB
}

func NewSentenceRepository(db *gorm.DB) *SentenceRepository {
	return &SentenceRepository{db: db}
}

func (r *SentenceRepository) Create(sentence *model.Sentence) error {
	return r.db.Model(&model.Sentence{}).Create(sentence).Error
}
