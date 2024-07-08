package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func UserLogin(c *gin.Context) {
	// 校验用户名密码
	var user models.UserLoginRequest
	err := c.ShouldBindJSON(&user)

	// 判断JSON解析是否正确
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验字段是否必填，范围是否正确
	err = validate.Struct(user)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(
				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"Request error",
				c,
			)
			return
		}

		common.ReqBadFailWithMessage(err.Error(), c)
		return
	}

	dbUser, err := models.CheckUserPassword(&user)
	if err != nil {
		sc.Logger.Error("登录失败，用户名不存在或密码错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("用户名不存在或密码错误:%v", err.Error()), c)
		return
	}
	// 生成token，并返回给前端
	models.TokenNext(dbUser, c)
}

// 登录以后获取用户信息
func getUserInfoAfterLogin(c *gin.Context) {

	// 我得拿到 userCliams
	jwtClaims, ok := c.MustGet(common.GIN_CTX_JWT_CLAIM).(*models.UserCustomClaims)
	if !ok {
		// 如果获取失败,则返回错误
		common.FailWithMessage("获取用户信息失败", c)
		return
	}
	common.OkWithDetailed(jwtClaims.User, "ok", c)

	/*	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("通过token解析到的userName去数据库中找User失败",
				zap.Error(err),
			)
			common.ReqBadFailWithMessage(fmt.Sprintf("通过token解析到的userName去数据库中找User失败:%v", err.Error()), c)
			return
		}
		common.OkWithDetailed(dbUser, "ok", c)*/
}
