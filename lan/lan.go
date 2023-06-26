package lan

import (
	"github.com/bytedance/sonic"
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

func MapTo(m any, v any) (err error) {
	if b, err := sonic.Marshal(m); err == nil {
		err = sonic.Unmarshal(b, v)
	}
	return
}
