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
	afterLoginApiGroup.Use(middleware.JWTAuthMiddleware())
	{
		afterLoginApiGroup.GET("/getUserInfo", getUserInfoAfterLogin)
		afterLoginApiGroup.GET("/getPermCode", getPermCode)
	}
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

		// 用户相关
		systemApiGroup.POST("/createAccount", createAccount)
		systemApiGroup.POST("/accountExist", accountExist)
		systemApiGroup.POST("/updateAccount", updateAccount)
		systemApiGroup.POST("/changePassword", changePassword)
		systemApiGroup.GET("/getAccountList", getAccountList)
		systemApiGroup.GET("/getAllUserAndRoles", getAllUserAndRoles)
		systemApiGroup.DELETE("/deleteAccount/:id", deleteAccount)

		// 权限相关
		systemApiGroup.GET("/getApiList", getApiList)
	}
}

func getNowTs(c *gin.Context) {
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}

/*func longRequest(c *gin.Context) {
	fmt.Println("longRequest请求开始，休息5秒")
	time.Sleep(time.Second * 5)
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
*/
