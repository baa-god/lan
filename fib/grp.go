package fib

import (
	"github.com/gofiber/fiber/v2"
)

type Group struct {
	grp *fiber.Group
}

func (grp *Group) Use(args ...any) Router {
	for i := 0; i < len(args); i++ {
		switch args[i].(type) {
		case string, []string, fiber.Handler:
			continue
		}
		args[i] = HandlerFunc(args[i])
	}

	grp.grp.Use(args...)
	return grp
}

func (grp *Group) Get(path string, handlers ...Handler) Router {
	return grp.Add(fiber.MethodGet, path, handlers...)
}

func (grp *Group) Post(path string, handlers ...Handler) Router {
	return grp.Add(fiber.MethodPost, path, handlers...)
}

func (grp *Group) Group(path string, handlers ...Handler) Router {
	router := grp.grp.Group(path, HandlerFunc(handlers...)...)
	return &Group{router.(*fiber.Group)}
}

func (grp *Group) Add(method string, path string, handlers ...Handler) Router {
	grp.grp.Add(method, path, HandlerFunc(handlers...)...)
	return grp
}
