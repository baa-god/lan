package sharp

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type File struct {
	*os.File
}

func OpenFile(name string, flag int, perm os.FileMode) (f *File, err error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}
	return &File{file}, err
}

// ReadLastLine 读取 file 最后一行文本
func (f *File) ReadLastLine() (line string, err error) {
	offset, err := f.Seek(-1, io.SeekEnd)

	for buffer := make([]byte, 1); offset >= 0; offset-- {
		if _, err = f.ReadAt(buffer, offset); err != nil {
			return
		}

		char := string(buffer)                           // 每次读取的单个字符
		if line == "" && strings.TrimSpace(char) == "" { // 如果没获取到字符则丢弃行尾空格
			continue
		}

		if char == "\n" || char == "\r" { // 找到最后的换行符后退出
			break
		}

		line = char + line
	}

	return
}

func (f *File) ReadFirstLine() (line string, err error) {
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan() && line == "" {
		line = strings.TrimSpace(scan.Text())
	}

	return
}
