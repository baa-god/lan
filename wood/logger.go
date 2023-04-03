package wood

import (
	"golang.org/x/exp/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Trace(msg string, args ...slog.Attr) {
	l.Logger.LogAttrs(nil, slog.Level(LevelTrace), msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *Logger) Panic(msg string, args ...slog.Attr) {
	l.Logger.LogAttrs(nil, slog.Level(LevelPanic), msg, args...)
	panic(msg)
}

func (l *Logger) Fatal(msg string, args ...slog.Attr) {
	l.Logger.LogAttrs(nil, slog.Level(LevelFatal), msg, args...)
	os.Exit(1)
}
