package rec

import (
	"golang.org/x/exp/slog"
	"io"
	"strings"
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
}

func Fatal(msg string, args ...slog.Attr) {
	slog.LogAttrs(nil, LevelFatal, msg, args...)
}

func SetDefault(w io.Writer) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.String()
				index := strings.LastIndexByte(source, '/')
				index = strings.LastIndexByte(source[:index], '/')
				a.Value = slog.StringValue(source[index+1:])
			} else if a.Key == slog.TimeKey {
				value := a.Value.Time().Format("2006-01-02 15:04:05.000")
				a.Value = slog.StringValue(value)
			}
			return a
		},
	}

	h := &Handler{Handler: opts.NewJSONHandler(w)}
	slog.SetDefault(slog.New(h))

	return slog.Default()
}
