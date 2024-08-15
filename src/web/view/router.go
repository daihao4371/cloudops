package view

import (
	"cloudops/src/web/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

// 配置路由
func ConfigRouter(r *gin.Engine) {
	// 定义一个组
	// 共享中间件
	//base := r.Group("/basic-api")
	base := r.Group("/")
	{
		base.GET("/ping", ping)
		base.GET("/now", getNowTs)
		//base.GET("/long", longRequest)
		base.POST("/login", UserLogin)

	}
	// 登录后才能访问的路由
	afterLoginApiGroup := r.Group("/api")
	afterLoginApiGroup.Use(middleware.JWTAuthMiddleware()).Use(middleware.CasbinRbacMiddleware())
	{
		afterLoginApiGroup.GET("/getUserInfo", getUserInfoAfterLogin)
		afterLoginApiGroup.GET("/getPermCode", getPermCode)
	}
	// 底座模块
	systemApiGroup := afterLoginApiGroup.Group("/system")
	{
		// 菜单相关
		systemApiGroup.GET("/getMenuList", getMenuList)
		systemApiGroup.GET("/getMenuListAll", getMenuListAll)
		systemApiGroup.POST("/updateMenu", updateMenu)
		systemApiGroup.POST("/createMenu", createMenu)
		systemApiGroup.DELETE("/deleteMenu/:id", deleteMenu)

		// 角色相关
		systemApiGroup.GET("/getRoleListAll", getRoleListAll)
		systemApiGroup.POST("/createRole", createRole)
		systemApiGroup.POST("/updateRole", updateRole)
		systemApiGroup.POST("/setRoleStatus", setRoleStatus)
		systemApiGroup.DELETE("/deleteRole/:id", deleteRole)

		// 用户相关
		systemApiGroup.POST("/createAccount", createAccount)
		systemApiGroup.POST("/accountExist", accountExist)
		systemApiGroup.POST("/updateAccount", updateAccount)
		systemApiGroup.POST("/changePassword", changePassword)
		systemApiGroup.GET("/getAccountList", getAccountList)
		systemApiGroup.GET("/getAllUserAndRoles", getAllUserAndRoles)
		systemApiGroup.DELETE("/deleteAccount/:id", deleteAccount)

		// API相关
		systemApiGroup.GET("/getApiList", getApiList)
		systemApiGroup.GET("/getApiListAll", getApiListAll)
		systemApiGroup.POST("/createApi", createApi)
		systemApiGroup.POST("/updateApi", updateApi)
		systemApiGroup.DELETE("/deleteApi/:id", deleteApi)
	}

	// 服务树模块
	streeApiGroup := afterLoginApiGroup.Group("/stree")
	{
		streeApiGroup.GET("/getStreeNodeList", getStreeNodeList)
		//streeApiGroup.GET("/getTopStreeNodes", getTopStreeNodes)
		streeApiGroup.GET("/getTopStreeNodes", getTopStreeNodesUseCache)
		streeApiGroup.POST("/createStreeNode", createStreeNode)
		streeApiGroup.DELETE("/deleteStreeNode/:id", deleteStreeNode)
		streeApiGroup.GET("/getChildrenStreeNodes/:pid", getChildrenStreeNodes)
		streeApiGroup.POST("/updateStreeNode", updateStreeNode)

	}
}

// 测试路由
func getNowTs(c *gin.Context) {
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}

/*func longRequest(c *gin.Context) {
	fmt.Println("longRequest请求开始，休息5秒")
	time.Sleep(time.Second * 5)
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
*/
