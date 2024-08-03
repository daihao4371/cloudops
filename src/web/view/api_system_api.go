package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// 获取所有的api接口
func getApiList(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	apis, err := models.GetApiAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的api接口错误",
			zap.Error(err),
		)
		common.ReqBadFailWithMessage(fmt.Sprintf("去数据库中拿所有的api接口错误:%v", err.Error()), c)
		return
	}

	fatherApiMap := make(map[uint]*models.Api)
	for _, api := range apis {
		api := api
		// 这里做判断 如果role的状态是 禁用的，那么 就跳过这个角色的菜单
		api.Key = api.ID
		api.Value = api.ID
		if api.Pid == 0 {
			// 说明这个菜单是父级
			fatherApiMap[api.ID] = api
			continue
		}

		// 说明menu是子集
		fatherApi, err := models.GetApiById(api.Pid)
		if err != nil {
			sc.Logger.Error("通过Pid找Api错误", zap.Error(err))
			continue
		}

		fatherApi.Key = fatherApi.ID
		fatherApi.Value = fatherApi.ID
		load, ok := fatherApiMap[fatherApi.ID]

		if !ok {
			//之前还没设置过 这个父级
			fatherApi.Children = make([]*models.Api, 0)
			fatherApi.Children = append(fatherApi.Children, api)
			fatherApiMap[fatherApi.ID] = fatherApi
		} else {
			// 存在的话 我们直接把menu塞入 Children
			load.Children = append(load.Children, api)
		}
	}

	finalApis := make([]*models.Api, 0)
	// 最终遍历	 fatherMenuMap
	for _, m := range fatherApiMap {
		m := m
		finalApis = append(finalApis, m)
	}
	common.OkWithDetailed(finalApis, "ok", c)
}

// 获取所有的API接口列表
func getApiListAll(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	apis, err := models.GetApiAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的api接口错误",
			zap.Error(err),
		)
		common.ReqBadFailWithMessage(fmt.Sprintf("去数据库中拿所有的api接口错误:%v", err.Error()), c)
		return
	}
	for _, api := range apis {
		api := api
		// 这里做判断 如果role的状态是 禁用的，那么 就跳过这个角色的菜单
		api.Key = api.ID
		api.Value = api.ID
	}
	common.OkWithDetailed(apis, "ok", c)
}

// createApi 创建api接口
func createApi(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	var reqApi models.Api
	err := c.ShouldBindJSON(&reqApi)
	if err != nil {
		sc.Logger.Error("解析新增reqApi请求失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("解析参数错误:%v", err.Error()), c)
		return
	}
	// 字段校验，是否必填，范围是否正确
	err = validate.Struct(reqApi)
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
	err = reqApi.CreateOne()
	if err != nil {
		sc.Logger.Error("新增api接口失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf(err.Error()), c)
		return
	}
	common.OkWithDetailed(reqApi, "新增成功", c)
}

// 更新api接口
func updateApi(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	var reqApi models.Api
	err := c.ShouldBindJSON(&reqApi)
	if err != nil {
		sc.Logger.Error("解析更新reqApi请求失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.FailWithMessage(fmt.Sprintf("解析参数错误:%v", err.Error()), c)
		return
	}
	// 字段校验，是否必填，范围是否正确
	err = validate.Struct(reqApi)
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
	// 通过ID获取API对象
	existingApi, err := models.GetApiById(int(reqApi.ID))
	if err != nil {
		sc.Logger.Error("根据id找api接口失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.ReqBadFailWithMessage(err.Error(), c)
		return
	}
	// 更新现有API对象的字段，避免更新CreateAt
	existingApi.Path = reqApi.Path
	existingApi.Method = reqApi.Method
	existingApi.Pid = reqApi.Pid
	existingApi.Title = reqApi.Title
	existingApi.Type = reqApi.Type
	existingApi.UpdatedAt = time.Now()

	err = existingApi.UpdateOne()
	if err != nil {
		sc.Logger.Error("更新api接口失败", zap.Any("reqApi", reqApi), zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf(err.Error()), c)
		return
	}
	common.OkWithDetailed(existingApi, "更新成功", c)
}

// deleteApi 删除api接口
func deleteApi(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	id := c.Param("id")
	sc.Logger.Info("删除api接口", zap.String("id", id))

	//先去数据库中拿这个api接口
	// 将字符串ID转换为整型
	intVar, _ := strconv.Atoi(id)
	dbApi, err := models.GetApiById(intVar)
	if err != nil {
		sc.Logger.Error("根据id找api接口失败", zap.Any("api", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 执行删除操作
	err = dbApi.DeleteOne()
	if err != nil {
		sc.Logger.Error("根据ID删除api接口失败", zap.Any("api", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
	}
	common.OkWithDetailed(dbApi, "删除成功", c)
}
