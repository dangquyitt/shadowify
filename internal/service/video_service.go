package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"shadowify/internal/apperr"
	"shadowify/internal/database"
	"shadowify/internal/dto"
	"shadowify/internal/logger"
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
	logger.Infof("Starting video creation with raw input: %s", req.YoutubeRawInput)
	youtubeId, err := s.getYoutubeIdFromRawInput(req.YoutubeRawInput)
	if err != nil {
		return nil, err
	}

	yt, err := s.repo.GetByYoutubeId(ctx, youtubeId)
	if err != nil {
		return nil, apperr.NewAppErr("video.create.error", "Failed to check existing video").WithCause(err)
	}
	if yt != nil {
		return nil, apperr.NewAppErr("video.create.error", "Video already exists").WithCause(err)
	}

	logger.Infof("Starting download and extraction for YouTube ID: %s", youtubeId)
	metadata, filePath, err := s.ytDLPService.DownloadAndExtract(ctx, youtubeId)
	defer func() {
		if filePath == "" {
			return
		}
		err = os.Remove(filePath)
		logger.Errorf("Failed to remove file %s: %v", filePath, err)
	}()
	if err != nil {
		return nil, apperr.NewAppErr("video.create.error", "Failed to download and extract video").WithCause(err)
	}

	lang, err := s.whisperService.DetectLanguage(ctx, filePath)
	if lang != "en" {
		return nil, apperr.NewAppErr("video.create.error", "Only English videos are supported").WithCause(err)
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

	logger.Infof("Starting transcription for video: %s", video.Title)
	segments, err := s.whisperService.Transcribe(ctx, filePath)
	if err != nil {
		return nil, err
	}

	logger.Infof("Starting CEFR prediction for %d segments", len(segments))

	var requestBody struct {
		Sentences []string `json:"sentences"`
	}
	for _, segment := range segments {
		requestBody.Sentences = append(requestBody.Sentences, segment.Content)
	}
	requestBodyByte, err := json.Marshal(requestBody)
	if err != nil {
		return nil, apperr.NewAppErr("request.encode.error", "Failed to encode request body").WithCause(err)
	}
	request, err := http.NewRequest("POST", "http://localhost:5050/predict", bytes.NewBuffer(requestBodyByte))
	if err != nil {
		return nil, apperr.NewAppErr("request.create.error", "Failed to create request").WithCause(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, apperr.NewAppErr("request.execute.error", "Failed to execute request").WithCause(err)
	}
	defer response.Body.Close()

	var responseBody []struct {
		Cefr     string `json:"cefr"`
		Sentence string `json:"sentence"`
	}
	if response.StatusCode != http.StatusOK {
		return nil, apperr.NewAppErr("request.execute.error", "Failed to get response from server").WithCause(err)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, apperr.NewAppErr("response.read.error", "Failed to read response").WithCause(err)
	}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return nil, apperr.NewAppErr("response.decode.error", "Failed to decode response").WithCause(err)
	}

	for i := range segments {
		segments[i].Cefr = responseBody[i].Cefr
	}

	logger.Infof("CEFR prediction completed for %d segments", len(segments))
	err = s.repo.Create(ctx, video, segments)
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) GetById(ctx context.Context, id, userId string) (*model.VideoDetail, error) {
	video, err := s.repo.GetById(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	go func() {
		if err := s.repo.IncrementViewCount(context.TODO(), id); err != nil {
			logger.Warnf("Failed to increment view count for video %s: %v", id, err)
		}
	}()
	return video, nil
}

func (s *VideoService) List(ctx context.Context, filter *model.VideoFilter) ([]*model.Video, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *VideoService) GetFavoriteVideos(ctx context.Context, userId string, filter *model.FavoriteVideoFilter) ([]*model.Video, int64, error) {
	return s.repo.FindFavoriteVideos(ctx, userId, filter)
}

func (s *VideoService) Categories(ctx context.Context) ([]string, error) {
	return s.repo.DistinctCategories(ctx)
}
