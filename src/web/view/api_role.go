package view

import (
	"bytes"
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"strconv"
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

// 过滤掉空字符串的函数
func filterEmptyStrings(ids []interface{}) ([]int, error) {
	var filteredIds []int
	for _, id := range ids {
		switch v := id.(type) {
		case string:
			if v != "" {
				intId, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				filteredIds = append(filteredIds, intId)
			}
		case float64: // JSON numbers are unmarshalled into float64
			filteredIds = append(filteredIds, int(v))
		}
	}
	return filteredIds, nil
}

// 创建角色
func createRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	// 读取并打印请求体内容
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		sc.Logger.Error("读取请求体失败", zap.Error(err))
		common.FailWithMessage("读取请求体失败", c)
		return
	}
	//sc.Logger.Info("收到的创建角色请求数据", zap.ByteString("请求体", bodyBytes))

	// 将读取的请求体重新放回c.Request.Body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqRole struct {
		Status    string        `json:"status"`
		RoleName  string        `json:"roleName"`
		RoleValue string        `json:"roleValue"`
		MenuIds   []interface{} `json:"menuIds"`
		ApiIds    []interface{} `json:"apiIds"`
	}

	err = c.ShouldBindJSON(&reqRole)
	if err != nil {
		sc.Logger.Error("解析新增角色请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 打印解析后的结构体
	//sc.Logger.Info("解析后的角色数据", zap.Any("角色", reqRole))

	// 去除menuIds和apiIds中的空字符串
	menuIds, err := filterEmptyStrings(reqRole.MenuIds)
	if err != nil {
		sc.Logger.Error("过滤menuIds失败", zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	apiIds, err := filterEmptyStrings(reqRole.ApiIds)
	if err != nil {
		sc.Logger.Error("过滤apiIds失败", zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 创建角色模型并填充数据
	role := models.Role{
		Status:    reqRole.Status,
		RoleName:  reqRole.RoleName,
		RoleValue: reqRole.RoleValue,
		MenuIds:   menuIds,
		ApiIds:    apiIds,
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(role)
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
	// 遍历角色menu 列表 找到角色
	for _, menuId := range role.MenuIds {
		dbMenu, err := models.GetMenuById(menuId)
		if err != nil {
			sc.Logger.Error("根据id找菜单错误", zap.Any("菜单", role), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		menus = append(menus, dbMenu)
	}
	role.Menus = menus

	// 创建角色
	err = role.CreateOne()
	if err != nil {
		sc.Logger.Error("创建角色错误", zap.Any("角色", role), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 获取完整的角色信息，包括关联的用户
	roleWithDetails, err := models.GetRoleById(role.ID)
	if err != nil {
		sc.Logger.Error("获取角色详情错误", zap.Uint("角色ID", uint(role.ID)), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 返回包含详细信息的响应
	common.OkWithDetailed(roleWithDetails, "创建成功", c)
}
