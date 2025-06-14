package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"shadowify/internal/model"

	"github.com/google/uuid"
)

type YTDLPService struct {
}

func NewYTDLPService() *YTDLPService {
	return &YTDLPService{}
}

func (s *YTDLPService) DownloadAndExtract(ctx context.Context, youtubeId string) (*model.YoutubeMetadata, string, error) {
	uid := uuid.New().String()
	outputBase := filepath.Join("./tmp", uid)

	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-x",
		"--audio-format", "wav",
		"--write-info-json",
		"-o", outputBase,
		"https://www.youtube.com/watch?v="+youtubeId,
	)

	if err := cmd.Run(); err != nil {
		return nil, "", fmt.Errorf("yt-dlp command failed: %w", err)
	}

	jsonPath := outputBase + ".info.json"
	defer os.Remove(jsonPath)

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read metadata json file: %w", err)
	}

	metadata := &model.YoutubeMetadata{}
	if err := json.Unmarshal(jsonData, metadata); err != nil {
		return nil, "", fmt.Errorf("failed to parse yt-dlp metadata JSON: %w", err)
	}

	audioPath := outputBase + ".wav"

	return metadata, audioPath, nil
}
