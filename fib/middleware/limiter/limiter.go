package limiter

import (
	"github.com/baa-god/lan/lan"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
	"time"
)

type Config struct {
	Limit    rate.Limit
	Burst    int
	Duration time.Duration
	Skip     func(*fiber.Ctx) bool
	Pre      func(c *fiber.Ctx) bool
	Key      func(*fiber.Ctx) string
	Reached  func(*fiber.Ctx, time.Duration) error
}

type Limiter struct {
	*rate.Limiter
	Time time.Time
}

var caches = lan.SyncMap[string, *Limiter]{}

func New(config ...Config) fiber.Handler {
	if config == nil {
		config = append(config, Config{})
	}

	f := config[0]
	return func(c *fiber.Ctx) error {
		if f.Skip != nil && f.Skip(c) {
			return c.Next()
		}

		if f.Pre != nil && f.Pre(c) {
			return nil
		}

		key := f.Key(c)
		limit, _ := caches.Load(key)

		if limit == nil {
			limit = &Limiter{
				Limiter: rate.NewLimiter(f.Limit, f.Burst),
				Time:    time.Now(),
			}

			caches.Store(key, limit)
			time.AfterFunc(f.Duration, func() {
				caches.Delete(key)
			})
		}

		if !limit.Allow() {
			limit.SetBurst(0)
			c.Status(fiber.StatusTooManyRequests)
			return f.Reached(c, f.Duration-time.Since(limit.Time))
		}

		return c.Next()
	}
}
