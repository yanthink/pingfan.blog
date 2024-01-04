package bootstrap

import (
	"blog/app/http/middleware"
	"blog/config"
	"blog/routes"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	_ = r.SetTrustedProxies(config.App.Proxies)

	if config.App.Env == "local" && config.App.Debug {
		r.Use(gin.Logger())
	}

	r.Use(middleware.Cors())
	r.Use(middleware.ErrorHandler())

	routes.RegisterApiRouter(r)
	routes.RegisterWebsocketRouter(r)

	return r
}
