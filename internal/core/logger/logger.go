package core_logger

import "context"

// Logger — абстракция над логгером приложения.
// Не зависит от конкретной библиотеки (zap, slog и т.д.).
type Logger interface {
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Close()
}

type contextKey struct{}

var loggerKey contextKey

func ContextWithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		panic("logger not found in context")
	}
	return log
}
