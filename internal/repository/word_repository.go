package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/ftsearch"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) Create(ctx context.Context, word *model.Word) error {
	if err := r.db.WithContext(ctx).Model(&model.Word{}).Create(word).Error; err != nil {
		return apperr.NewAppErr("word.create.error", "Failed to create word").WithCause(err)
	}
	return nil
}

func (r *WordRepository) List(ctx context.Context, filter *model.WordFilter) ([]*model.Word, int64, error) {
	var words []*model.Word
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Word{}).Where("user_id = ?", filter.UserId)

	if filter.Q != nil && *filter.Q != "" {
		tsquery := gorm.Expr("to_tsquery('simple', ?)", ftsearch.BuildTsqueryExpression(*filter.Q, ftsearch.WithPrefixMatching()))
		// Search across title, full_title, and description
		ftSearch := "to_tsvector('simple', coalesce(meaning_en,'') || ' ' || coalesce(meaning_vi,'')) @@ ?"
		query = query.Where(ftSearch, tsquery)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperr.NewAppErr("word.list.error", "Failed to count words").WithCause(err)
	}

	err := query.Order("created_at DESC").Offset(filter.Offset()).Limit(filter.Limit()).Find(&words).Error
	if err != nil {
		return nil, 0, apperr.NewAppErr("word.list.error", "Failed to list words").WithCause(err)
	}
	return words, total, nil
}

func (r *WordRepository) FindByWord(ctx context.Context, word string, userId string) (*model.Word, error) {
	var foundWord model.Word
	if err := r.db.WithContext(ctx).Model(&model.Word{}).
		Where("meaning_en = ? AND user_id = ?", word, userId).
		First(&foundWord).Error; err != nil {
		return nil, apperr.NewAppErr("word.find.error", "Failed to find word").WithCause(err)
	}
	return &foundWord, nil
}

func (r *WordRepository) DeleteByWord(ctx context.Context, word string, userId string) error {
	if err := r.db.WithContext(ctx).Model(&model.Word{}).
		Where("meaning_en = ? AND user_id = ?", word, userId).
		Delete(&model.Word{}).Error; err != nil {
		return apperr.NewAppErr("word.delete.error", "Failed to delete word").WithCause(err)
	}
	return nil
}
