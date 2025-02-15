package service

import (
	"context"
	"encoding/json"

	"github.com/dangquyitt/shadowify/video-service/internal/model"
	pb "github.com/dangquyitt/shadowify/video-service/proto"
)

type TranscriptionService interface {
	GetTranscript(ctx context.Context, videoURL string) ([]model.Transcript, error)
}

type transcriptionService struct {
	client pb.TranscriptionServiceClient
}

func NewTranscriptionService(client pb.TranscriptionServiceClient) *transcriptionService {
	return &transcriptionService{client: client}
}

func (s *transcriptionService) GetTranscript(ctx context.Context, videoURL string) ([]model.Transcript, error) {
	req := &pb.GetTranscriptRequest{VideoUrl: videoURL}
	res, err := s.client.GetTranscript(ctx, req)
	if err != nil {
		return nil, err
	}

	var result []model.Transcript
	json.Unmarshal([]byte(res.Transcript), &result)

	return result, nil
}
