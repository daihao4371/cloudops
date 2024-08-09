package middleware

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CasbinRbacMiddleware 权限校验中间件
func CasbinRbacMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 先从jwt中拿出cliaim信息
		// 解析出user
		// 遍历 roles
		userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("[casibin]通过token解析到的userName去数据库中找User失败",
				zap.Error(err),
			)
			common.ReqBadFailWithMessage(fmt.Sprintf("[casibin]通过token解析到的userName去数据库中找User失败:%v", err.Error()), c)
			c.Abort()
			return
		}

		// 在这里支持服务账号
		if dbUser.AccountType == 2 {
			c.Next()
			return
		}

		// 获取当前请求的 path 和method
		path := c.Request.URL.Path
		method := c.Request.Method

		// 遍历roles 判断是否有权限
		pass := false
		for _, role := range dbUser.Roles {
			role := role
			//

			ok, err := models.CasbinCheckPermission(role.RoleValue, path, method)
			if err != nil {
				sc.Logger.Error("casbin校验出错",
					zap.Error(err),
					zap.String("userName", userName),
					zap.String("RoleValue", role.RoleValue),
					zap.String("path", path),
					zap.String("method", method),
				)
				common.ReqBadFailWithMessage(fmt.Sprintf("[casibin]casbin校验出错:%v", err.Error()), c)
				c.Abort()
			}
			sc.Logger.Debug("单个role-casbin校验结果",
				zap.String("userName", userName),
				zap.String("RoleValue", role.RoleValue),
				zap.String("path", path),
				zap.String("method", method),
				zap.Bool("pass", ok),
			)

			if ok {
				pass = true
				// 只要有一个权限就跳出循环
				break
			}

		}
		if !pass {
			sc.Logger.Error("casbin校验结果不通过",
				zap.Error(err),
				zap.String("userName", userName),
				zap.String("path", path),
				zap.String("method", method),
			)
			common.Req403WithWithMessage("[casibin]casbin校验结果不通过", c)
			c.Abort()
		}

		sc.Logger.Debug("casbin校验结果通过",
			zap.String("userName", userName),
			zap.String("path", path),
			zap.String("method", method),
		)

		c.Next()

	}
}
