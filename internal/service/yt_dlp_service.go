package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"shadowify/internal/model"
)

type YTDLPService struct {
}

func NewYTDLPService() *YTDLPService {
	return &YTDLPService{}
}

func (s *YTDLPService) DownloadAndExtract(ctx context.Context, youtubeId string) (*model.YoutubeMetadata, string, error) {
	url := "https://www.youtube.com/watch?v=" + youtubeId
	outputPath := filepath.Join("../../tmp", youtubeId+".wav")

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "wav",
		"-o", outputPath,
		"--print-json",
		url,
	)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, "", fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, stdoutPipe)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read stdout: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, "", fmt.Errorf("yt-dlp command failed: %w", err)
	}

	metadata := &model.YoutubeMetadata{}
	err = json.Unmarshal(buf.Bytes(), metadata)
	if err != nil {
		scanner := bufio.NewScanner(bytes.NewReader(buf.Bytes()))
		for scanner.Scan() {
			line := scanner.Bytes()
			if err := json.Unmarshal(line, metadata); err == nil {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, "", fmt.Errorf("error scanning stdout: %w", err)
		}
	}

	return metadata, outputPath, nil
}
