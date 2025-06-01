package repository

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/model"

	"gorm.io/gorm"
)

type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) *LanguageRepository {
	return &LanguageRepository{
		db: db,
	}
}

func (r *LanguageRepository) Create(ctx context.Context, language *model.Language) error {
	err := r.db.WithContext(ctx).Create(language).Error
	if err != nil {
		return apperr.NewAppErr("language.create.error", "Failed to create language").WithCause(err)
	}
	return nil
}

func (r *LanguageRepository) GetByID(ctx context.Context, id string) (*model.Language, error) {
	var language model.Language
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&language).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperr.NewAppErr("language.not_found", "Language not found")
		}
		return nil, apperr.NewAppErr("language.get_by_id.error", "Failed to get language by ID").WithCause(err)
	}
	return &language, nil
}

func (r *LanguageRepository) GetAll(ctx context.Context) ([]*model.Language, error) {
	var languages []*model.Language
	err := r.db.WithContext(ctx).Find(&languages).Error
	if err != nil {
		return nil, apperr.NewAppErr("language.get_all.error", "Failed to get all languages").WithCause(err)
	}
	return languages, nil
}

func (r *LanguageRepository) Update(ctx context.Context, language *model.Language) error {
	err := r.db.WithContext(ctx).Model(&model.Language{}).Where("id = ?", language.Id).Updates(language).Error
	if err != nil {
		return apperr.NewAppErr("language.update.error", "Failed to update language").WithCause(err)
	}
	return nil
}

func (r *LanguageRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Language{}).Error
	if err != nil {
		return apperr.NewAppErr("language.delete.error", "Failed to delete language").WithCause(err)
	}
	return nil
}
