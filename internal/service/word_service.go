package service

import (
	"context"
	"shadowify/internal/apperr"
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type WordService struct {
	wordRepository    *repository.WordRepository
	translatorService *TranslatorService
}

func NewWordService(wordRepository *repository.WordRepository, translatorService *TranslatorService) *WordService {
	return &WordService{
		wordRepository:    wordRepository,
		translatorService: translatorService,
	}
}

func (s *WordService) GetByWord(ctx context.Context, word string, userId string) (*model.Word, error) {
	if word == "" {
		return nil, apperr.NewAppErr("bad_request", "Word is required")
	}

	if userId == "" {
		return nil, apperr.NewAppErr("bad_request", "User ID is required")
	}

	return s.wordRepository.FindByWord(ctx, word, userId)
}

func (s *WordService) List(ctx context.Context, filter *model.WordFilter) ([]*model.Word, int64, error) {
	return s.wordRepository.List(ctx, filter)
}

func (s *WordService) Create(ctx context.Context, word *model.Word) error {
	if word.UserId == "" {
		return apperr.NewAppErr("bad_request", "User ID is required")
	}

	if word.MeaningEN == "" {
		return apperr.NewAppErr("bad_request", "Meaning in English is required")
	}

	w, _ := s.wordRepository.FindByWord(ctx, word.MeaningEN, word.UserId)
	if w != nil {
		return apperr.NewAppErr("word_exists", "Word already exists").WithCause(nil)
	}

	meaningVI, err := s.translatorService.Translate(ctx, &model.TranslateInput{
		Text: word.MeaningEN,
	})
	if err != nil {
		return apperr.NewAppErr("translation_error", "Failed to translate meaning").WithCause(err)
	}
	word.MeaningVI = meaningVI.Text

	return s.wordRepository.Create(ctx, word)
}

func (s *WordService) DeleteByWord(ctx context.Context, word string, userId string) error {
	if word == "" {
		return apperr.NewAppErr("bad_request", "Word is required")
	}

	return s.wordRepository.DeleteByWord(ctx, word, userId)
}
