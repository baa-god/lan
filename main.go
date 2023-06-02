package main

import (
	"fmt"
	"github.com/baa-god/lan/fib"
	"github.com/gofiber/fiber/v2"
)

func main() {
	r := fib.New()
	r.Post("test/post", func(c *fib.Ctx) {
		fmt.Println(c.Arg("name"))
		fmt.Println(c.Get("name"))
	})

	r.Get("test/get", func(c *fib.Ctx) {
		c.Route().Handlers = []fiber.Handler{}
		fmt.Println("================")
	}, func(c *fib.Ctx) (err error) {
		fmt.Println("argName:", c.Query("name"))
		err = c.SendString("6666")
		return err
	})

	r.Listen(":800")
}
