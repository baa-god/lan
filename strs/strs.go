package strs

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/baa-god/lan/lan"
	"math/big"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func BaseN(path string, n int) string {
	var noHas bool
	var base string
	path = filepath.ToSlash(path)

	for i, index := 0, 0; n > i && index != -1; i++ {
		if index = strings.LastIndexByte(path, '/'); index > 0 {
			base = path[index:] + base
			path = path[:index]
			continue
		}
		noHas = true
	}

	if noHas {
		return path
	}

	return strings.TrimPrefix(base, "/")
}

func Rand(n int) (s string) {
	length := len(lan.Asciis)
	for i := 0; i < n; i++ {
		max := big.NewInt(int64(length)) // [0, length-1]
		bigI, _ := rand.Int(rand.Reader, max)
		index, _ := strconv.Atoi(bigI.String())
		s += fmt.Sprintf("%c", lan.Asciis[index])
	}
	return
}

func SHA256(text string) string {
	b := sha256.Sum256([]byte(text))
	return strings.ToUpper(fmt.Sprintf("%x", b))
}

func HasSuffix(s string, a ...string) bool {
	for _, v := range a {
		if strings.HasSuffix(s, v) {
			return true
		}
	}
	return false
}

func TrimPrefix(s string, pattern string) string {
	re := regexp.MustCompile("^" + pattern)
	return re.ReplaceAllString(s, "")
}

func TrimSuffix(s string, pattern string) string {
	re := regexp.MustCompile(pattern + "$")
	return re.ReplaceAllString(s, "")
}
