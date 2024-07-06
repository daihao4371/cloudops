package middleware

import (
	"bytes"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

func NewGinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,

		// Context 的作用可以给zap添加一些k,v对，做一些额外的信息记录
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			//先从响应头里面拿到request_id，如果为空则从请求头里面获取
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			//从请求体里面获取body  在post请求的时候，查看参数格式
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			// 获取请求头里面的Authorization
			authHeader := c.Request.Header.Get("Authorization")
			fields = append(fields, zap.String("body", string(body)))
			fields = append(fields, zap.String("Authorization", authHeader))
			return fields
		}),
	})
}
