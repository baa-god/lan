package lan

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"path/filepath"
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

func Or[T any](v T, e T) T {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Invalid &&
		value.IsValid() && !value.IsZero() {
		return v
	}
	return e
}

func ErrOr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func If[T any](b bool, v T, e T) T {
	if b {
		return v
	}
	return e
}

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

func JsonUnmarshal(data any, output any) (err error) {
	if v, err := json.Marshal(data); err == nil {
		err = json.Unmarshal(v, output)
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

func TrimPrefix(s string, prefix ...string) string {
	for _, x := range prefix {
		s = strings.TrimPrefix(s, x)
	}
	return s
}

func JwtSigned(key any, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return s
}
