package sharp

import (
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
