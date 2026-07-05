package core_logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
	file   *os.File
}

func NewLogger(config Config) (Logger, error) {
	zapLvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.00000")
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	encoderConfig := newEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	return &zapLogger{
		// Skip(1): caller указывает на место вызова log.Debug/Warn/..., а не обёртку zapLogger.
		logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)),
		file:   logFile,
	}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(toZapFields(fields)...),
		file:   l.file,
	}
}

func (l *zapLogger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger file:", err)
	}
}

func toZapFields(fields []Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	result := make([]zap.Field, len(fields))
	for i, field := range fields {
		result[i] = toZapField(field)
	}

	return result
}

func toZapField(field Field) zap.Field {
	switch val := field.val.(type) {
	case string:
		return zap.String(field.key, val)
	case int:
		return zap.Int(field.key, val)
	case error:
		return zap.Error(val)
	case time.Time:
		return zap.Time(field.key, val)
	case time.Duration:
		return zap.Duration(field.key, val)
	default:
		return zap.Any(field.key, val)
	}
}
