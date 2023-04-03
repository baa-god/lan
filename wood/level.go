package wood

import (
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
)

const (
	LevelTrace = Level(-8)
	LevelDebug = Level(-4)
	LevelInfo  = Level(0)
	LevelWarn  = Level(4)
	LevelError = Level(8)
	LevelPanic = Level(12)
	LevelFatal = Level(16)
)

type Level slog.Level

func (l Level) Level() slog.Level {
	return slog.Level(l)
}

func (l Level) String() string {
	switch {
	case l < LevelDebug:
		return "TRACE"
	case l < LevelInfo:
		return "DEBUG"
	case l < LevelWarn:
		return "INFO"
	case l < LevelError:
		return "WARN"
	case l < LevelPanic:
		return "ERROR"
	case l < LevelFatal:
		return "PANIC"
	default:
		return "FATAL"
	}
}

func (l Level) ColorString() string {
	switch level := l.String(); l {
	case LevelTrace:
		return color.Hex("#7970A9").Sprint(level)
	case LevelDebug:
		return color.Hex("#808080").Sprint(level)
	case LevelInfo:
		return color.Green.Sprint(level + " ")
	case LevelWarn:
		return color.Yellow.Sprint(level + " ")
	case LevelError:
		return color.Hex("#FF3800").Sprint(level)
	case LevelPanic:
		return color.Hex("#F998CC").Sprint(level)
	default:
		return color.Hex("#FE4EDA").Sprint(level)
	}
}
