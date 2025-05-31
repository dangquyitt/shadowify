package service

import (
	"context"
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type LanguageService struct {
	repo *repository.LanguageRepository
}

func NewLanguageService(repo *repository.LanguageRepository) *LanguageService {
	return &LanguageService{
		repo: repo,
	}
}

// Create creates a new language
func (s *LanguageService) Create(ctx context.Context, req *model.CreateLanguageRequest) (*model.Language, error) {
	language := &model.Language{
		Code:    req.Code,
		Name:    req.Name,
		FlagURL: req.FlagURL,
	}

	err := s.repo.Create(ctx, language)
	if err != nil {
		return nil, err
	}

	return language, nil
}

// GetByID retrieves a language by its ID
func (s *LanguageService) GetByID(ctx context.Context, id string) (*model.Language, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAll retrieves all languages
func (s *LanguageService) GetAll(ctx context.Context) ([]*model.Language, error) {
	return s.repo.GetAll(ctx)
}

// Update updates a language
func (s *LanguageService) Update(ctx context.Context, id string, req *model.UpdateLanguageRequest) (*model.Language, error) {
	language, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only update fields that are provided
	if req.Code != "" {
		language.Code = req.Code
	}
	if req.Name != "" {
		language.Name = req.Name
	}
	if req.FlagURL != "" {
		language.FlagURL = req.FlagURL
	}

	err = s.repo.Update(ctx, language)
	if err != nil {
		return nil, err
	}

	return language, nil
}

// Delete deletes a language by its ID
func (s *LanguageService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
