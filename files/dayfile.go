package files

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

type Dayfile struct {
	Filename string
	file     *os.File
	format   string
	perm     os.FileMode

	time string
	mu   *sync.Mutex
}

func NewDayfile(name, format string, perm os.FileMode, lock bool) (f *Dayfile, err error) {
	if err = os.Mkdir(name, perm); err == nil || errors.Is(err, os.ErrExist) {
		f = &Dayfile{Filename: name, format: format, time: time.Now().Format(format)}
		if lock {
			f.mu = &sync.Mutex{}
		}
	}
	return
}

func (f *Dayfile) Write(b []byte) (n int, err error) {
	return f.WriteWithLock(b, f.mu != nil)
}

func (f *Dayfile) WriteWithLock(b []byte, lock bool) (n int, err error) {
	_time := time.Now().Format(f.format)
	name := strings.Join([]string{f.Filename, "_", _time, ".log"}, "")

	if f.time != _time {
		if lock && f.mu != nil {
			f.mu.Lock()
			defer f.mu.Unlock()
		}

		if f.file != nil { // 打开过的文件
			_ = f.file.Close()
		}

		f.time = _time
		f.file, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, f.perm)
		if err == nil {
			n, err = f.file.Write(b)
		}
	}

	return
}
