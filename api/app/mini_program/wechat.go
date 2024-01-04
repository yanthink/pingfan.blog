package mini_program

import (
	"blog/app"
	"blog/app/cache"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Wechat struct {
	Appid  string
	Secret string
}

func (w *Wechat) Code2Session(code string, anonymousCode string) Session {
	params := url.Values{
		"appid":      {w.Appid},
		"secret":     {w.Secret},
		"js_code":    {code},
		"grant_type": {"authorization_code"},
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?%s", params.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		SessionKey string `json:"session_key"`
		Unionid    string `json:"unionid"`
		Openid     string `json:"openid"`
		Errcode    int64  `json:"errcode"`
		Errmsg     string `json:"errmsg"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	if result.Errcode != 0 {
		panic(result.Errmsg)
	}

	return Session{
		SessionKey: result.SessionKey,
		Unionid:    result.Unionid,
		Openid:     result.Openid,
	}
}

type GetUnlimitedQRCodeParams struct {
	Scene      string          `json:"scene"`
	Page       string          `json:"page,omitempty"`
	CheckPath  bool            `json:"check_path"`
	EnvVersion string          `json:"env_version,omitempty"`
	Width      int64           `json:"width,omitempty"`
	AutoColor  bool            `json:"auto_color,omitempty"`
	LineColor  *map[string]any `json:"line_color,omitempty"`
	IsHyaline  bool            `json:"is_hyaline"`
}

func (w *Wechat) GetUnlimitedQRCode(params *GetUnlimitedQRCodeParams) (buf bytes.Buffer) {
	gatewayUrl := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", w.getAccessToken())
	jsonParams, _ := json.Marshal(params)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(gatewayUrl, "application/json", bytes.NewBuffer(jsonParams))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(&buf, resp.Body); err != nil {
		panic(err)
	}

	var result struct {
		Errcode int64  `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err == nil && result.Errcode != 0 {
		if result.Errcode == 40001 {
			w.clearAccessToken() // todo retry
		}
		app.Logger.Sugar().Debug("获取小程序码失败：", params, result)
		panic(result.Errmsg)
	}

	return
}

func (w *Wechat) MsgSecCheck(content, openid string) bool {
	gatewayUrl := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", w.getAccessToken())

	params := map[string]any{
		"content": content,
		"version": 2,
		"scene":   2, // 评论
		"openid":  openid,
	}

	jsonParams, _ := json.Marshal(params)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(gatewayUrl, "application/json", bytes.NewBuffer(jsonParams))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Errcode int64            `json:"errcode"`
		Errmsg  string           `json:"errmsg"`
		Detail  []map[string]any `json:"detail"`
		Result  struct {
			Suggest string `json:"suggest"`
			Label   int64  `json:"label"`
		} `json:"result"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	if result.Errcode != 0 {
		if result.Errcode == 40001 {
			w.clearAccessToken() // todo retry
		}
		panic(result.Errmsg)
	}

	if result.Result.Suggest != "pass" {
		app.Logger.Sugar().Debugf("内容安全检测未通过：%+v", result)
	}

	return result.Result.Suggest == "pass"
}

func (w *Wechat) getAccessToken() (accessToken string) {
	strCache := cache.New[string]()
	cacheKey := fmt.Sprintf("access_token:%s", w.Appid)
	ctx := context.Background()

	var err error

	if accessToken, err = strCache.Get(ctx, cacheKey); err == nil {
		return
	}

	params := url.Values{
		"appid":      {w.Appid},
		"secret":     {w.Secret},
		"grant_type": {"client_credential"},
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?%s", params.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Errcode     int64  `json:"errcode"`
		Errmsg      string `json:"errmsg"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	if result.Errcode != 0 {
		panic(result.Errmsg)
	}

	accessToken = result.AccessToken
	_ = strCache.Put(ctx, cacheKey, accessToken, time.Duration(result.ExpiresIn-15)*time.Second)

	return
}

func (w *Wechat) clearAccessToken() {
	strCache := cache.New[string]()
	cacheKey := fmt.Sprintf("access_token:%s", w.Appid)
	ctx := context.Background()

	_ = strCache.Forget(ctx, cacheKey)
}
