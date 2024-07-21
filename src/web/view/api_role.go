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

func getRoleListAll(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 从数据库中获取所有的菜单
	roles, err := models.GetRoleAll()
	if err != nil {
		sc.Logger.Error("获取角色列表错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("获取角色列表错误:%v", err.Error()), c)
		return
	}
	for _, role := range roles {
		role := role
		for _, menu := range role.Menus {
			menu.Key = menu.ID
			menu.Value = menu.ID
		}
		for _, api := range role.Apis {
			api.Key = api.ID
			api.Value = api.ID
		}
	}
	common.OkWithDetailed(roles, "获取角色列表成功", c)
}

func createRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	var reqRole models.Role
	err := c.ShouldBind(&reqRole)
	if err != nil {
		sc.Logger.Error("创建角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 字段校验，是否必填，范围是否正确
	err = validate.Struct(reqRole)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(
				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.ReqBadFailWithMessage(err.Error(), c)
		return

	}
	menus := make([]*models.Menu, 0)
	for _, menuId := range reqRole.MenuIds {
		dbmenu, err := models.GetMenuById(menuId)
		if err != nil {
			sc.Logger.Error("根据ID找菜单错误", zap.Any("菜单", menuId), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
		}
		menus = append(menus, dbmenu)
	}
	reqRole.Menus = menus
	err = reqRole.CreateOne()
	if err != nil {
		sc.Logger.Error("创建角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("创建角色成功", c)
}
