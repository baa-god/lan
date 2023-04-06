package wood

import (
	"context"
	"fmt"
	"github.com/baa-god/lan/strs"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"os"
	"runtime"
)

type Handler struct {
	slog.Handler
	Attrs []slog.Attr
	Group slog.Attr
	isStd bool
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	pc, file, line, _ := runtime.Caller(4)
	if r.PC = pc; !h.isStd {
		if err = h.Handler.Handle(ctx, r); err != nil {
			return err
		}
	}

	level := Level(r.Level)
	prefix := r.Time.Format("2006-01-02 15:04:05.000")
	prefix = color.HEX("#A9B7C6").Sprint(prefix)
	prefix += " | " + level.ColorString()

	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) { attrs = append(attrs, a) })

	if h.Group.Key != "" {
		h.addGroupAttr(attrs...)
		attrs = []slog.Attr{h.Group}
	}

	s := attrString(append(h.Attrs, attrs...)...)
	s = color.Cyan.Sprint(s)

	source := fmt.Sprintf("%s:%d", file, line)
	source = strs.BaseN(source, 2)
	fmt.Printf("%s | %s > %s %s\n", prefix, source, r.Message, s)

	if level == LevelPanic {
		panic(r.Message)
	} else if level == LevelFatal {
		os.Exit(1)
	}

	return
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := h.Handler.WithAttrs(attrs)
	handler := &Handler{Handler: newHandler, Attrs: h.Attrs, Group: h.Group}

	if handler.Group.Key != "" {
		handler.addGroupAttr(attrs...)
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
		handler.addGroupAttr(group)
	}

	return handler
}

func (h *Handler) addGroupAttr(attrs ...slog.Attr) {
	if v := &lastGroup(&h.Group).Value; v.Kind() == slog.KindGroup {
		*v = slog.GroupValue(append(v.Group(), attrs...)...)
	}
}
