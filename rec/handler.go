package rec

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"runtime"
	"strings"
)

type Handler struct {
	slog.Handler
	attrs []slog.Attr
	group slog.Attr
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	pc, file, line, _ := runtime.Caller(4)
	r.PC += pc + 4

	if err = h.Handler.Handle(ctx, r); err != nil {
		return err
	}

	level := r.Level.String()
	prefix := r.Time.Format("2006-01-02 15:05:05.000")
	prefix = color.HEX("#A9B7C6").Sprint(prefix)

	switch r.Level {
	case LevelTrace:
		level = color.Hex("#7970A9").Sprint("TRACE")
	case slog.LevelDebug:
		level = color.Hex("#808080").Sprint(level)
	case slog.LevelInfo:
		level = color.Green.Sprint(level + " ")
	case slog.LevelWarn:
		level = color.Yellow.Sprint(level + " ")
	case slog.LevelError:
		level = color.Hex("#ff3800").Sprint(level)
	case LevelPanic:
		level = color.Hex("#F998CC").Sprint("PANIC")
	case LevelFatal:
		level = color.Hex("#FE4EDA").Sprint("FATAL")
	}

	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) { attrs = append(attrs, a) })

	if h.group.Key != "" {
		h.AddGroupAttr(attrs...)
		attrs = []slog.Attr{h.group}
	}

	attrStr := AttrString(append(h.attrs, attrs...)...)
	attrStr = color.Cyan.Sprint(attrStr)

	fmt.Printf(
		"%s | %s | %s:%d > %s %s\n",
		prefix, level, file, line, r.Message, attrStr,
	)

	if r.Level == LevelPanic {
		// panic(r.Message)
	} else if r.Level == LevelFatal {
		// os.Exit(1)
	}

	return
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := h.Handler.WithAttrs(attrs)
	handler := &Handler{Handler: newHandler, attrs: h.attrs, group: h.group}

	if handler.group.Key != "" {
		handler.AddGroupAttr(attrs...)
	} else {
		handler.attrs = append(h.attrs, attrs...)
	}

	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newHandler := h.Handler.WithGroup(name)
	handler := &Handler{Handler: newHandler, attrs: h.attrs, group: h.group}

	if group := slog.Group(name); h.group.Key == "" {
		handler.group = group
	} else {
		handler.AddGroupAttr(group)
	}

	return handler
}

func AttrString(attrs ...slog.Attr) (s string) {
	for _, x := range attrs {
		if v, ok := x.Value.Any().([]slog.Attr); ok {
			s += x.Key + "{" + AttrString(v...) + "} "
			continue
		}

		v, kind := x.Value.String(), x.Value.Kind()
		if kind != slog.KindInt64 && kind != slog.KindUint64 &&
			kind != slog.KindBool && kind != slog.KindDuration {
			v = `"` + v + `"`
		}

		s += fmt.Sprintf("%s=%s ", x.Key, v)
	}

	return strings.TrimSuffix(s, " ")
}

func LastGroup(attr *slog.Attr) *slog.Attr {
	if v, _ := attr.Value.Any().([]slog.Attr); v != nil {
		for i := len(v) - 1; i >= 0; i-- {
			if x := &v[i]; x.Value.Kind() == slog.KindGroup {
				return LastGroup(x)
			}
		}
	}
	return attr
}

func (h *Handler) AddGroupAttr(attrs ...slog.Attr) {
	v := &LastGroup(&h.group).Value
	if v.Kind() == slog.KindGroup {
		*v = slog.GroupValue(append(v.Group(), attrs...)...)
	}
}
