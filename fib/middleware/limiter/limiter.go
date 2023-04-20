package limiter

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Config struct {
	Limit    rate.Limit
	Burst    int
	Duration time.Duration
	Skip     func(*fiber.Ctx) bool
	Key      func(*fiber.Ctx) string
	Reached  func(*fiber.Ctx, time.Duration) error
	Add429   bool
}

var (
	mu        sync.Mutex
	blacklist = map[string]*Limiter{}
)

type Limiter struct {
	*rate.Limiter
	Time time.Time
}

func New(config ...Config) fiber.Handler {
	if config == nil {
		config = append(config, Config{})
	}

	f := config[0]
	return func(c *fiber.Ctx) error {
		if f.Skip != nil && f.Skip(c) {
			return c.Next()
		}

		key := f.Key(c)
		limit, _ := blacklist[key]

		if limit == nil {
			mu.Lock()
			defer mu.Unlock()

			limit = &Limiter{
				Limiter: rate.NewLimiter(f.Limit, f.Burst),
				Time:    time.Now(),
			}

			go func() {
				time.Sleep(f.Duration)
				mu.Lock()
				mu.Unlock()
				delete(blacklist, key)
			}()

			blacklist[key] = limit
		}

		if !limit.Allow() {
			if limit.SetBurst(0); f.Add429 {
				c.Status(fiber.StatusTooManyRequests)
			}
			return f.Reached(c, f.Duration-time.Since(limit.Time))
		}

		return c.Next()
	}
}
