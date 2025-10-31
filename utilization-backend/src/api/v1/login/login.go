package login

import (
	"net/http"
	"utilization-backend/src/api/dao"
	"utilization-backend/src/api/models"
	"utilization-backend/src/api/result"
	"utilization-backend/src/api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleLogin(c *gin.Context, user models.LoginUser) {
	// 验证用户登录信息
	isValid, err := dao.ValidateUser(&user)
	if err != nil {
		// 记录错误日志
		zap.L().Error("数据库查询失败", zap.Error(err))
		result.Fail(c, http.StatusInternalServerError, "服务器内部错误")
		return
	}

	if isValid {
		// 生成JWT token
		token, err := utils.GenerateToken(user.Username)
		if err != nil {
			zap.L().Error("生成token失败", zap.Error(err))
			result.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		result.Success(c, gin.H{
			"token": token,
			"user": gin.H{
				"username": user.Username,
			},
		})
		return
	}

	result.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
}
