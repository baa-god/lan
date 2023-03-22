package sharp

import (
	"io"
	"os"
	"strings"
)

// ReadLastLine 读取 file 最后一行文本
func ReadLastLine(file *os.File) (line string, err error) {
	offset, err := file.Seek(-1, io.SeekEnd)

	for buffer := make([]byte, 1); offset >= 0; {
		_, err = file.ReadAt(buffer, offset)
		offset -= 1            // 负偏移索引
		char := string(buffer) // 每次读取的单个字符

		// 如果还没有获取到字符则丢弃行尾空格
		if line == "" && strings.TrimSpace(char) == "" {
			continue
		}

		if char == "\n" { // 找到最后的换行符后退出
			break
		}

		line = char + line
	}

	return
}
