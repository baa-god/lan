package sharp

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

type DayDir struct {
	Dir    string
	Format string

	lock bool
	open string
	file *os.File
	mu   sync.Mutex
}

func NewDayDir(dir string, lock bool, perm ...os.FileMode) (d *DayDir, err error) {
	_perm := append(perm, 0)[0]
	if err = os.Mkdir(dir, _perm); err != nil && !errors.Is(err, os.ErrExist) {
		return
	}

	d = &DayDir{Dir: dir, Format: "20060102", lock: lock}
	return
}

func (d *DayDir) Write(b []byte) (n int, err error) {
	format := time.Now().Format(d.Format)
	fileName := strings.TrimRight(d.Dir, "/") + "/day_" + format + ".log"

	if d.open != fileName {
		if d.file != nil {
			_ = d.file.Close()
		}

		d.open = fileName
		d.file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}

	if err == nil {
		if d.lock {
			d.mu.Lock()
			defer d.mu.Unlock()
		}
		n, err = d.file.Write(b)
	}

	return
}
