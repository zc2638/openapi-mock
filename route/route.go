package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zctod/go-tool/common/util_server"
	"mock/config"
	"mock/middleware"
	"net/http"
	"time"
)

/**
 * Created by zc on 2019-10-24.
 */
func Start() {

	g := gin.Default()

	// 注册全局中间件
	g.Use(middleware.Cors())

	// 设置静态路由
	g.Static("/uploads", "uploads")

	// 注册路由
	routeApi(g)

	// 启动服务
	startServer(g)
}

func startServer(g *gin.Engine) {

	server := &http.Server{
		Addr:           "127.0.0.1:" + config.Cfg.Port,
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// 平滑退出，先结束所有在执行的任务
	util_server.GracefulExitWeb(server)
}