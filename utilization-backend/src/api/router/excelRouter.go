package router

import (
	"utilization-backend/src/api/v1/excel"

	"github.com/gin-gonic/gin"
)

func ExcelRouter(r *gin.Engine) *gin.RouterGroup {
	apiExcel := r.Group("/api")
	{
		excel := apiExcel.Group("/excel")
		excel.GET("/download", handleDownloadExcel)
	}
	return apiExcel
}

func handleDownloadExcel(ctx *gin.Context) {
	excel.DownLoadExcel(ctx)
}
