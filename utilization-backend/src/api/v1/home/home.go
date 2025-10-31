package home

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"utilization-backend/src/api/dao"
	"utilization-backend/src/api/models"
	"utilization-backend/src/api/result"
	"utilization-backend/src/api/vo"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/**
 * 获取月度稼动率列表
 * @param currentTime 当前时间
 * @return 当月稼动率列表
 */
func GetUtilizationListByMonth(currentTime string, ctx *gin.Context) {
	fmt.Printf("currentTime: %v\n", currentTime)
	var utilizationVoList []vo.UtilizationVo
	year, month := CurrentTimeToYearMonth(currentTime)
	// 从数据库获取数据
	utilizationWorkingHourList, utilizationExceptionList :=
		dao.GetUtilizationListByMonth(year, month)

	//处理数据
	fmt.Println("handleUtilizationList 开始执行")

	utilizationVoList = handleUtilizationList(
		utilizationWorkingHourList,
		utilizationExceptionList,
		currentTime)

	// 返回数据
	result.Success(ctx, utilizationVoList)
}

func GetAllUtilizationList(ctx *gin.Context) {
	currentYear := time.Now().Year()
	fmt.Println("currentYear", currentYear)
	var utilizationAllVoList []vo.UtilizationAllVo
	for i := 1; i <= 12; i++ {
		currentTime := fmt.Sprintf("%d-%02d", currentYear, i)
		fmt.Printf("currentTime: %v\n", currentTime)
		year, month := CurrentTimeToYearMonth(currentTime)
		// 从数据库获取数据
		utilizationWorkingHourList, utilizationExceptionList :=
			dao.GetUtilizationListByMonth(year, month)
		//处理数
		valueAll, valueActual, valueRate := HandleAllUtilizationList(utilizationWorkingHourList, utilizationExceptionList)
		utilizationAllVoList = append(utilizationAllVoList, vo.UtilizationAllVo{
			Date:        currentTime,
			Value:       valueRate,
			ValueAll:    valueAll,
			ValueActual: valueActual,
		})
	}
	fmt.Println("utilizationVoList", utilizationAllVoList)
	result.Success(ctx, utilizationAllVoList)
}

func HandleAllUtilizationList(utilizationWorkingHourList []models.UtilizationWorkingHour, utilizationExceptionList []models.UtilizationException) (float32, float32, float32) {
	var utilizationList0 []models.Utilization
	var utilizationList1 []models.Utilization
	var value0 float32
	var value1 float32
	for _, utilizationWorkingHour := range utilizationWorkingHourList {
		// 正常工作时间
		if utilizationWorkingHour.Type == "0" {
			utilizationList0 = append(utilizationList0, models.Utilization{
				Date:  utilizationWorkingHour.Time,
				Value: StringToFloat32(utilizationWorkingHour.WorkingHours),
				Line:  utilizationWorkingHour.Line,
			})
		} else {
			var AbnormalStopTime float32
			for _, exception := range utilizationExceptionList {
				//算出停线时间
				if exception.Time == utilizationWorkingHour.Time {
					AbnormalStopTime = CalculateTimeInterval(exception.AbnormalStopTime)
					fmt.Printf("异常停线时间%v:%v      %v\n", exception.Time, AbnormalStopTime, exception.AbnormalStopTime)
				}
			}
			utilizationList1 = append(utilizationList1, models.Utilization{
				Date:  utilizationWorkingHour.Time,
				Value: StringToFloat32(utilizationWorkingHour.WorkingHours) - AbnormalStopTime,
				Line:  utilizationWorkingHour.Line,
			})
		}

	}
	for _, utilization := range utilizationList0 {
		value0 += utilization.Value
	}
	fmt.Println("utilizationList0", utilizationList0)
	fmt.Println("utilizationList1", utilizationList1)
	for _, utilization := range utilizationList1 {
		value1 += utilization.Value
	}
	fmt.Println("value1", value1)
	fmt.Println("value0", value0)
	if value0 == 0 {
		return 0, 0, 0
	} else {
		return value0, value1, value1 / value0 / 3
	}
}

func DeleteEquipmentRepairCompletionRate(names []string, ctx *gin.Context) {
	err := dao.DeleteEquipmentRepairCompletionRate(names)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, "删除设备维护完成率数据成功")
}

func GetEquipmentRepairCompletionRateList(ctx *gin.Context) {

	yearAndMonth := time.Now().Format("2006-01")
	equipmentRepairCompletionRateList, err := dao.GetEquipmentRepairCompletionRateList(yearAndMonth)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
	}
	result.Success(ctx, equipmentRepairCompletionRateList)
}

