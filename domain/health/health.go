package health

import (
	"net/http"

	"github.com/chxfantasy/go_bootstrap/conf"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(app *gin.Engine) {
	app.Any("/health", func(ctx *gin.Context) {
		conf.BizLogger.Info("health")
		ctx.String(http.StatusOK, "SUCCESS")
	})

	app.Any("/health2", func(ctx *gin.Context) {
		conf.BizLogger.Info("health2 with trace log")
		ctx.String(http.StatusOK, "SUCCESS")
	})
}
