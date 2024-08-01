package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERROR  = 7
	SUCESS = 0
)

// 定义一个通用的返回结构体
type BaseResp struct {
	Code    int         `json:"code"` // 前后端交互的字段码
	Data    interface{} `json:"result"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
}

// http的200响应，但是不代表code是OK的
func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// Ok 是一个用于返回操作成功结果的 HTTP 响应的函数
func Ok(c *gin.Context) {
	Result(SUCESS, map[string]interface{}{}, "操作成功", c)
}

// OkWithMessage 是一个用于返回操作成功结果，并附带消息的 HTTP 响应的函数
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCESS, map[string]interface{}{}, message, c)
}

// OkWithData 是一个用于返回操作成功结果，并附带数据体的 HTTP 响应的函数
func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCESS, data, "查询成功", c)
}

// OkWithDetailed 是一个用于返回操作成功结果，并附带数据体和消息的 HTTP 响应的函数
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCESS, data, message, c)
}

// Fail 是一个用于返回操作失败的 HTTP 响应的函数
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

// FailWithMessage 是一个用于返回操作失败，并附带消息的 HTTP 响应的函数
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

// FailWithData 是一个用于返回操作失败，并附带数据体的 HTTP 响应的函数
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

// 参数错误
func Result400(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusBadRequest, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

func Result403(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusForbidden, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// 没权限
func Result401(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusUnauthorized, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// 服务器错误
func Result5XX(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusInternalServerError, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// ReqBadFailWithMessage 函数向客户端返回HTTP状态码为400的响应，表示请求错误
func ReqBadFailWithMessage(message string, c *gin.Context) {
	Result400(ERROR, map[string]interface{}{}, message, c)
}

// ReqBadFailWithWithDetailed 函数用于处理HTTP请求失败的情况，并返回详细的错误信息
func ReqBadFailWithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result400(ERROR, data, message, c)
}

// Req401WithWithDetailed 是一个用于返回HTTP 401 Unauthorized响应的函数

func Req401WithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(ERROR, data, message, c)
}

// Req403WithWithMessage 响应HTTP 403 Forbidden错误，并返回自定义的错误信息
func Req403WithWithMessage(message string, c *gin.Context) {
	Result403(ERROR, map[string]interface{}{}, message, c)
}

// Req5XXWithWithDetailed 是一个处理HTTP 5XX系列错误的函数
func Req5XXWithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(ERROR, data, message, c)
}
