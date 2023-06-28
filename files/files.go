package files

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// LastLine 读取 file 最后一行文本
func LastLine(f *os.File) (line string, err error) {
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

// FirstLine 读取 file 第一行文本
func FirstLine(f *os.File) (line string, err error) {
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return
	}

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan() && line == "" {
		line = strings.TrimSpace(scan.Text())
	}

	return
}
