package wood

import (
	"fmt"
	"golang.org/x/exp/slog"
	"strings"
)

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
