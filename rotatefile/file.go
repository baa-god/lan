package rotatefile

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type File struct {
	Filename string
	file     *os.File
	perm     os.FileMode

	next   time.Time
	mu     *sync.Mutex
	rotate time.Duration
}

func New(name string, rotate time.Duration, perm os.FileMode, lock bool) (f *File, err error) {
	err = os.MkdirAll(filepath.Dir(name), perm)
	if err == nil || errors.Is(err, os.ErrExist) {
		f = &File{Filename: name, rotate: rotate}
		if err = nil; lock {
			f.mu = &sync.Mutex{}
		}
	}
	return
}

func (f *File) WriteWithLock(b []byte, lock bool) (n int, err error) {
	if lock && f.mu != nil {
		f.mu.Lock()
		defer f.mu.Unlock()
	}

	if now := time.Now(); now.After(f.next) {
		if f.file != nil { // 打开过的文件
			_ = f.file.Close()
		}

		dir := filepath.Dir(f.Filename)
		base := filepath.Base(f.Filename)     // name.ext
		ext := filepath.Ext(f.Filename)       // .ext
		name := strings.TrimSuffix(base, ext) // name

		format := now.Format(time.DateOnly + "_15")
		if f.rotate >= time.Hour*24 {
			format = now.Format(time.DateOnly)
		}

		f.next = now.Add(f.rotate).Truncate(f.rotate)
		f.file, err = os.OpenFile(
			strings.Join([]string{dir, "/", name, "_", format, ext}, ""),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			f.perm,
		)
	}

	if err == nil {
		n, err = f.file.Write(b)
	}

	return
}

func (f *File) Writer(b []byte) (n int, err error) {
	return f.WriteWithLock(b, f.mu != nil)
}
