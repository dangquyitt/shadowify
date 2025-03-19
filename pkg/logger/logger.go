package logger

import "context"

type Logger interface {
	Debug(ctx context.Context, v any)
	Debugf(ctx context.Context, format string, v ...any)
	Info(ctx context.Context, v any)
	Infof(ctx context.Context, format string, v ...any)
	Warn(ctx context.Context, v any)
	Warnf(ctx context.Context, format string, v ...any)
	Error(ctx context.Context, v any)
	Errorf(ctx context.Context, format string, v ...any)
}
