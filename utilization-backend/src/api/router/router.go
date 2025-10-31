package router

import (
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.Logger())
	//初始化日志
	config.InitLogger()

	// config.Logger.Info("信息日志", zap.String("key", "value"))
	// config.Logger.Error("错误日志", zap.Error(err))
	// config.Logger.Debug("调试日志")
	// config.Logger.Warn("警告日志")

	err := config.InitDB()
	if err != nil {
		config.Logger.Fatal("数据库连接失败",
			zap.Error(err))
	}

	LoginRouter(r)
	HomeRouter(r)
	ExcelRouter(r)

}
