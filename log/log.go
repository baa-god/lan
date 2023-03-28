package log

import (
	"github.com/baa-god/sharp/sharp"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.Level(-4)
	LevelInfo  = slog.Level(0)
	LevelWarn  = slog.Level(4)
	LevelError = slog.Level(8)
	LevelPanic = slog.Level(12)
	LevelFatal = slog.Level(16)
)

func Trace(msg string, args ...slog.Attr) {
	slog.LogAttrs(nil, LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Panic(msg string, args ...slog.Attr) {
	slog.LogAttrs(nil, LevelPanic, msg, args...)
	panic(msg)
}

func Fatal(msg string, args ...slog.Attr) {
	slog.LogAttrs(nil, LevelFatal, msg, args...)
	os.Exit(1)
}

func SetDefault(w io.Writer) {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(sharp.BaseN(a.Value.String(), 2))
			} else if a.Key == slog.TimeKey {
				value := a.Value.Time().Format("2006-01-02 15:04:05.000")
				a.Value = slog.StringValue(value)
			}
			return a
		},
	}

	h := &Handler{Handler: opts.NewJSONHandler(w)}
	slog.SetDefault(slog.New(h))
}
