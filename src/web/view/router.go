package view

import (
	"github.com/gin-gonic/gin"
	"time"
)

// 配置路由
func ConfigRouter(r *gin.Engine) {
	// 定义一个组
	// 共享中间件
	base := r.Group("/basic-api")
	{
		base.GET("/ping", ping)
		base.GET("/now", getNowTs)
		//base.GET("/long", longRequest)

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
