package main

import (
	"github.com/baa-god/lan/fib"
)

func main() {
	app := fib.New()

	app.Get("/test", func(c *fib.Ctx) error {
		return c.Status(5000).SendString("hi")
	})

	app.Listen(":3000")
}
