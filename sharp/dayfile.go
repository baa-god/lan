package sharp

import (
	"context"
	"fmt"
	"github.com/gookit/color"
	"golang.org/x/exp/slog"
	"os"
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
	fileName := strings.TrimRight(f.Dir, "/") + "_day" + format + ".log"

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

func (f *DayHandler) Handle(ctx context.Context, r slog.Record) error {
	level := "[" + r.Level.String() + "]"
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
		fields += fmt.Sprintf("%s=", a.Key)
		if _, ok := a.Value.Any().(string); ok {
			fields += `"` + a.Value.String() + `" `
		} else {
			fields += a.Value.String() + " "
		}
	})

	prefix := r.Time.Format("[2006-01-02 15:05:05.000]")
	prefix = color.HEX("7970A9").Sprint(prefix)

	msg := color.White.Sprint(r.Message)
	fields = color.Cyan.Sprint(fields)

	fmt.Println(prefix+" "+level, msg, fields)
	return nil
}

func NewDayHandle(dir string) *DayHandler {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				v := a.Value.String()
				last := strings.LastIndexByte(v, '/')
				if index := strings.LastIndex(v[:last], "/"); index != -1 {
					v = v[index+1:]
				} else {
					v = v[last:]
				}
				a.Value = slog.StringValue(v)
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
