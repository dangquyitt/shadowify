package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"shadowify/internal/apperr"
	"shadowify/internal/config"
	"shadowify/internal/model"
)

type TranslatorService struct {
	cfg    config.AzureTranslatorConfig
	client *http.Client
}

func NewTranslatorService(cfg config.AzureTranslatorConfig) *TranslatorService {
	return &TranslatorService{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func (s *TranslatorService) Translate(ctx context.Context, input *model.TranslateInput) (*model.TranslateOutput, error) {
	u, _ := url.Parse(s.cfg.URI)
	q := u.Query()
	q.Add("from", "en")
	q.Add("to", "vi")
	u.RawQuery = q.Encode()

	body := []struct {
		Text string
	}{
		{Text: input.Text},
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, apperr.NewAppErr("translator.request.error", "Failed to create request").WithCause(err)
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", s.cfg.APIKey)
	req.Header.Add("Ocp-Apim-Subscription-Region", s.cfg.Region)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, apperr.NewAppErr("translator.request.error", "Failed to send request").WithCause(err)
	}

	type TranslationResponse []struct {
		Translations []struct {
			Text string `json:"text"`
			To   string `json:"to"`
		} `json:"translations"`
	}

	var result TranslationResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	translatedText := ""
	if len(result) > 0 && len(result[0].Translations) > 0 {
		translatedText = result[0].Translations[0].Text
	}

	return &model.TranslateOutput{
		Text: translatedText,
	}, nil // Placeholder for actual translation logic
}
