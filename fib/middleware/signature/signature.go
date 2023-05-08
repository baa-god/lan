package signature

import (
	"bytes"
	"fmt"
	"github.com/baa-god/lan/lan"
	"github.com/baa-god/lan/strs"
	"github.com/elliotchance/pie/v2"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
	"math"
	"time"

	"strings"
	"sync"
	// "time"
)

type Param struct {
	Authorization string `json:"Authorization"`
	Milli         int64  `json:"Milli"`
	Nonce         string `json:"Nonce"`
	Signed        string `json:"Signed"`
	Params        map[string]any
}

var (
	mu    sync.Mutex
	saved = map[string]bool{}
)

func New(secret string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := Param{
			Authorization: strs.TrimPrefix(c.Get("Authorization"), "Bearer ?"),
			Milli:         cast.ToInt64(c.Get("Milli")),
			Nonce:         c.Get("Nonce"),
			Signed:        fmt.Sprint(c.Get("Signed")),
			Params:        map[string]any{},
		}

		if p.Signed == "" {
			p.Authorization = strs.TrimPrefix(c.Query("Authorization"), "Bearer ?")
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
			dec := jsoniter.NewDecoder(bytes.NewReader(c.Body()))
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
		msec := time.Now().UnixMilli()
		if math.Abs(float64(msec-p.Milli)) > 1000*60 { // 超时1min
			return c.Status(fiber.StatusForbidden).SendString("request expired")
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

		signed := strings.Join(append(signs, "KEY="), "&")

		if strs.SHA256(signed+secret) != p.Signed {
			fmt.Println("path:", c.Path())
			fmt.Println("signed+secret:", signed+secret)
			return c.Status(fiber.StatusForbidden).SendString("signed invalid")
		}

		if err != nil {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}

		if _, ok := saved[p.Signed]; ok {
			return c.Status(fiber.StatusForbidden).SendString("signed expired")
		}

		mu.Lock()
		defer mu.Unlock()
		saved[p.Signed] = true

		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			mu.Unlock()
			delete(saved, p.Signed)
		}()

		return c.Next()
	}
}