func AddOrUpdateEquipmentRepairCompletionRate(rates []models.EquipmentRepairCompletionRate, ctx *gin.Context) {
	yearAndMonth := time.Now().Format("2006-01")
	equipmentRepairCompletionRateList, err := dao.GetEquipmentRepairCompletionRateList(yearAndMonth)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var mapList = make(map[string]vo.EquipmentRepairCompletionRateVo)

	for _, rate := range equipmentRepairCompletionRateList {
		mapList[rate.Name] = vo.EquipmentRepairCompletionRateVo{
			Name:       rate.Name,
			Value:      rate.Value,
			CreateTime: rate.CreateTime,
		}
	}

	today := time.Now().Format("2006-01")
	for _, rate := range rates {
		value, exists := mapList[rate.Name]
		if exists {
			// 检查是否是同一天的数据
			recordDate := strings.Split(strings.Split(value.CreateTime, " ")[0], "-")[0] + "-" + strings.Split(strings.Split(value.CreateTime, " ")[0], "-")[1]
			if recordDate == today {
				// 同一天则更新
				var updateRate = models.EquipmentRepairCompletionRate{
					Name:  rate.Name,
					Value: rate.Value,
				}
				updateTime := time.Now().Format("2006-01-02 15:04:05")
				err := dao.UpdateEquipmentRepairCompletionRate(updateRate, updateTime)
				if err != nil {
					result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("更新设备 %s 的维护完成率失败: %v", rate.Name, err))
					return
				}
			} else {
				// 不是同一天则新增
				createTime := time.Now().Format("2006-01-02 15:04:05")
				err := dao.AddEquipmentRepairCompletionRate(rate, createTime)
				if err != nil {
					result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("添加设备 %s 的维护完成率失败: %v", rate.Name, err))
					return
				}
			}
		} else {
			// 设备不存在则新增
			createTime := time.Now().Format("2006-01-02 15:04:05")
			err := dao.AddEquipmentRepairCompletionRate(rate, createTime)
			if err != nil {
				result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("添加设备 %s 的维护完成率失败: %v", rate.Name, err))
				return
			}
		}
	}

	result.Success(ctx, "设备维护完成率数据更新成功")
}

func AddOrUpdateEquipmentRepairCompletionTimes(times []models.EquipmentRepairCompletionTimes, ctx *gin.Context) {
	yearAndMonth := time.Now().Format("2006-01")
	equipmentRepairCompletionTimesList, err := dao.GetEquipmentRepairCompletionTimes(yearAndMonth)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
	}

	var mapList = make(map[string]vo.EquipmentRepairCompletionTimesVo)

	for _, time := range equipmentRepairCompletionTimesList {
		mapList[time.Name] = vo.EquipmentRepairCompletionTimesVo{
			Name:       time.Name,
			Value:      time.Value,
			CreateTime: time.CreateTime,
		}
	}
	fmt.Println("mapList", mapList)

	for _, rate := range times {
		value, exists := mapList[rate.Name]
		if exists {
			// 检查是否是同一天的数据
			recordDate := strings.Split(strings.Split(value.CreateTime, " ")[0], "-")[0] + "-" + strings.Split(strings.Split(value.CreateTime, " ")[0], "-")[1]
			today := time.Now().Format("2006-01")
			if recordDate == today {
				// 同一天则更新
				var updateRate = models.EquipmentRepairCompletionTimes{
					Name:  rate.Name,
					Value: rate.Value,
				}
				updateTime := time.Now().Format("2006-01-02 15:04:05")
				err := dao.UpdateEquipmentRepairCompletionTimes(updateRate, updateTime)
				if err != nil {
					result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("更新设备 %s 的维护完成率失败: %v", rate.Name, err))
					return
				}
			} else {
				// 不是同一天则新增
				createTime := time.Now().Format("2006-01-02 15:04:05")
				err := dao.AddEquipmentRepairCompletionTimes(rate, createTime)
				if err != nil {
					result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("添加设备 %s 的维护完成率失败: %v", rate.Name, err))
					return
				}
			}
		} else {
			// 设备不存在则新增
			createTime := time.Now().Format("2006-01-02 15:04:05")
			err := dao.AddEquipmentRepairCompletionTimes(rate, createTime)
			if err != nil {
				result.Fail(ctx, http.StatusInternalServerError, fmt.Sprintf("添加设备 %s 的维护完成率失败: %v", rate.Name, err))
				return
			}
		}
		result.Success(ctx, "设备维修次数数据更新成功")
	}

}

func GetEquipmentRepairCompletionTimes(ctx *gin.Context) {
	yearAndMonth := time.Now().Format("2006-01")
	equipmentRepairCompletionTimes, err := dao.GetEquipmentRepairCompletionTimes(yearAndMonth)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
	}
	result.Success(ctx, equipmentRepairCompletionTimes)
}

func CurrentTimeToYearMonth(currentTime string) (int, int) {
	monthTime := strings.Split(currentTime, "-")
	year, err := strconv.Atoi(monthTime[0])

	if err != nil {
		zap.L().Error("year string转换int失败", zap.Error(err))
	}
	month, err := strconv.Atoi(monthTime[1])
	if err != nil {
		zap.L().Error("month string转换int失败", zap.Error(err))
	}
	return year, month
}

