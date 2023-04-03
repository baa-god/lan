package wood

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Trace(msg any, args ...slog.Attr) {
	m := fmt.Sprint(msg)
	l.Logger.LogAttrs(nil, slog.Level(LevelTrace), m, args...)
}

func (l *Logger) Debug(msg any, args ...any) {
	l.Logger.Debug(fmt.Sprint(msg), args...)
}

func (l *Logger) Info(msg any, args ...any) {
	l.Logger.Info(fmt.Sprint(msg), args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	l.Logger.Warn(fmt.Sprint(msg), args...)
}

func (l *Logger) Error(msg any, args ...any) {
	l.Logger.Error(fmt.Sprint(msg), args...)
}

func (l *Logger) Panic(msg any, args ...slog.Attr) {
	m := fmt.Sprint(msg)
	l.Logger.LogAttrs(nil, slog.Level(LevelPanic), m, args...)
	panic(msg)
}

func (l *Logger) Fatal(msg any, args ...slog.Attr) {
	m := fmt.Sprint(msg)
	l.Logger.LogAttrs(nil, slog.Level(LevelFatal), m, args...)
	os.Exit(1)
}
