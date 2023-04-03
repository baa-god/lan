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
	m := fmt.Sprint(msg)
	l.Logger.Debug(m, args...)
}

func (l *Logger) Info(msg any, args ...any) {
	m := fmt.Sprint(msg)
	l.Logger.Info(m, args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	m := fmt.Sprint(msg)
	l.Logger.Warn(m, args...)
}

func (l *Logger) Error(msg any, args ...any) {
	m := fmt.Sprint(msg)
	l.Logger.Error(m, args...)
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
