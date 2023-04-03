package wood

import (
	"context"
	"fmt"
	"github.com/baa-god/lan/lan"
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
	r.PC -= 1
	frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()

	if err = h.Handler.Handle(ctx, r); err != nil {
		return err
	}

	prefix := r.Time.Format("2006-01-02 15:05:05.000")
	prefix = color.HEX("#A9B7C6").Sprint(prefix)
	prefix += " | " + Level(r.Level).ColorString()

	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) { attrs = append(attrs, a) })

	if h.Group.Key != "" {
		h.AddGroupAttr(attrs...)
		attrs = []slog.Attr{h.Group}
	}

	s := AttrString(append(h.Attrs, attrs...)...)
	s = color.Cyan.Sprint(s)

	source := fmt.Sprintf("%s:%d", frame.File, frame.Line)
	source = lan.BaseN(source, 2)
	fmt.Printf("%s | %s > %s %s\n", prefix, source, r.Message, s)

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
