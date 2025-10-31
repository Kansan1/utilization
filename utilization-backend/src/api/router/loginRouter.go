package router

import (
	"net/http"
	"utilization-backend/src/api/models"
	"utilization-backend/src/api/result"
	"utilization-backend/src/api/v1/login"

	"github.com/gin-gonic/gin"
)

// LoginRouter 注册登录相关路由
func LoginRouter(r *gin.Engine) *gin.RouterGroup {
	apiUser := r.Group("/api")
	{
		user := apiUser.Group("/user")
		user.POST("/login", handleLogin)

	}
	return apiUser
}

// handleLogin 处理登录请求
func handleLogin(ctx *gin.Context) {

	var user models.LoginUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	login.HandleLogin(ctx, user)
}
