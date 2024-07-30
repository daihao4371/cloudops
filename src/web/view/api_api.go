package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

// 获取所有的api接口
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
