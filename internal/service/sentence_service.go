package service

import (
	"shadowify/internal/model"
	"shadowify/internal/repository"
)

type SentenceService struct {
	sentenceRepository *repository.SentenceRepository
}

func NewSentenceService(sentenceRepository *repository.SentenceRepository) *SentenceService {
	return &SentenceService{sentenceRepository: sentenceRepository}
}

func (s *SentenceService) Create(sentence *model.Sentence) error {
	return s.sentenceRepository.Create(sentence)
}
