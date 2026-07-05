package core_logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	levelWidth  = 5
	callerWidth = 28
	consoleSep  = "  "
)

func bracketedTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	formatted := t.Format("2006-01-02T15:04:05.000000")
	enc.AppendString(fmt.Sprintf("[==%s==]", formatted))
}

func paddedColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	label := level.CapitalString()
	if len(label) < levelWidth {
		label += strings.Repeat(" ", levelWidth-len(label))
	}
	enc.AppendString(colorizeLevel(label, level))
}

func colorizeLevel(label string, level zapcore.Level) string {
	const reset = "\x1b[0m"

	switch level {
	case zapcore.DebugLevel:
		return "\x1b[36m" + label + reset
	case zapcore.InfoLevel:
		return "\x1b[34m" + label + reset
	case zapcore.WarnLevel:
		return "\x1b[33m" + label + reset
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return "\x1b[31m" + label + reset
	default:
		return label
	}
}

func paddedShortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	s := caller.TrimmedPath()
	if len(s) > callerWidth {
		s = s[len(s)-callerWidth:]
	}
	enc.AppendString(fmt.Sprintf("%-*s", callerWidth, s))
}

func newEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = bracketedTimeEncoder
	cfg.EncodeLevel = paddedColorLevelEncoder
	cfg.EncodeCaller = paddedShortCallerEncoder
	cfg.ConsoleSeparator = consoleSep
	return cfg
}
