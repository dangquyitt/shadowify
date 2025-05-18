package repository

import (
	"context"
	"shadowify/internal/apperr"
	extractor "shadowify/proto"
)

type ExtractorRepository struct {
	client extractor.ExtractorServiceClient
}

func NewExtractorRepository(client extractor.ExtractorServiceClient) *ExtractorRepository {
	return &ExtractorRepository{
		client: client,
	}
}

func (r *ExtractorRepository) ExtractMetadata(ctx context.Context, youtubeId string) (*extractor.YoutubeMetadataResponse, error) {
	req := &extractor.YoutubeRequest{
		YoutubeId: youtubeId,
	}
	resp, err := r.client.ExtractYoutubeMetadata(ctx, req)
	if err != nil {
		return nil, apperr.NewAppErr("internal_error", "failed to extract youtube metadata").WithCause(err)
	}
	return resp, nil
}

func (r *ExtractorRepository) ExtractTranscript(ctx context.Context, youtubeId string) (*extractor.YoutubeTranscriptResponse, error) {
	req := &extractor.YoutubeRequest{
		YoutubeId: youtubeId,
	}
	resp, err := r.client.ExtractYoutubeTranscript(ctx, req)
	if err != nil {
		return nil, apperr.NewAppErr("internal_error", "failed to extract youtube transcript").WithCause(err)
	}
	return resp, nil
}
