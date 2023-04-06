package fib

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

type Ctx struct {
	*fiber.Ctx
}

func (c *Ctx) Form(key string, or ...string) (v string) {
	first := pie.FirstOr(or, v)
	if c.Get(fiber.HeaderContentType) == "" {
		return c.Query(key, first)
	}
	return c.FormValue(key, first)
}

func (c *Ctx) FormInt(key string, or ...int) int {
	if or == nil {
		or = append(or, 0)
	}
	return cast.ToInt(c.Form(key, fmt.Sprint(or[0])))
}

func (c *Ctx) FormInt64(key string, or ...int64) int64 {
	if or == nil {
		or = append(or, 0)
	}
	return cast.ToInt64(c.Form(key, fmt.Sprint(or[0])))
}

func (c *Ctx) FormBool(key string) bool {
	return cast.ToBool(c.Form(key))
}