// 处理数据
func handleUtilizationList(
	utilizationWorkingHourList []models.UtilizationWorkingHour,
	utilizationExceptionList []models.UtilizationException,
	currentTime string) []vo.UtilizationVo {

	var utilizationList0 []models.Utilization
	var utilizationList1 []models.Utilization
	var utilizationVoList []vo.UtilizationVo
	for _, utilizationWorkingHour := range utilizationWorkingHourList {
		// 正常工作时间
		if utilizationWorkingHour.Type == "0" {
			utilizationList0 = append(utilizationList0, models.Utilization{
				Date:  utilizationWorkingHour.Time,
				Value: StringToFloat32(utilizationWorkingHour.WorkingHours),
				Line:  utilizationWorkingHour.Line,
			})
		} else {
			var AbnormalStopTime float32
			for _, exception := range utilizationExceptionList {
				//算出停线时间
				if exception.Time == utilizationWorkingHour.Time {
					AbnormalStopTime = CalculateTimeInterval(exception.AbnormalStopTime)
					fmt.Printf("异常停线时间%v:%v      %v\n", exception.Time, AbnormalStopTime, exception.AbnormalStopTime)
				}
			}
			utilizationList1 = append(utilizationList1, models.Utilization{
				Date:  utilizationWorkingHour.Time,
				Value: StringToFloat32(utilizationWorkingHour.WorkingHours) - AbnormalStopTime,
				Line:  utilizationWorkingHour.Line,
			})
		}
	}

	year, month := CurrentTimeToYearMonth(currentTime)
	days := DaysInMonth(year, time.Month(month))

	dailyData := make(map[string]vo.UtilizationVo)

	// 初始化每日数据
	for i := 1; i <= days; i++ {
		date := fmt.Sprintf("%d/%d/%d", year, month, i)
		dailyData[date] = vo.UtilizationVo{
			Date:  date,
			Value: 0,
		}
	}

	for _, utilization := range utilizationList0 {
		value0, count0 := CalculateValueOfDate(utilizationList0, utilization.Date)
		value1, _ := CalculateValueOfDate(utilizationList1, utilization.Date)

		if count0 > 0 && value0 != 0 {
			dailyData[utilization.Date] = vo.UtilizationVo{
				Date:  utilization.Date,
				Value: value1 / value0 / float32(count0),
			}
		}
	}
	// 将 map 转换为切片
	for i := 1; i <= days; i++ {
		date := fmt.Sprintf("%d/%d/%d", year, month, i)
		if data, exists := dailyData[date]; exists {
			date = strings.ReplaceAll(date, "/", "-")
			data.Date = date
			utilizationVoList = append(utilizationVoList, data)
		}
	}

	fmt.Println("utilizationVoList", utilizationVoList)
	return utilizationVoList
}

func CalculateValueOfDate(utilizationList []models.Utilization, date string) (float32, int) {
	var value float32
	lines := make(map[string]struct{})
	for _, utilization := range utilizationList {
		if utilization.Date == date {
			lines[utilization.Line] = struct{}{}
			value += utilization.Value
		}
	}
	return value, len(lines)
}

// DaysInMonth 获取指定年月的天数
func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

func StringToFloat32(s string) float32 {
	// 尝试将字符串转换为float32
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		// 如果转换失败，记录错误并返回0
		zap.L().Error("string转换float32失败",
			zap.String("value", s),
			zap.Error(err))
		return 0
	}
	return float32(value)
}

func CalculateTimeInterval(timeStr string) float32 {
	// 分割开始和结束时间
	times := strings.Split(timeStr, "-")
	if len(times) != 2 {
		zap.L().Error("时间格式错误",
			zap.String("value", timeStr))
		return 0
	}

	startTime := times[0]
	endTime := times[1]

	// 解析开始时间
	startParts := strings.Split(startTime, ":")
	if len(startParts) != 2 {
		zap.L().Error("开始时间格式错误",
			zap.String("value", startTime))
		return 0
	}
	startHour, err1 := strconv.Atoi(startParts[0])
	startMinute, err2 := strconv.Atoi(startParts[1])
	if err1 != nil || err2 != nil {
		zap.L().Error("开始时间解析失败",
			zap.String("value", startTime))
		return 0
	}

	// 解析结束时间
	endParts := strings.Split(endTime, ":")
	if len(endParts) != 2 {
		zap.L().Error("结束时间格式错误",
			zap.String("value", endTime))
		return 0
	}
	endHour, err3 := strconv.Atoi(endParts[0])
	endMinute, err4 := strconv.Atoi(endParts[1])
	if err3 != nil || err4 != nil {
		zap.L().Error("结束时间解析失败",
			zap.String("value", endTime))
		return 0
	}

	// 创建时间对象
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMinute, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMinute, 0, 0, time.Local)

	// 如果结束时间小于开始时间，说明跨天，需要加一天
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	// 计算时间差（时钟）
	duration := end.Sub(start)
	return float32(duration.Minutes()) / 60
}
