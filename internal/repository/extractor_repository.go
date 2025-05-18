package repository

import (
	"context"
	"shadowify/internal/apperr"
	extractor "shadowify/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExtractorRepository struct {
	client extractor.ExtractorServiceClient
}

func NewExtractorRepository(grpcHost string) (*ExtractorRepository, error) {
	conn, err := grpc.NewClient(
		grpcHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &ExtractorRepository{
		client: extractor.NewExtractorServiceClient(conn),
	}, nil
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
