package sharp

import (
	"os"
	"sync"
	"time"
)

type DayFile struct {
	mu     sync.Mutex
	Dir    string
	Format string
}

func (f *DayFile) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileName := f.Dir + "/day_" + time.Now().Format(f.Format) + ".log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err == nil {
		n, err = file.Write(b)
	}

	return
}
