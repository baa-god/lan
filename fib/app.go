package fib

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Handler any

type Router interface {
	Use(args ...any) Router
	Get(path string, handlers ...Handler) Router
	Post(path string, handlers ...Handler) Router
	Group(path string, handlers ...Handler) Router
	Add(method string, path string, handlers ...Handler) Router
}

type App struct {
	*fiber.App
}

func New() *App {
	app := fiber.New()

	app.Use(
		recover.New(recover.Config{EnableStackTrace: true}),
		logger.New(logger.Config{
			Format:     "${time} |${status}| ${ip} | ${method} ${path} ${queryParams} ${body}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}),
	)

	return &App{App: app}
}

func (app *App) Use(args ...any) Router {
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case string, []string, fiber.Handler:
			continue
		}
		args[i] = HandlerFunc(args[i])
	}

	app.App.Use(args...)
	return app
}

func (app *App) Get(path string, handlers ...Handler) Router {
	return app.Add(fiber.MethodGet, path, handlers...)
}

func (app *App) Post(path string, handlers ...Handler) Router {
	return app.Add(fiber.MethodPost, path, handlers...)
}

func (app *App) Group(path string, handlers ...Handler) Router {
	router := app.App.Group(path, HandlerFunc(handlers...)...)
	return &Group{grp: router.(*fiber.Group)}
}

func (app *App) Add(method string, path string, handlers ...Handler) Router {
	app.App.Add(method, path, HandlerFunc(handlers...)...)
	return app
}
