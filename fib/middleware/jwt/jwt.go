package jwtware

import (
	"github.com/baa-god/lan/strs"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
)

type Config struct {
	Skip    func(*fiber.Ctx) bool
	Succeed func(*fiber.Ctx, jwt.Claims) error
	Claims  jwt.Claims
	Failed  func(*fiber.Ctx, error) error
	Secret  []byte
}

func New(config ...Config) fiber.Handler {
	f := pie.First(config)
	return func(c *fiber.Ctx) (err error) {
		if f.Skip != nil && f.Skip(c) {
			return c.Next()
		}

		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			token = c.Query(fiber.HeaderAuthorization)
		}

		if token = strs.TrimPrefix(token, "Bearer ?"); token == "" {
			return c.Next()
		}

		value := reflect.New(reflect.TypeOf(f.Claims))
		claims := value.Interface().(jwt.Claims)

		_, err = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) {
			return f.Secret, nil
		})

		if err != nil {
			if f.Failed == nil {
				return c.Status(403).SendString(err.Error())
			}
			return f.Failed(c, err)
		}

		if c.Locals("user", claims); f.Succeed == nil {
			return c.Next()
		}

		return f.Succeed(c, claims)
	}
}
