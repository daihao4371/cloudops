package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	c.JSON(200, BaseResp{
		Code:    200,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

const (
	ERROR  = 7
	SUCESS = 0
)

func Ok(c *gin.Context) {
	Result(SUCESS, map[string]interface{}{}, "operate sucess", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCESS, data, "Search successful", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "operate failed", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)

}

// 参数错误
func Result400(code int, dat interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, BaseResp{
		Code:    code,
		Data:    dat,
		Message: msg,
		Type:    "",
	})
}

// 权限错误
func Result401(code int, dat interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, BaseResp{
		Code:    code,
		Data:    dat,
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

func Result5XX(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusForbidden, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

func ReqBadFailWithMessage(message string, c *gin.Context) {
	Result400(ERROR, map[string]interface{}{}, message, c)
}
func ReqBadFailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result400(ERROR, data, message, c)
}
func ReqBadFail(c *gin.Context) {
	Result400(ERROR, map[string]interface{}{}, "operate failed", c)
}

func Req401WithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(ERROR, data, message, c)
}

func Req403WithWithMessage(message string, c *gin.Context) {
	Result403(ERROR, map[string]interface{}{}, message, c)
}

func Req5XXWithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(ERROR, data, message, c)
}
