package logger

import (
	"context"
	"os"
	"shadowify/pkg/config"

	"github.com/rs/zerolog"
)

var zerologLevelMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
}

type zerologAdapter struct {
	l *zerolog.Logger
}

func NewZerologAdapter(cfg config.LoggerConfig) *zerologAdapter {
	level, ok := zerologLevelMap[cfg.Level]
	if !ok {
		level = zerolog.InfoLevel
	}

	l := zerolog.New(os.Stdout).Level(level)

	return &zerologAdapter{l: &l}
}

func (l *zerologAdapter) Debug(ctx context.Context, v any) {
	l.l.Debug().Ctx(ctx).Msgf("%v", v)
}

func (l *zerologAdapter) Debugf(ctx context.Context, format string, v ...any) {
	l.l.Debug().Ctx(ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Info(ctx context.Context, v any) {
	l.l.Info().Ctx(ctx).Msgf("%v", v)
}

func (l *zerologAdapter) Infof(ctx context.Context, format string, v ...any) {
	l.l.Info().Ctx(ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Warn(ctx context.Context, v any) {
	l.l.Warn().Ctx(ctx).Msgf("%v", v)
}

func (l *zerologAdapter) Warnf(ctx context.Context, format string, v ...any) {
	l.l.Warn().Ctx(ctx).Msgf(format, v...)
}

func (l *zerologAdapter) Error(ctx context.Context, v any) {
	l.l.Error().Ctx(ctx).Msgf("%v", v)
}

func (l *zerologAdapter) Errorf(ctx context.Context, format string, v ...any) {
	l.l.Error().Ctx(ctx).Msgf(format, v...)
}
