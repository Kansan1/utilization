package utils

import (
	"net/http"
	"time"
	"utilization-backend/src/api/dao"
	"utilization-backend/src/api/result"
)

// RunMonthlyTask 运行每月任务
func RunMonthlyTask() {
	// 获取当前时间
	now := time.Now()

	// 计算下个月1号的时间
	nextMonth := now.AddDate(0, 1, 0)
	nextMonthFirstDay := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	// 计算距离下个月1号的时间间隔
	duration := nextMonthFirstDay.Sub(now)

	// 创建定时器
	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		// 执行清空操作
		err := dao.DeleteEquipmentRepairCompletionRateEachMonth()
		if err != nil {
			// 记录错误日志
			result.Fail(nil, http.StatusInternalServerError, "自动清空设备维护完成率数据失败")
		}

		// 重新设置定时器，继续执行下个月的任务
		RunMonthlyTask()
	}()
}
