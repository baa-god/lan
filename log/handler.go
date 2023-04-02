package log

import (
	"context"
	"fmt"
	"github.com/baa-god/sharp/sharp"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"runtime"
)

type Handler struct {
	slog.Handler
	Attrs []slog.Attr
	Group slog.Attr
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	pc, file, line, _ := runtime.Caller(4)
	r.PC = pc + 4

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

	if h.Group.Key != "" {
		h.AddGroupAttr(attrs...)
		attrs = []slog.Attr{h.Group}
	}

	s := AttrString(append(h.Attrs, attrs...)...)
	s = color.Cyan.Sprint(s)

	source := fmt.Sprintf("%s:%d", file, line)
	source = sharp.BaseN(source, 2)
	fmt.Printf("%s | %s | %s > %s %s\n", prefix, level, source, r.Message, s)

	return
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := h.Handler.WithAttrs(attrs)
	handler := &Handler{Handler: newHandler, Attrs: h.Attrs, Group: h.Group}

	if handler.Group.Key != "" {
		handler.AddGroupAttr(attrs...)
	} else {
		handler.Attrs = append(h.Attrs, attrs...)
	}

	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	newHandler := h.Handler.WithGroup(name)
	handler := &Handler{Handler: newHandler, Attrs: h.Attrs, Group: h.Group}

	if group := slog.Group(name); h.Group.Key == "" {
		handler.Group = group
	} else {
		handler.AddGroupAttr(group)
	}

	return handler
}

func (h *Handler) AddGroupAttr(attrs ...slog.Attr) {
	if v := &LastGroup(&h.Group).Value; v.Kind() == slog.KindGroup {
		*v = slog.GroupValue(append(v.Group(), attrs...)...)
	}
}
