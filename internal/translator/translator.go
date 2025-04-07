package translator

import (
	"context"
	"shadowify/pkg/config"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Translator interface {
	Translate(ctx context.Context, code string, lang string)
	TranslateWithParams(ctx context.Context, code string, lang string, params map[string]any)
}

type i18nTranslator struct {
	bundle *i18n.Bundle
}

func NewI18nTranslator(cfg *config.I18nConfig) *i18nTranslator {
	return &i18nTranslator{}
}
