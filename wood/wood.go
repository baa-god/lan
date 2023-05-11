package wood

import (
	"fmt"
	"github.com/baa-god/lan/strs"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"sync/atomic"
)

var (
	defaultLogger atomic.Value
)

func init() {
	defaultLogger.Store(New(os.Stdout))
}

type Logger struct {
	*slog.Logger
}

func (l *Logger) Trace(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelTrace), fmt.Sprint(msg), args...)
}

func (l *Logger) Tracef(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelTrace), fmt.Sprintf(msg, a...))
}

func (l *Logger) Debug(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelDebug), fmt.Sprint(msg), args...)
}

func (l *Logger) Debugf(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelDebug), fmt.Sprintf(msg, a...))
}

func (l *Logger) Infof(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelInfo), fmt.Sprintf(msg, a...))
}

func (l *Logger) Info(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelInfo), fmt.Sprint(msg), args...)
}

func (l *Logger) Warn(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelWarn), fmt.Sprint(msg), args...)
}

func (l *Logger) Warnf(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelWarn), fmt.Sprintf(msg, a...))
}

func (l *Logger) Error(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelError), fmt.Sprint(msg), args...)
}

func (l *Logger) Errorf(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelError), fmt.Sprintf(msg, a...))
}

func (l *Logger) Panic(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelPanic), fmt.Sprint(msg), args...)
}

func (l *Logger) Panicf(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelPanic), fmt.Sprintf(msg, a...))
}

func (l *Logger) Fatal(msg any, args ...any) {
	l.Log(nil, slog.Level(LevelFatal), fmt.Sprint(msg), args...)
}

func (l *Logger) Fatalf(msg string, a ...any) {
	l.Log(nil, slog.Level(LevelFatal), fmt.Sprintf(msg, a...))
}

func Trace(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelTrace), fmt.Sprint(msg), args...)
}

func Tracef(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelTrace), fmt.Sprintf(msg, a...))
}

func Debug(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelDebug), fmt.Sprint(msg), args...)
}

func Debugf(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelDebug), fmt.Sprintf(msg, a...))
}

func Info(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelInfo), fmt.Sprint(msg), args...)
}

func Infof(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelInfo), fmt.Sprintf(msg, a...))
}

func Warn(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelWarn), fmt.Sprint(msg), args...)
}

func Warnf(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelWarn), fmt.Sprintf(msg, a...))
}

func Error(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelError), fmt.Sprint(msg), args...)
}

func Errorf(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelError), fmt.Sprintf(msg, a...))
}

func Panic(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelPanic), fmt.Sprint(msg), args...)
}

func Panicf(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelPanic), fmt.Sprintf(msg, a...))
}

func Fatal(msg any, args ...any) {
	Default().Log(nil, slog.Level(LevelFatal), fmt.Sprint(msg), args...)
}

func Fatalf(msg string, a ...any) {
	Default().Log(nil, slog.Level(LevelFatal), fmt.Sprintf(msg, a...))
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
	handler := &Handler{Handler: slog.NewJSONHandler(w, &opts), isStd: isStd}

	return &Logger{Logger: slog.New(handler)}
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func Default() *Logger {
	return defaultLogger.Load().(*Logger)
}
