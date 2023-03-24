package sharp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type DayFile struct {
	Dir    string
	Format string

	mu       sync.Mutex
	file     *os.File
	lastOpen string
}

type DayHandler struct {
	slog.Handler
}

func (f *DayFile) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	format := time.Now().Format(f.Format)
	fileName := strings.TrimRight(f.Dir, "/") + "/day_" + format + ".log"

	if f.lastOpen != fileName {
		if f.file != nil {
			_ = f.file.Close()
		}
		f.lastOpen = fileName
		f.file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}

	if err == nil {
		n, err = f.file.Write(b)
	}

	return
}

func (f *DayHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	if err = f.Handler.Handle(ctx, r); err != nil {
		return err
	}

	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = color.Magenta.Sprint(level)
	case slog.LevelInfo:
		level = color.Green.Sprint(level)
	case slog.LevelWarn:
		level = color.Yellow.Sprint(level)
	case slog.LevelError:
		level = color.Red.Sprint(level)
	}

	var fields string
	r.Attrs(func(a slog.Attr) {
		var b []byte
		if b, err = json.Marshal(a.Value.Any()); err != nil {
			return
		}
		fields += fmt.Sprintf("%s=%v ", a.Key, string(b))
	})

	prefix := r.Time.Format("2006-01-02 15:05:05.000")
	prefix = color.HEX("A9B7C6").Sprint(prefix)

	_, file, line, _ := runtime.Caller(3)
	source := LastSource(fmt.Sprint(file, ":", line))

	fields = color.Cyan.Sprint(fields)
	fmt.Println(prefix, "|", level, "|", source, ">", r.Message, fields)

	return
}

func NewDayHandle(dir string) *DayHandler {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Value = slog.StringValue(LastSource(a.Value.String()))
			} else if a.Key == slog.TimeKey {
				value := time.Now().Format("2006-01-02 15:04:05.000")
				a.Value = slog.StringValue(value)
			}
			return a
		},
	}

	w := &DayFile{Dir: dir, Format: "060102"}
	return &DayHandler{Handler: opts.NewJSONHandler(w)}
}
