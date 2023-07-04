package sign

import (
	"bytes"
	"fmt"
	"github.com/baa-god/lan/strs"
	"github.com/bytedance/sonic"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
	"time"

	"sync"
	// "time"
)

type Param struct {
	Auth   string `json:"Authorization"`
	Milli  int64  `json:"Milli"`
	Nonce  string `json:"Nonce"`
	Signed string `json:"Signed"`
	Params map[string]any
}

var (
	mu     sync.Mutex
	caches = map[string]bool{}
)

func New(secrets []string, retSecret func(*fiber.Ctx, string)) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := Param{
			Auth:   strs.TrimPrefix(c.Get("Authorization"), "Bearer ?"),
			Milli:  cast.ToInt64(c.Get("Milli")),
			Nonce:  c.Get("Nonce"),
			Signed: fmt.Sprint(c.Get("Signed")),
			Params: map[string]any{},
		}

		if p.Signed == "" {
			p.Auth = strs.TrimPrefix(c.Query("Authorization"), "Bearer ?")
			p.Milli = cast.ToInt64(c.Query("Milli"))
			p.Nonce = c.Query("Nonce")
			p.Signed = fmt.Sprint(c.Query("Signed"))
		}

		args := map[string]any{}
		if c.Method() == "GET" {
			c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
				args[string(key)] = string(value)
			})
		} else if c.Method() == "POST" {
			dec := sonic.ConfigFastest.NewDecoder(bytes.NewReader(c.Body()))
			dec.UseNumber()
			_ = dec.Decode(&args)
		}

		for k, v := range args {
			if k == "Authorization" || k == "Nonce" || k == "Milli" || k == "Signed" {
				continue
			}
			p.Params[k] = v
		}

		// 验证时间戳
		// msec := time.Now().UnixMilli()
		// if math.Abs(float64(msec-p.Milli)) > 1000*60 { // 超时1min
		// 	return c.Status(fiber.StatusForbidden).SendString("request expired")
		// }

		params := map[string]any{"Authorization": p.Auth, "Nonce": p.Nonce, "Milli": p.Milli}
		maps.Copy(params, p.Params)

		var allow bool
		var useSecret string
		var sortedArgs string
		sortedKeys := pie.Sort(maps.Keys(params))

		for _, key := range sortedKeys {
			sortedArgs += fmt.Sprintf("%s=%v&", key, params[key])
		}

		sortedArgs += "KEY="
		for _, x := range secrets {
			if allow = strs.SHA256(sortedArgs+x) == p.Signed; allow {
				useSecret = x
				break
			}
		}

		if _, ok := caches[p.Signed]; ok || !allow {
			return c.Status(403).SendString("sign error!")
		}

		mu.Lock()
		defer mu.Unlock()
		caches[p.Signed] = true

		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			defer mu.Unlock()
			delete(caches, p.Signed)
		}()

		retSecret(c, useSecret)
		return c.Next()
	}
}
