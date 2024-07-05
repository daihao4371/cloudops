package web

import (
	"cloudops/src/config"
	"cloudops/src/web/view"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// StartGIn 启动gin
func StartGIn(sc *config.ServeConfig) error {
	// 初始化引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 禁用控制台颜色
	gin.DisableConsoleColor()
	// 初始化路由
	view.ConfigRouter(r)

	// http 读写超时参数
	s := &http.Server{
		Addr:           sc.HttpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
