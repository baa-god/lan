package jwtware

import (
	"github.com/baa-god/lan/strs"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
)

type Config struct {
	Skip    func(*fiber.Ctx, string) bool
	Succeed func(*fiber.Ctx, jwt.Claims) error
	Claims  jwt.Claims
	Failed  func(*fiber.Ctx, error) error
	Secret  []byte
}

func New(config ...Config) fiber.Handler {
	f := pie.First(config)
	return func(c *fiber.Ctx) (err error) {
		token := c.Get(fiber.HeaderAuthorization)
		if token == "" {
			token = c.Query(fiber.HeaderAuthorization)
		}

		token = strs.TrimPrefix(token, "Bearer ?")
		if f.Skip != nil && f.Skip(c, token) {
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

		if f.Succeed != nil {
			if err = f.Succeed(c, claims); err == nil {
				c.Locals("user", claims)
			}
		}

		return
	}
}
