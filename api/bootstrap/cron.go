package bootstrap

import (
	"blog/app"
	"blog/app/resource"
	"blog/app/services"
	"blog/app/storage"
	"blog/app/websocket"
	"blog/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

func SetupCron() *cron.Cron {
	if config.App.SkipCron {
		return nil
	}

	c := cron.New()

	// 每小时同步浏览量并且更新热度
	c.AddFunc("@hourly", func() {
		services.Article.SyncViewCountAndUpdateHotness(time.Now().Add(-5 * time.Minute))
		services.Article.UpdateHotnessInDecayHours()
	})

	// 每天删除前天的上传目录
	c.AddFunc("@daily", func() {
		rType := resource.ArticleImage
		path := rType.UploadPath(time.Now().Add(-25 * time.Hour))
		if err := storage.Disk().DelDir(path); err != nil {
			app.Logger.Debug("上传目录删除失败", zap.Error(err))
		}
	})

	// 每30分钟清理一次心跳超时的连接
	c.AddFunc("@every 30m", func() {
		websocket.ClearTimeoutConnections()
	})

	// 启动定时任务
	c.Start()

	return c
}
