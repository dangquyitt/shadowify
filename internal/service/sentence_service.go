package service

import (
	"context"
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type SentenceService struct {
	sentenceRepository *repository.SentenceRepository
	translatorService  *TranslatorService
}

func NewSentenceService(sentenceRepository *repository.SentenceRepository, translatorService *TranslatorService) *SentenceService {
	return &SentenceService{
		sentenceRepository: sentenceRepository,
		translatorService:  translatorService,
	}
}

func (s *SentenceService) Create(sentence *model.Sentence) error {
	meaningVI, err := s.translatorService.Translate(context.Background(), &model.TranslateInput{
		Text: sentence.MeaningEN,
	})
	if err != nil {
		return err
	}
	sentence.MeaningVI = meaningVI.Text

	return s.sentenceRepository.Create(sentence)
}

func (s *SentenceService) GetByUserIdAndSegmentId(ctx context.Context, userId string, segmentId string) (*model.Sentence, error) {
	return s.sentenceRepository.FindByUserIdAndSegmentId(ctx, userId, segmentId)
}

func (s *SentenceService) List(ctx context.Context, filter *model.SentenceFilter) ([]*model.Sentence, int64, error) {
	return s.sentenceRepository.List(ctx, filter)
}

func (s *SentenceService) DeleteByUserIdAndSegmentId(ctx context.Context, userId string, segmentId string) error {
	return s.sentenceRepository.DeleteByUserIdAndSegmentId(ctx, userId, segmentId)
}
