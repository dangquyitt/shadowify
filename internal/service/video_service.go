package service

import (
	"context"
	"net/url"
	"os"
	"shadowify/internal/apperr"
	"shadowify/internal/database"
	"shadowify/internal/dto"
	"shadowify/internal/model"
	"shadowify/internal/repository"
	"strings"
)

type VideoService struct {
	repo           *repository.VideoRepository
	segmentRepo    *repository.SegmentRepository
	whisperService *WhisperService
	ytDLPService   *YTDLPService
}

func NewVideoService(repo *repository.VideoRepository, segmentRepo *repository.SegmentRepository, whisperService *WhisperService, ytDLPService *YTDLPService) *VideoService {
	return &VideoService{
		repo:           repo,
		segmentRepo:    segmentRepo,
		whisperService: whisperService,
		ytDLPService:   ytDLPService,
	}
}

func (s *VideoService) getYoutubeIdFromRawInput(rawInput string) (string, error) {
	if len(rawInput) == 0 {
		return "", apperr.NewAppErr("bad_request", "youtube id is required")
	}

	if strings.Contains(rawInput, "youtube.com/watch?v=") {
		url, err := url.Parse(rawInput)
		if err != nil {
			return "", apperr.NewAppErr("bad_request", "invalid youtube url").WithCause(err)
		}
		queryParams := url.Query()
		return queryParams.Get("v"), nil
	}

	if strings.HasPrefix(rawInput, "https://youtu.be/") {
		url, err := url.Parse(rawInput)
		if err != nil {
			return "", apperr.NewAppErr("bad_request", "invalid youtube url").WithCause(err)
		}
		return strings.TrimPrefix(url.Path, "/"), nil
	}

	return strings.TrimSpace(rawInput), nil
}

func (s *VideoService) Create(ctx context.Context, req *dto.CreateVideoRequest) (*model.Video, error) {
	youtubeId, err := s.getYoutubeIdFromRawInput(req.YoutubeRawInput)
	if err != nil {
		return nil, err
	}

	metadata, filePath, err := s.ytDLPService.DownloadAndExtract(ctx, youtubeId)
	if err != nil {
		return nil, err
	}

	video := &model.Video{
		Title:          metadata.Title,
		FullTitle:      metadata.FullTitle,
		Description:    metadata.Description,
		YoutubeId:      metadata.Id,
		Duration:       metadata.Duration,
		DurationString: metadata.DurationString,
		Thumbnail:      metadata.Thumbnail,
		Tags:           database.JSONType[[]string]{Data: metadata.Tags},
		Categories:     database.JSONType[[]string]{Data: metadata.Categories},
	}

	err = s.repo.Create(ctx, video)
	if err != nil {
		return nil, err
	}

	segments, err := s.whisperService.Transcribe(ctx, filePath)
	if err != nil {
		return nil, err
	}

	for i := range segments {
		segments[i].VideoId = video.Id
	}

	err = s.segmentRepo.Create(ctx, segments)
	if err != nil {
		return nil, err
	}

	err = os.Remove(filePath)
	if err != nil {
		return nil, apperr.NewAppErr("file.remove.error", "Failed to remove downloaded file").WithCause(err)
	}

	return video, nil
}

func (s *VideoService) GetById(ctx context.Context, id string) (*model.Video, error) {
	video, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (s *VideoService) List(ctx context.Context, filter *model.VideoFilter) ([]*model.Video, int64, error) {
	return s.repo.List(ctx, filter)
}
