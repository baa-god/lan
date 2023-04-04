package lan

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
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

func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) M1 {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func Unmarshal(data any, output any) (err error) {
	if v, err := json.Marshal(data); err == nil {
		err = json.Unmarshal(v, output)
	}
	return
}

func JwtSigned(key any, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return s
}
