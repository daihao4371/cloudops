package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
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
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServeConfig)
	// 校验字段是否必填，范围是否正确
	err = validate.Struct(user)
	if err != nil {

		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithDetailed(
				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"Request error",
				c,
			)
			return
		}

		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(user, c)

	dbUser, err := models.CheckUserPassword(&user)
	if err != nil {
		sc.Logger.Error("用户登录失败，用户名不存在或密码错误", zap.Error(err))
		common.ReqBadFailWithMessage("用户名不存在或密码错误", c)
		return
	}
	// 生成token，并返回给前端
	models.TokenNext(dbUser, c)
}

// 登录以后获取用户信息
func getUserInfoAfterLoign(c *gin.Context) {
	jwtClaim := c.MustGet(common.GIN_CTX_JWT_CLAIM).(*models.UserCustomClaims)
	common.OkWithDetailed(jwtClaim.User, "ok", c)
}
