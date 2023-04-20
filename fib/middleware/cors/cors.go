package cors

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func New(allowOrigins, allowHeaders string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var allow bool
		origin := c.Get(fiber.HeaderOrigin)

		for _, x := range strings.Split(allowOrigins, ",") {
			if x == origin {
				allow = true
				break
			}
		}

		if origin != "" && !allow {
			return c.SendStatus(fiber.StatusForbidden)
		}

		allowHeaders = strings.ReplaceAll(allowHeaders, " ", "")
		allowOrigins = strings.ReplaceAll(allowOrigins, " ", "")

		c.Set(fiber.HeaderAccessControlAllowOrigin, allowOrigins)
		c.Set(fiber.HeaderAccessControlAllowHeaders, allowHeaders)

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
