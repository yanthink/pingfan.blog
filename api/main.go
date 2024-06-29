package main

import (
	"blog/app/events"
	"blog/bootstrap"
	"blog/config"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

func main() {
	if config.App.Env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化日志
	bootstrap.SetupLogger()
	// 初始化数据库
	bootstrap.SetupDatabase()
	// 初始化 Redis
	bootstrap.SetupRedis()
	// 初始化雪花ID
	bootstrap.SetupSnowflake()
	// 初始化事件
	bootstrap.SetupEvent()
	// 运行 websocket 服务
	bootstrap.SetupWebsocket()
	// 定时任务
	cron := bootstrap.SetupCron()

	// 初始化路由
	router := bootstrap.SetupRouter()

	endless.DefaultHammerTime = 1 * time.Second
	if err := endless.ListenAndServe(fmt.Sprintf(":%s", config.App.Port), router); err != nil {
		log.Fatalf("listen: %s\n", err)
	}

	if cron != nil {
		cron.Stop()
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		events.Wg.Wait()
	}()

	finish := make(chan struct{})
	go func() {
		defer close(finish)
		wg.Wait()
	}()

	select {
	case <-finish:
		return
	case <-time.After(60 * time.Second):
		return
	}
}
