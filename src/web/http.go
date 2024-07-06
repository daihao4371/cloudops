package web

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/web/middleware"
	"cloudops/src/web/view"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// StartGIn 启动gin
// 设置中间件
func StartGIn(sc *config.ServeConfig) error {
	// 初始化引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 把logger插入
	varMap := map[string]interface{}{}
	//varMap[common.GIN_CTX_CONFIG_LOGGER] = sc.Logger
	varMap[common.GIN_CTX_CONFIG_CONFIG] = sc

	// 添加中间件 打印请求耗时
	r.Use(middleware.TimeCost())
	// 传递变量
	r.Use(middleware.ConfigMiddleware(varMap))

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
