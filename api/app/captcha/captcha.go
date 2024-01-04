package captcha

import (
	"blog/app/cache"
	"blog/config"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"time"
)

const Email = "email"

type Captcha struct {
	Type      string
	Expire    time.Duration
	Sensitive bool
}

func New() *Captcha {
	return &Captcha{
		Type:      Email,
		Expire:    5 * time.Minute,
		Sensitive: false,
	}
}

func (c *Captcha) SetType(t string) *Captcha {
	c.Type = t

	return c
}

func (c *Captcha) SetExpire(duration time.Duration) *Captcha {
	c.Expire = duration

	return c
}

func (c *Captcha) SetSensitive(sensitive bool) *Captcha {
	c.Sensitive = sensitive

	return c
}

func (c *Captcha) Generate(account string) (code string, hash string) {
	code = random.RandNumeral(6)
	hash = random.RandString(16)

	if config.App.Env == "local" && config.App.Debug {
		code = "888888"
	}

	strCache := cache.New[string]()

	key := c.buildKey(account)
	_ = strCache.Set(context.Background(), key, code+hash, c.Expire)

	return
}

func (c *Captcha) Check(account string, code string, hash string) bool {
	strCache := cache.New[string]()

	ctx := context.Background()
	key := c.buildKey(account)

	if v, _ := strCache.Pull(ctx, key); v == code+hash {
		return true
	}

	return false
}

func (c *Captcha) buildKey(account string) string {
	return fmt.Sprintf("%s_captcha:%s", c.Type, account)
}
