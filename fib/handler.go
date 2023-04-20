package fib

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

var (
	validate     = validator.New()
	fiberCtxType = reflect.TypeOf(&fiber.Ctx{})
)

func HandlerFunc(handlers ...Handler) []fiber.Handler {
	return pie.Map(handlers, func(f Handler) fiber.Handler {
		if handler, ok := f.(fiber.Handler); ok {
			return handler
		}

		fv := reflect.ValueOf(f)       // 具体的处理方法
		args, ctx := fv.Type(), &Ctx{} // 方法类型

		return func(c *fiber.Ctx) (err error) {
			var in []reflect.Value // 处理参数 | 1: 上下文, 2: 输入
			if first := args.In(0); first == fiberCtxType {
				in = append(in, reflect.ValueOf(c))
			} else {
				if ctx.Ctx != c {
					ctx.Ctx = c
					ctx.getArgs()
				}
				in = append(in, reflect.ValueOf(ctx))
			}

			if args.NumIn() == 2 {
				p := reflect.New(args.In(1)) // param-2: 输入
				input := p.Interface()

				if c.Get(fiber.HeaderContentType) == "" {
					if err = c.QueryParser(input); err != nil {
						return err
					}
				} else if err = c.BodyParser(input); err != nil {
					return err
				}

				if err = validate.Struct(input); err != nil { // 验证输入
					return err
				}

				in = append(in, p.Elem())
			}

			return handleResult(c, fv.Call(in)...)
		}
	})
}

func handleResult(c *fiber.Ctx, out ...reflect.Value) (err error) {
	if out == nil {
		return
	}

	st, result := 0, out[0].Interface()
	if result == nil {
		return nil
	}

	if st, _ = result.(int); st > 0 {
		if len(out) == 1 {
			return c.SendStatus(st)
		}
		c.Status(st)
	}

	if len(out) == 2 {
		result = out[1].Interface()
	}

	if err, _ = result.(error); err != nil && st == 0 {
		return err
	}

	fmt.Println("accept --------------:", c.Get("Accept"))

	if accept := c.Get("Accept"); accept != "" && accept != "*/*" {
		return c.Format(result)
	}

	return c.JSON(result)
}
