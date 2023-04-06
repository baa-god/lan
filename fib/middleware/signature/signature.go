package signature

import (
	"errors"
	"fmt"
	"github.com/baa-god/lan/lan"
	"github.com/baa-god/lan/strs"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
	"math"
	"strings"
	"sync"
	"time"
)

type Param struct {
	Authorization string `json:"Authorization"`
	Milli         int64  `json:"Milli"`
	Nonce         string `json:"Nonce"`
	Signature     string `json:"Signature"`
	Params        map[string]any
}

var (
	mu    sync.Mutex
	saved = map[string]bool{}
)

func New(secret string, skip func(*fiber.Ctx) bool) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if skip(c) {
			return c.Next()
		}

		token := c.Get("Authorization")
		token = strs.TrimSuffix(token, "Bearer ?")

		p := &Param{
			Authorization: token,
			Milli:         cast.ToInt64(c.Get("Milli")),
			Nonce:         c.Get("Nonce"),
			Signature:     c.Get("Signature"),
			Params:        map[string]any{},
		}

		if strs.HasSuffix(c.Path(), "conn") {
			p.Authorization = c.Query("Authorization")
			p.Milli = cast.ToInt64(c.Query("Milli"))
			p.Nonce = c.Query("Nonce")
			p.Signature = c.Query("Signature")
		}

		args := c.Request().URI().QueryArgs()
		if c.Method() == "POST" {
			args = c.Request().PostArgs()
		}

		args.VisitAll(func(key, value []byte) {
			k, v := string(key), string(value)
			if k == "Authorization" || k == "Nonce" || k == "Milli" || k == "Signature" {
				return
			}
			p.Params[k] = v
		})

		// 验证时间戳
		mill := time.Now().UnixMilli()
		if math.Abs(float64(mill-p.Milli)) > 1000*60 { // 超时1min
			return errors.New("request expired")
		}

		params := lan.CopyMap(map[string]any{
			"Authorization": p.Authorization,
			"Nonce":         p.Nonce,
			"Milli":         p.Milli,
		}, p.Params)

		keys := pie.Sort(maps.Keys(params))
		signs := pie.Map(keys, func(key string) string {
			return fmt.Sprintf("%s=%v", key, params[key])
		})

		sign := strings.Join(append(signs, "KEY="), "&")
		if strs.SHA256(sign+secret) != p.Signature {
			return jwt.ErrTokenSignatureInvalid
		}

		if err != nil {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}

		if _, ok := saved[p.Signature]; ok {
			return c.SendStatus(fiber.StatusForbidden)
		}

		mu.Lock()
		defer mu.Unlock()
		saved[p.Signature] = true

		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			mu.Unlock()
			delete(saved, p.Signature)
		}()

		return c.Next()
	}
}
