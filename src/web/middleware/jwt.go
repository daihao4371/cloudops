package middleware

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 在header里面去获取authorization
		authHeaderString := c.Request.Header.Get("Authorization")
		if authHeaderString == "" {
			common.Req401WithWithDetailed(gin.H{"reload": true}, "未登录或非法登录没有Authorization", c)
			c.Abort()
			return
		}
		// 获取到authorization校验
		parts := strings.SplitN(authHeaderString, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.Req401WithWithDetailed(gin.H{"reload": true}, "请求头中的auth格式错误", c)
			// 阻止调用后续的函数
			c.Abort()
			return
		}
		//解析token字符串jwt解析
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		userClaims, err := models.ParseToken(parts[1], sc)
		if err != nil {
			common.Req401WithWithDetailed(gin.H{"reload": true}, fmt.Sprintf("parseToken 解析token包含信息错误：%v", err.Error()), c)
			// 阻止调用后续的函数
			c.Abort()
			return
		}
		// token续期的逻辑
		// 判断token是否过期，定义1个临时窗口
		// 重新生成一个token, 返回给前端
		if userClaims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix() < int64(sc.JWTC.BufferDuration/time.Second) {
			sc.Logger.Info("token即将过期，重新生成token",
				zap.String("user", userClaims.Username),
				zap.String("老token过期时间", userClaims.RegisteredClaims.ExpiresAt.String()),
				zap.String("临期窗口 ", sc.JWTC.BufferDuration.String()),
			)
			// 重新生成token
			newToken, err := models.GenJwtToken(userClaims.User, sc)
			if err != nil {
				common.Result5XX(0, gin.H{}, fmt.Sprintf("parseToken 解析token信息错误：%v", err.Error()), c)
				// 阻止调用后续的函数
				c.Abort()
			}
			// 设置到header中
			//c.Header("Authorization", "Bearer "+newToken)
			c.Header("new-token", newToken)
		} else {
			/*			sc.Logger.Info("token未过期，无需重新生成token",
						zap.String("user", userClaims.Username),
						zap.String("老token过期时间", userClaims.RegisteredClaims.ExpiresAt.String()),
						zap.String("临期窗口 ", sc.JWTC.BufferDuration.String()),
					)*/
		}

		// 设置cliam 后续的处理逻辑就能拿到
		c.Set(common.GIN_CTX_JWT_USER_NAME, userClaims.Username)
		c.Next()

	}
}
