package main

import (
	"fmt"
	"github.com/baa-god/lan/fib"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var Ws = func(c *fiber.Ctx) (err error) {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}
	return c.Next()
}

type Response struct {
	Event   string `json:"event"`
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func main() {
	app := fib.New()

	var Conn = websocket.New(func(conn *websocket.Conn) {
		b := websocket.FormatCloseMessage(1000, "close_reason")
		err := conn.WriteMessage(websocket.CloseMessage, b)
		fmt.Println("write err:", err)

		defer conn.Close()

		var r Response

		for {
			if err := conn.ReadJSON(&r); err != nil {
				fmt.Println("read err:", err)
				break
			}
		}

		fmt.Println("quit!")
	})

	app.Get("/conn", Ws, Conn)

	app.Listen(":3000")
}
