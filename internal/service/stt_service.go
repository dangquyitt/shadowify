package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"shadowify/internal/apperr"
	"shadowify/internal/logger"
	"shadowify/internal/model"

	"github.com/google/uuid"
)

type STTService struct {
	whisperService    *WhisperService
	translatorService *TranslatorService
}

func NewSTTService(whisperService *WhisperService, translatorService *TranslatorService) *STTService {
	return &STTService{
		whisperService:    whisperService,
		translatorService: translatorService,
	}
}

func (s *STTService) EvaluateAudio(ctx context.Context, input *model.EvaluateInput) (*model.EvaluateOutput, error) {
	audioData, err := base64.StdEncoding.DecodeString(input.AudioBase64)
	if err != nil {
		return nil, apperr.NewAppErr("stt.decode.error", "Failed to decode audio base64").WithCause(err)
	}

	filePath := filepath.Join("./tmp", fmt.Sprintf("%s.m4a", uuid.NewString()))
	err = os.WriteFile(filePath, audioData, 0644)
	if err != nil {
		return nil, apperr.NewAppErr("stt.write.error", "Failed to write audio file").WithCause(err)
	}
	// defer os.Remove(filePath)
	// logger.Infof("Evaluating audio file: %s", filePath)
	// lang, err := s.whisperService.DetectLanguage(ctx, filePath)
	// if lang != "en" {
	// 	return nil, apperr.NewAppErr("video.create.error", "Only English videos are supported").WithCause(err)
	// }

	meaningEN, err := s.whisperService.TranscribeNoTimestamps(ctx, filePath)
	if err != nil {
		return nil, apperr.NewAppErr("stt.transcribe.error", "Failed to transcribe audio").WithCause(err)
	}

	tranOutput, err := s.translatorService.Translate(ctx, &model.TranslateInput{Text: meaningEN})
	if err != nil {
		return nil, apperr.NewAppErr("translator.translate.error", "Failed to translate meaning").WithCause(err)
	}

	output := &model.EvaluateOutput{
		MeaningEN: meaningEN,
		MeaningVI: tranOutput.Text,
	}

	var requestBody struct {
		Sentences []string `json:"sentences"`
	}

	requestBody.Sentences = []string{meaningEN}

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

	output.Cefr = responseBody[0].Cefr
	return output, nil
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
