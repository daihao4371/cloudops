package view

import (
	"cloudops/src/common"
	"cloudops/src/models"
	"github.com/gin-gonic/gin"
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
	// 校验字段是否必填，范围是否正确
	err = validate.Struct(user)
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithData(user, c)
}
