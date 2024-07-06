package middleware

import "github.com/gin-gonic/gin"

// ConfigMiddleware 返回一个 gin.HandlerFunc 类型的中间件函数，该函数用于将配置映射中的配置项设置到 gin.Context 中
// 参数 m 是一个包含配置项键值对的 map[string]interface{} 类型的映射
// 返回值是一个 gin.HandlerFunc 类型的中间件函数
func ConfigMiddleware(m map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 遍历配置映射
		for k, v := range m {
			// 将配置项设置到 gin.Context 中
			c.Set(k, v)
		}
		// 继续执行后续的中间件或处理函数
		c.Next()
	}
}
