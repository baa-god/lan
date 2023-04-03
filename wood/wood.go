package wood

import (
	"fmt"
	"github.com/baa-god/lan/strs"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

var (
	defaultLogger = New(os.Stdout)
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Trace(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelTrace), fmt.Sprint(msg), args...)
}

func (l *Logger) Debug(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelDebug), fmt.Sprint(msg), args...)
}

func (l *Logger) Info(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelInfo), fmt.Sprint(msg), args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelWarn), fmt.Sprint(msg), args...)
}

func (l *Logger) Error(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelError), fmt.Sprint(msg), args...)
}

func (l *Logger) Panic(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelPanic), fmt.Sprint(msg), args...)
}

func (l *Logger) Fatal(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelFatal), fmt.Sprint(msg), args...)
}

func Trace(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelTrace), fmt.Sprint(msg), args...)
}

func Debug(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelDebug), fmt.Sprint(msg), args...)
}

func Info(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelInfo), fmt.Sprint(msg), args...)
}

func Warn(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelWarn), fmt.Sprint(msg), args...)
}

func Error(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelError), fmt.Sprint(msg), args...)
}

func Panic(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelPanic), fmt.Sprint(msg), args...)
}

func Fatal(msg any, args ...any) {
	defaultLogger.Log(nil, slog.Level(LevelFatal), fmt.Sprint(msg), args...)
}

func New(w io.Writer) *Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(strs.BaseN(a.Value.String(), 2))
			} else if a.Key == slog.TimeKey {
				value := a.Value.Time().Format("2006-01-02 15:04:05.000")
				a.Value = slog.StringValue(value)
			} else if a.Key == slog.LevelKey {
				level := Level(a.Value.Any().(slog.Level))
				a.Value = slog.StringValue(level.String())
			}
			return a
		},
	}

	isStd := w == os.Stdout || w == os.Stderr
	handler := &Handler{Handler: opts.NewJSONHandler(w), isStd: isStd}

	return &Logger{Logger: slog.New(handler)}
}

func SetDefault(logger *Logger) {
	defaultLogger = logger
}

func Default() *Logger {
	return defaultLogger
}
