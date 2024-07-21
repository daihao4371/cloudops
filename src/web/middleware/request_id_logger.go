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

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func NewGinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		//TraceID:    true,

		// Context 的作用可以给zap添加一些k,v对，做一些额外的信息记录
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			//先从响应头里面拿到request_id，如果为空则从请求头里面获取
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			fields = append(fields, zap.String("body", string(body)))
			c.Request.Body = io.NopCloser(&buf)

			// 获取请求头里面的Authorization 记录认证信息 太长了我就关了
			/*		authHeader := c.Request.Header.Get("Authorization")
					fields = append(fields, zap.String("body", string(body)))
					fields = append(fields, zap.String("Authorization", authHeader))*/
			//return fields  // 返回的字段会添加到日志里面
			return fields // 返回的字段不会添加到日志里面
		}),
	})
}
