package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"shadowify/internal/logger"
	"shadowify/internal/model"

	"github.com/google/uuid"
)

type STTService struct {
	whisperService *WhisperService
}

func NewSTTService(whisperService *WhisperService) *STTService {
	return &STTService{
		whisperService: whisperService,
	}
}

func (s *STTService) Transcribe(ctx context.Context, input *model.TranscribeInput) (*model.TranscribeOutput, error) {
	audioData, err := base64.StdEncoding.DecodeString(input.AudioBase64)
	if err != nil {
		return nil, err
	}

	outputPath := filepath.Join("./tmp", fmt.Sprintf("%s.m4a", uuid.NewString()))
	err = os.WriteFile(outputPath, audioData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write audio file: %w", err)
	}

	logger.Infof("Transcribing audio file: %s", outputPath)

	text, err := s.whisperService.TranscribeNoTimestamps(ctx, outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	return &model.TranscribeOutput{
		Text: text,
	}, nil
}
