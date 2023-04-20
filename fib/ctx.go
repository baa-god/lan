package fib

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"strconv"
)

type Ctx struct {
	*fiber.Ctx
	args map[string]any
}

func (c *Ctx) Arg(key string, or ...string) string {
	v, ok := c.args[key]
	if !ok && or != nil {
		return fmt.Sprint(v)
	}
	return fmt.Sprint(v)
}

func (c *Ctx) ArgInt(key string, or ...int) int {
	v, err := strconv.Atoi(c.Arg(key))
	if err != nil && or != nil {
		v = or[0]
	}
	return v
}

func (c *Ctx) ArgInt64(key string, or ...int64) int64 {
	v, err := strconv.ParseInt(c.Arg(key), 10, 64)
	if err != nil && or != nil {
		v = or[0]
	}
	return v
}

func (c *Ctx) ArgBool(key string) bool {
	return cast.ToBool(c.Arg(key))
}

func (c *Ctx) getArgs() {
	if c.args = map[string]any{}; c.Method() == "GET" {
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			c.args[string(key)] = string(value)
		})
	} else if c.Method() == "POST" {
		_ = c.App().Config().JSONDecoder(c.Body(), &c.args)
	}
}
