package sharp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

func OpenFile(name string, flag int, perm os.FileMode) (f *File, err error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}
	return &File{file}, err
}

func Or[T any](r1 T, e T) T {
	v := reflect.ValueOf(r1)
	if v.IsValid() && !v.IsZero() {
		return r1
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

func LastSource(file string) string {
	last := strings.LastIndexByte(file, '/')
	index := strings.LastIndexByte(file[:last], '/')
	return file[index+1:]
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

func JwtSigned(key any, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return s
}
