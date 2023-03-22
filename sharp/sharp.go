package sharp

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	Digits         = "0123456789"
	AsciiLowercase = "abcdefghijklmnopqrstuvwxyz"
	AsciiUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Punctuation    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Asciis         = Digits + AsciiLowercase + AsciiUppercase + Punctuation
)

type Pair[Key any, Value any] struct {
	First  Key
	Second Value
}

func Or[T any](r1 T, e T) T {
	v := reflect.ValueOf(r1)
	if v.IsValid() && !v.IsZero() {
		return r1
	}
	return e
}

func If[T any](b bool, v T, e T) T {
	if b {
		return v
	}
	return e
}

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

func Rand(n int) (s string) {
	length := len(Asciis)
	for i := 0; i < n; i++ {
		// bigI: [0, length-1]
		bigI, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		index, _ := strconv.Atoi(bigI.String())
		s += fmt.Sprintf("%c", Asciis[index])
	}
	return
}

func Sha256(text string) string {
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	return strings.ToUpper(sum)
}

func HasSuffix(s string, a ...string) bool {
	for _, v := range a {
		if strings.HasSuffix(s, v) {
			return true
		}
	}
	return false
}
