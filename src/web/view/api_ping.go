package view

import "github.com/gin-gonic/gin"

// 测试接口
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
