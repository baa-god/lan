package fib

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"net/url"
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

func NewCtx(ctx *fiber.Ctx) *Ctx {
	c := &Ctx{Ctx: ctx, args: map[string]any{}}

	if c.Method() == "GET" {
		query, _ := url.ParseQuery(c.Context().QueryArgs().String())
		for k, v := range query {
			c.args[k] = v[0]
		}
	} else if c.Method() == "POST" {
		dec := sonic.ConfigFastest.NewDecoder(bytes.NewReader(c.Body()))
		dec.UseNumber()
		_ = dec.Decode(&c.args)
	}

	return c
}

func (c *Ctx) Abort() error {
	c.Route().Handlers = nil
	return nil
}
