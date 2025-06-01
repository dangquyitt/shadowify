package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"shadowify/internal/model"
	"strings"
)

type WhisperService struct {
}

func NewWhisperService() *WhisperService {
	return &WhisperService{}
}

func (s *WhisperService) DetectLanguage(ctx context.Context, audioFilePath string) (string, error) {
	wd, _ := os.Getwd()
	cmd := exec.Command(filepath.Join(wd, "lib/whisper-cli"),
		"-m", filepath.Join(wd, "lib/ggml-large-v3-turbo.bin"),
		"-f", audioFilePath,
		"-np",
		"-dl",
		"-t", fmt.Sprintf("%d", runtime.NumCPU()-2),
		"-oj",
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run whisper-cli: %w\nOutput: %s", err, audioFilePath)
	}

	jsonPath := audioFilePath + ".json"
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to read json output: %w", err)
	}

	var parsed struct {
		Result struct {
			Language string `json:"language"`
		} `json:"result"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}

	if err := os.Remove(jsonPath); err != nil {
		return "", fmt.Errorf("failed to delete json file: %w", err)
	}

	return parsed.Result.Language, nil
}

func (s *WhisperService) Transcribe(ctx context.Context, audioFilePath string) ([]*model.Segment, error) {
	wd, _ := os.Getwd()
	cmd := exec.Command(filepath.Join(wd, "lib/whisper-cli"),
		"-m", filepath.Join(wd, "lib/ggml-large-v3-turbo.bin"),
		"-f", audioFilePath,
		"-np",
		"-t", fmt.Sprintf("%d", runtime.NumCPU()-2),
		"-oj",
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run whisper-cli: %w\nOutput: %s", err, audioFilePath)
	}

	jsonPath := audioFilePath + ".json"
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read json output: %w", err)
	}

	var parsed struct {
		Transcription []struct {
			Text    string `json:"text"`
			Offsets struct {
				From float32 `json:"from"`
				To   float32 `json:"to"`
			}
		} `json:"transcription"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	segments := make([]*model.Segment, len(parsed.Transcription))
	for i, t := range parsed.Transcription {
		segments[i] = &model.Segment{
			StartSec: t.Offsets.From / 1000,
			EndSec:   t.Offsets.To / 1000,
			Content:  strings.TrimSpace(t.Text),
		}
	}

	if err := os.Remove(jsonPath); err != nil {
		return nil, fmt.Errorf("failed to delete json file: %w", err)
	}

	return segments, nil
}
