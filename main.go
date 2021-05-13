package main

import (
	"fmt"

	"github.com/chxfantasy/go_bootstrap/domain/health"
	"github.com/chxfantasy/go_bootstrap/middleware"

	config "github.com/chxfantasy/go_bootstrap/conf"
	"github.com/chxfantasy/go_bootstrap/conf/scheduler"
	"github.com/chxfantasy/go_bootstrap/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello go bootstrap")
	config.Gin.Use(middleware.Boss(config.AppConf.Server.Name, config.TraceLogger))

	initRouter(config.Gin)
	go scheduler.InitScheduler()

	serverPort := config.AppConf.Server.Port
	if serverPort <= 0 {
		serverPort = 7001
	}
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", serverPort)
	config.BizLogger.Info("listen address:", "addr", addr)
	err := utils.ListenAndServe(addr, config.Gin)
	if err != nil {
		config.BizLogger.Errorf("start listener err: %v, addr: %v", err, addr)
		return
	}
}
func initRouter(app *gin.Engine) {
	health.InitRouter(app)
}
