package lan

import (
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
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

func MapTo(m map[string]any, v any) (err error) {
	if b, err := jsoniter.Marshal(m); err == nil {
		err = jsoniter.Unmarshal(b, v)
	}
	return
}

func CopyMap[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) M1 {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func JwtSigned(key any, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return s
}
