package middleware

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"time"
)

func TimeCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求获取当前时间
		nowTime := time.Now()
		// 请求处理
		c.Next()

		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		sc.Logger.Info("耗时中间件打印结果",
			zap.String("URL", c.Request.URL.String()),
			zap.Duration("耗时", time.Since(nowTime)),
		)

		log.Printf("the request URL %s const %v", c.Request.URL, time.Since(nowTime))
	}
}
