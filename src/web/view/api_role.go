package view

import (
	"bytes"
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"strconv"
)

// 函数用于获取所有角色列表
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

// 更新角色状态
func setRoleStatus(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	var reqRole struct {
		Id     uint   `json:"id" validate:"required"`     // Id 改回 uint 类型
		Status string `json:"status" validate:"required"` // 根据你的需要添加其他字段
	}

	err := c.ShouldBindJSON(&reqRole)
	if err != nil {
		sc.Logger.Error("解析请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 做字段校验，是否必填，范围是否正确
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

	dbRole, err := models.GetRoleById(reqRole.Id)
	if err != nil {
		sc.Logger.Error("根据id找角色错误", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	dbRole.Status = reqRole.Status
	err = dbRole.UpdateMenus(dbRole.Menus)
	if err != nil {
		sc.Logger.Error("更新角色和关联的菜单错误", zap.Any("角色", dbRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)
	//sc.Logger.Info("更新角色状态成功", zap.Any("角色", dbRole))
}

// 删除角色
func deleteRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 获取路径参数中的id
	id := c.Param("id")
	sc.Logger.Info("删除角色", zap.String("id", id))

	// 将id从字符串转换为整数
	intVar, err := strconv.Atoi(id)
	if err != nil {
		sc.Logger.Error("ID 转换错误", zap.String("id", id), zap.Error(err))
		common.FailWithMessage("ID 转换错误", c)
		return
	}

	// 根据id查找角色
	dbRole, err := models.GetRoleById(uint(intVar))
	if err != nil {
		sc.Logger.Error("根据ID找角色错误", zap.Uint("id", uint(intVar)), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 删除关联表 role_menus 中的记录
	err = models.DB.Model(&dbRole).Association("Menus").Clear()
	if err != nil {
		sc.Logger.Error("删除关联表记录错误", zap.Uint("id", uint(intVar)), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 删除角色，明确设置模型
	//err = models.DB.Model(&models.Role{}).Delete(dbRole, clause.Associations).Error
	err = models.DB.Delete(&dbRole).Error
	if err != nil {
		sc.Logger.Error("根据ID删除角色错误", zap.Uint("id", uint(intVar)), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("删除成功", c)
	sc.Logger.Info("删除角色成功", zap.Uint("id", uint(intVar)))
}

// 更新角色信息
func updateRole(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)

	var reqRole models.Role

	// 读取并打印原始请求体
	body, err := io.ReadAll(c.Request.Body)
	/*		if err == nil {
			sc.Logger.Info("原始请求体", zap.String("body", string(body)))
		}*/
	// 重置请求体以便后续操作
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 手动绑定JSON，先绑定非嵌入字段
	err = json.Unmarshal(body, &reqRole)
	if err != nil {
		sc.Logger.Error("解析更新角色请求失败", zap.Any("角色", reqRole), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 手动绑定Model中的ID字段
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err == nil {
		if id, ok := data["id"].(float64); ok {
			reqRole.ID = uint(id)
		}
	}

	// 打印绑定后的角色信息
	sc.Logger.Info("绑定后的角色信息", zap.Any("角色", reqRole))

	// 检查ID是否有效
	if reqRole.ID == 0 {
		sc.Logger.Error("角色ID无效", zap.Uint("id", reqRole.ID))
		common.FailWithMessage("角色ID无效", c)
		return
	}

	role, err := models.GetRoleById(reqRole.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sc.Logger.Error("角色不存在", zap.Uint("id", reqRole.ID))
			common.FailWithMessage("角色不存在", c)
		} else {
			sc.Logger.Error("根据ID找角色错误", zap.Uint("id", reqRole.ID), zap.Error(err))
			common.FailWithMessage("数据库错误", c)
		}
		return
	}

	// 字段校验
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
	// 遍历角色menu列表，找到菜单
	for _, menuId := range reqRole.MenuIds {
		dbMenu, err := models.GetMenuById(menuId)
		if err != nil {
			sc.Logger.Error("根据ID找菜单错误", zap.Any("menu", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		menus = append(menus, dbMenu)
	}

	apis := make([]*models.Api, 0)
	// 遍历角色apis列表，找到Apisid
	for _, apiId := range reqRole.ApiIds {
		dbApi, err := models.GetApiById(apiId)
		if err != nil {
			sc.Logger.Error("根据id找api错误", zap.Any("api", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		apis = append(apis, dbApi)
	}

	if len(menus) > 0 {
		err = role.UpdateMenus(menus)
		if err != nil {
			msg := fmt.Sprintf("更新角色和关联的菜单错误 err:%v", err.Error())
			sc.Logger.Error(msg, zap.Any("角色", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
	}

	if len(apis) > 0 {
		err = role.UpdateApis(apis, sc)
		if err != nil {
			msg := fmt.Sprintf("更新角色和关联的api错误 err:%v", err.Error())
			sc.Logger.Error(msg, zap.Any("角色", reqRole), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
	}

	common.OkWithMessage("更新成功", c)
	//sc.Logger.Info("更新角色成功", zap.Any("角色", reqRole))
}
