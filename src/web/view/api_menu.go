package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getMenuList(c *gin.Context) {
	// 拿到用户对应的role列表 遍历role列表，找到menu list
	//  在拼接父子结构，返回的是数组，第一层是父级，第二层是子级
	useName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	dbUser, err := models.GetUserByUserName(useName)
	if err != nil {
		sc.Logger.Error("get user by username error", zap.Any("err", err))
		common.ReqBadFailWithMessage(fmt.Sprintf("get user by username error: %s", err), c)
	}
	roles := dbUser.Roles
	for _, role := range roles {
		role := role
		fmt.Println(role)
	}
}
