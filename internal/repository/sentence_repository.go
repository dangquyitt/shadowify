package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/ftsearch"
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

func (r *SentenceRepository) List(ctx context.Context, filter *model.SentenceFilter) ([]*model.Sentence, int64, error) {
	var sentences []*model.Sentence
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Sentence{}).Where("user_id = ?", filter.UserId)

	if filter.Q != nil && *filter.Q != "" {
		tsquery := gorm.Expr("to_tsquery('simple', ?)", ftsearch.BuildTsqueryExpression(*filter.Q, ftsearch.WithPrefixMatching()))
		ftSearch := "to_tsvector('simple', coalesce(meaning_en,'') || ' ' || coalesce(meaning_vi,'')) @@ ?"
		query = query.Where(ftSearch, tsquery)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperr.NewAppErr("sentence.list.error", "Failed to count sentences").WithCause(err)
	}

	err := query.Order("created_at DESC").Offset(filter.Offset()).Limit(filter.Limit()).Find(&sentences).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("sentence.list.error", "Failed to list sentences").WithCause(err)
	}

	return sentences, total, nil
}

func (r *SentenceRepository) FindByUserIdAndSegmentId(ctx context.Context, userId string, segmentId string) (*model.Sentence, error) {
	var sentence model.Sentence
	err := r.db.WithContext(ctx).Where("user_id = ? AND segment_id = ?", userId, segmentId).First(&sentence).Error
	if err != nil {
		return nil, apperr.NewAppErr("sentence.find.error", "Failed to find sentence").WithCause(err)
	}
	return &sentence, nil
}

func (r SentenceRepository) DeleteByUserIdAndSegmentId(ctx context.Context, userId string, segmentId string) error {
	if err := r.db.WithContext(ctx).Where("user_id = ? AND segment_id = ?", userId, segmentId).Delete(&model.Sentence{}).Error; err != nil {
		return apperr.NewAppErr("sentence.delete.error", "Failed to delete sentence").WithCause(err)
	}
	return nil
}
