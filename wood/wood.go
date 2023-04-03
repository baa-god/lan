package wood

import (
	"github.com/baa-god/lan/lan"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

var (
	defaultLogger *Logger
)

func init() {
	defaultLogger = New(os.Stdout)
}

func Trace(msg string, args ...slog.Attr) {
	defaultLogger.Trace(msg, args...)
}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func Panic(msg string, args ...slog.Attr) {
	defaultLogger.Panic(msg, args...)
}

func Fatal(msg string, args ...slog.Attr) {
	defaultLogger.Fatal(msg, args...)
}

func New(w io.Writer) *Logger {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(lan.BaseN(a.Value.String(), 2))
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

	handler := &Handler{Handler: opts.NewJSONHandler(w)}
	return &Logger{Logger: slog.New(handler)}
}
