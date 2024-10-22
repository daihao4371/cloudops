package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"sort"
	"strconv"
)

/*
getMenuList 函数用于获取当前用户的菜单列表
参数：

	c *gin.Context：Gin框架的上下文对象

返回值：

	无，直接通过HTTP响应返回结果
*/
func getMenuList(c *gin.Context) {
	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("通过token解析到的userName去数据库中找User失败",
			zap.Error(err),
		)
		common.ReqBadFailWithMessage(fmt.Sprintf("通过token解析到的userName去数据库中找User失败:%v", err.Error()), c)
		return
	}

	// 	遍历 role列表 找到 Menu list
	// 在拼装父子结构 返回的是数组 第一层 father 第2层children

	fatherMenuMap := make(map[uint]*models.Menu)
	uniqueChildMap := make(map[uint]*models.Menu)
	roles := dbUser.Roles
	for _, role := range roles {
		role := role
		//sc.Logger.Info("遍历user的role打印", zap.String("role", role.RoleName),
		//	zap.Any("role的menulist详情", role.Menus),
		//)
		// 这里做判断 如果role的状态是 禁用的，那么 就跳过这个角色的菜单
		if role.Status == common.COMMON_STATUS_DISABLE {
			sc.Logger.Info("用户的角色禁用，跳过筛选菜单",
				zap.String("用户", dbUser.Username),
				zap.String("角色", role.RoleName),
			)
			continue
		}

		for _, menu := range role.Menus {
			menu := menu

			// 如果菜单是禁用的那么 跳过这个菜单
			if menu.Status == common.COMMON_STATUS_DISABLE {
				// 超管角色下面的用户能看到所有菜单
				if role.RoleValue != "super" {
					sc.Logger.Info("菜单禁用，跳过筛选菜单",
						zap.String("用户", dbUser.Username),
						zap.String("角色", role.RoleName),
						zap.String("菜单", menu.Name),
					)
					continue
				}

			}

			// 拼接前端依赖的字段
			menu.Meta = &models.MenuMeta{}
			menu.Meta.Icon = menu.Icon
			menu.Meta.Title = menu.Title
			menu.Key = menu.ID
			menu.Value = menu.ID
			menu.Meta.ShowMenu = common.COMMON_SHOW_MAP[menu.Show]
			//menu.Meta.ShowMenu = true
			menu.Meta.HideMenu = !common.COMMON_SHOW_MAP[menu.Show]

			//if menu.Path == "stree" {
			//	menu.Meta.IgnoreKeepAlive = true
			//}
			menu.Meta.IgnoreKeepAlive = true
			// 拼接小id 给前端的树形结构的

			if menu.Pid == 0 {
				// 说明这个菜单是父级
				fatherMenuMap[menu.ID] = menu
				continue
			}

			// 说明menu是子集
			fatherMenu, err := models.GetMenuById(menu.Pid)
			if err != nil {
				sc.Logger.Error("通过Pid找menu错误", zap.Error(err))
				continue
			}
			fatherMenu.Meta = &models.MenuMeta{}
			fatherMenu.Meta.Icon = fatherMenu.Icon
			fatherMenu.Meta.Title = fatherMenu.Title
			fatherMenu.Key = fatherMenu.ID
			fatherMenu.Value = fatherMenu.ID
			fatherMenu.Meta.ShowMenu = common.COMMON_SHOW_MAP[fatherMenu.Show]
			fatherMenu.Meta.HideMenu = !common.COMMON_SHOW_MAP[fatherMenu.Show]

			// 判断子菜单是否是重复的
			_, ok := uniqueChildMap[menu.ID]
			if ok {
				continue
			}
			// 否则的话先塞入
			uniqueChildMap[menu.ID] = menu

			load, ok := fatherMenuMap[fatherMenu.ID]

			if !ok {
				//之前还没设置过 这个父级
				fatherMenu.Children = make([]*models.Menu, 0)
				fatherMenu.Children = append(fatherMenu.Children, menu)
				fatherMenuMap[fatherMenu.ID] = fatherMenu
			} else {
				// 存在的话 我们直接把menu塞入 Children
				load.Children = append(load.Children, menu)
			}

		}

	}

	finalMenus := make([]*models.Menu, 0)
	finalMenuIds := []int{}

	for id := range fatherMenuMap {
		//mId:=mId
		finalMenuIds = append(finalMenuIds, int(id))
	}
	sort.Ints(finalMenuIds)

	// 最终遍历	 fatherMenuMap
	for _, id := range finalMenuIds {
		//mId:=mId
		m := fatherMenuMap[uint(id)]
		finalMenus = append(finalMenus, m)
	}
	common.OkWithDetailed(finalMenus, "ok", c)
}

/*
getMenuListAll 函数从数据库中获取所有的菜单列表，并返回给前端
参数：
c *gin.Context：gin框架的上下文对象
返回值：
无，通过HTTP响应返回结果
*/
func getMenuListAll(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 从数据库中获取所有的菜单
	menus, err := models.GetMenuAll()
	if err != nil {
		sc.Logger.Error("获取菜单列表错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("获取菜单列表错误:%v", err.Error()), c)
		return
	}
	for _, menu := range menus {
		menu := menu
		// 拼接前端依赖的字段
		menu.Meta = &models.MenuMeta{}
		menu.Meta.Icon = menu.Icon
		menu.Meta.Title = menu.Title
		menu.Meta.ShowMenu = common.COMMON_SHOW_MAP[menu.Show]
		menu.Key = menu.ID
		menu.Value = menu.ID

	}
	common.OkWithDetailed(menus, "ok", c)
}

// 更新菜单
func updateMenu(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	var reqMenu models.Menu

	err := c.ShouldBindJSON(&reqMenu)
	if err != nil {
		sc.Logger.Error("解析更新菜单请求失败", zap.Any("菜单", reqMenu), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqMenu)
	if err != nil {

		// 这里为什么要判断错误是否是 ValidationErrors
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

	_, err = models.GetMenuById(int(reqMenu.ID))
	if err != nil {
		sc.Logger.Error("根据id找menu错误", zap.Any("菜单", reqMenu), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	err = reqMenu.UpdateOne()
	if err != nil {
		sc.Logger.Error("根据id更新menu错误", zap.Any("菜单", reqMenu), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)

}

// 创建菜单
func createMenu(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu 字段
	var reqMenu models.Menu
	err := c.ShouldBindJSON(&reqMenu)
	if err != nil {
		sc.Logger.Error("解析新增菜单请求失败", zap.Any("菜单", reqMenu), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqMenu)
	if err != nil {
		// 这里为什么要判断错误是否是 ValidationErrors
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

	err = reqMenu.CreateOne()
	if err != nil {
		sc.Logger.Error("创建menu错误", zap.Any("菜单", reqMenu), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("创建成功", c)
}

// 删除菜单
func deleteMenu(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	id := c.Param("id")
	sc.Logger.Info("删除菜单", zap.Any("id", id))

	// 先 去db中根据id找到这个user
	intVar, _ := strconv.Atoi(id)
	dbMenu, err := models.GetMenuById(intVar)
	if err != nil {
		sc.Logger.Error("根据id找菜单错误", zap.Any("菜单", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	err = dbMenu.DeleteOne()
	if err != nil {
		sc.Logger.Error("根据id删除菜单错误", zap.Any("菜单", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("删除成功", c)
}
