package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一返回结构体
type Result struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Result{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败返回
func Fail(c *gin.Context, code int, message string) {
	c.JSON(code, Result{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// FailWithData 失败返回带数据
func FailWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Result{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
