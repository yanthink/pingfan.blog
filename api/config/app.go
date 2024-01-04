package config

import (
	"encoding/hex"
	"github.com/spf13/viper"
)

type app struct {
	Env      string
	Debug    bool
	Name     string
	Port     string
	Key      []byte
	Url      string
	SiteUrl  string
	Proxies  []string
	SkipCron bool
}

var App app

func loadAppConfig() {
	// UnmarshalKey 无法读取环境变量配置
	// _ = viper.UnmarshalKey("app", &App)

	hexKey := GetString("app.key")
	key, _ := hex.DecodeString(hexKey)

	App = app{
		Env:      GetString("app.env"),
		Debug:    GetBool("app.debug", false),
		Name:     GetString("app.name"),
		Port:     GetString("app.port", "8888"),
		Key:      key,
		Url:      GetString("app.url"),
		SiteUrl:  GetString("app.site_url"),
		Proxies:  viper.GetStringSlice("app.proxies"),
		SkipCron: GetBool("app.skip_cron", false),
	}
}
