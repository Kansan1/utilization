package router

import (
	"net/http"
	"strconv"
	"time"
	"utilization-backend/src/api/dao"
	"utilization-backend/src/api/models"
	"utilization-backend/src/api/result"
	"utilization-backend/src/api/v1/home"
	"utilization-backend/src/socket"

	"github.com/gin-gonic/gin"
)

func HomeRouter(r *gin.Engine) *gin.RouterGroup {
	apiHome := r.Group("/api")
	{
		home := apiHome.Group("/home")
		home.GET("/utilization/list", handleUtilizationList)
		home.GET("/utilization/all", handleAllUtilizationList)
		home.POST("/equipmentRepairCompletionRate/addOrUpdate", handleAddOrUpdateEquipmentRepairCompletionRate)
		home.GET("/equipmentRepairCompletionRate/list", handleGetEquipmentRepairCompletionRateList)
		//		home.DELETE("/equipmentRepairCompletionRate/delete", handleDeleteEquipmentRepairCompletionRate)
		home.GET("/equipmentRepairCompletionTimes/list", handleGetEquipmentRepairCompletionTimes)
		home.POST("/equipmentRepairCompletionTimes/addOrUpdate", handleAddOrUpdateEquipmentRepairCompletionTimes)
		//		home.DELETE("/equipmentRepairCompletionTimes/delete", handleDeleteEquipmentRepairCompletionTimes)
		home.GET("/equipment/list", GetAllEquipments)
		home.POST("/equipment/add", AddEquipment)
		home.PUT("/equipment/update", UpdateEquipment)
		home.DELETE("/equipment/delete", DeleteEquipment)

		home.GET("/equipmentPlan/list", GetAllEquipmentPlans)
		home.POST("/equipmentPlan/add", AddEquipmentPlan)
		home.POST("/equipmentPlan/update", UpdateEquipmentPlan)
		home.POST("/equipmentPlan/delete", DeleteEquipmentPlan)

		home.GET("/planSpecific/list", GetAllEquipmentPlanSpecific)
		home.POST("/planSpecific/add", AddEquipmentPlanSpecific)
		home.PUT("/planSpecific/update", UpdateEquipmentPlanSpecific)
		home.DELETE("/planSpecific/delete", DeleteEquipmentPlanSpecific)
		home.GET("/planSpecific/day", GetTodayMergedPlan)
		home.GET("/planSpecific/week", GetThisWeekMergedPlan)

		home.GET("/equipment/GetEquipmentTypeMaintainStats", GetEquipmentTypeMaintainStats)

		home.POST("/dailyRepair/add", AddDailyRepairTask)
		home.PUT("/dailyRepair/update", UpdateDailyRepairTask)
		home.DELETE("/dailyRepair/delete/:id", DeleteDailyRepairTask)
		home.GET("/dailyRepair/list", GetAllDailyRepairTasks)
		home.GET("/dailyRepair/todayList", GetTodayDailyRepairTasks)

		home.GET("/monthlyFaultFrequency/list", GetMonthlyFaultFrequencyList)         // 查询列表（当月）
		home.POST("/monthlyFaultFrequency/add", AddMonthlyFaultFrequency)             // 新增
		home.PUT("/monthlyFaultFrequency/update", UpdateMonthlyFaultFrequency)        // 更新
		home.DELETE("/monthlyFaultFrequency/delete/:id", DeleteMonthlyFaultFrequency) // 删除
		home.GET("/monthlyFaultFrequency/merged", GetMergedMonthlyFaultFrequencies)

	}
	return apiHome
}

func handleUtilizationList(ctx *gin.Context) {
	currentTime := ctx.Query("currentTime")
	home.GetUtilizationListByMonth(currentTime, ctx)

}

func handleAllUtilizationList(ctx *gin.Context) {
	home.GetAllUtilizationList(ctx)

}

func handleGetEquipmentRepairCompletionRateList(ctx *gin.Context) {
	home.GetEquipmentRepairCompletionRateList(ctx)
}

func handleDeleteEquipmentRepairCompletionRate(ctx *gin.Context) {
	var names []string
	if err := ctx.ShouldBindJSON(&names); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if len(names) == 0 {
		result.Fail(ctx, http.StatusBadRequest, "设备名称列表不能为空")
		return
	}

	home.DeleteEquipmentRepairCompletionRate(names, ctx)
}

func handleAddOrUpdateEquipmentRepairCompletionRate(ctx *gin.Context) {
	var rates []models.EquipmentRepairCompletionRate
	if err := ctx.ShouldBindJSON(&rates); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
	}
	home.AddOrUpdateEquipmentRepairCompletionRate(rates, ctx)
}

func handleGetEquipmentRepairCompletionTimes(ctx *gin.Context) {
	home.GetEquipmentRepairCompletionTimes(ctx)
}

func handleAddOrUpdateEquipmentRepairCompletionTimes(ctx *gin.Context) {
	var times []models.EquipmentRepairCompletionTimes
	if err := ctx.ShouldBindJSON(&times); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
	}
	home.AddOrUpdateEquipmentRepairCompletionTimes(times, ctx)
}

// 查询所有设备
func GetAllEquipments(ctx *gin.Context) {
	list, err := dao.GetAllEquipments()
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "查询设备失败："+err.Error())
		return
	}

	result.Success(ctx, list)
}

// 添加设备
func AddEquipment(ctx *gin.Context) {
	var e models.Equipment
	if err := ctx.ShouldBindJSON(&e); err != nil {
		result.Fail(ctx, http.StatusBadRequest, "参数解析失败："+err.Error())
		return
	}
	if err := dao.AddEquipment(e); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "添加设备失败："+err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "添加成功")
}

// 更新设备
func UpdateEquipment(ctx *gin.Context) {
	var e models.Equipment
	if err := ctx.ShouldBindJSON(&e); err != nil {
		result.Fail(ctx, http.StatusBadRequest, "参数解析失败："+err.Error())
		return
	}
	if err := dao.UpdateEquipment(e); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "更新设备失败："+err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "修改成功")
}

// 删除设备（支持单个或多个）
func DeleteEquipment(ctx *gin.Context) {
	var ids []int
	if err := ctx.ShouldBindJSON(&ids); err == nil && len(ids) > 0 {
		// 批量删除
		if err := dao.DeleteEquipmentsByIDs(ids); err != nil {
			result.Fail(ctx, http.StatusInternalServerError, "批量删除失败："+err.Error())
			return
		}
		result.Success(ctx, "批量删除成功")
		return
	}

	// 若未传 JSON 则尝试 query 参数删除
	idStr := ctx.Query("id")
	if idStr == "" {
		result.Fail(ctx, http.StatusBadRequest, "请提供设备 ID")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Fail(ctx, http.StatusBadRequest, "设备 ID 格式错误")
		return
	}
	if err := dao.DeleteEquipmentByID(id); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "删除失败："+err.Error())
		return
	}

	socket.NotifyAllClients("Change", "")

	result.Success(ctx, "删除成功")
}

// 查询所有设备维护计划
func GetAllEquipmentPlans(ctx *gin.Context) {
	list, err := dao.GetEquipmentPlanVoList()
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, list)
}

func AddEquipmentPlan(ctx *gin.Context) {
	var plan models.EquipmentPlan
	if err := ctx.ShouldBindJSON(&plan); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 判断是否已有记录
	exists, err := dao.CheckEquipmentPlanExists(plan.EquipmentID, plan.Day)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	if exists {
		result.Success(ctx, "添加成功")
		return
	}

	// 插入
	err = dao.AddEquipmentPlan(plan)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "添加成功")
}

// 更新设备维护计划
func UpdateEquipmentPlan(ctx *gin.Context) {
	var plan models.EquipmentPlan
	if err := ctx.ShouldBindJSON(&plan); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := dao.UpdateEquipmentPlan(plan)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "更新成功")
}

// 删除设备维护计划
func DeleteEquipmentPlan(ctx *gin.Context) {
	var req struct {
		EquipmentID int `json:"equipment_id"`
		Day         int `json:"day"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Fail(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	//if req.EquipmentID == 0 || req.Day == 0 {
	//	result.Fail(ctx, http.StatusBadRequest, "设备ID和日期不能为空")
	//	return
	//}

	err := dao.DeleteEquipmentPlanByEquipmentIDAndDay(req.EquipmentID, req.Day)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	socket.NotifyAllClients("Change", "")

	result.Success(ctx, "删除成功")
}

func GetAllEquipmentPlanSpecific(ctx *gin.Context) {
	data, err := dao.GetAllEquipmentPlanSpecific()
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, data)
}

func AddEquipmentPlanSpecific(ctx *gin.Context) {
	var p models.EquipmentPlanSpecific
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := dao.AddEquipmentPlanSpecific(p); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "添加成功")
}

func UpdateEquipmentPlanSpecific(ctx *gin.Context) {
	var p models.EquipmentPlanSpecific
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 先查有没有这个ID的记录
	existing, err := dao.GetEquipmentPlanSpecificByID(p.ID)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if existing == nil {
		// 不存在，新增
		err = dao.AddEquipmentPlanSpecific(p)
	} else {
		// 存在，更新
		err = dao.UpdateEquipmentPlanSpecific(p)
	}

	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "操作成功")
}

func DeleteEquipmentPlanSpecific(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Fail(ctx, http.StatusBadRequest, "ID参数无效")
		return
	}
	if err := dao.DeleteEquipmentPlanSpecific(id); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "删除成功")
}

func GetTodayMergedPlan(ctx *gin.Context) {
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()

	list, err := dao.GetMergedPlanByDay(year, month, day)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, list)
}

func GetThisWeekMergedPlan(ctx *gin.Context) {
	now := time.Now()
	year, month, today := now.Year(), int(now.Month()), now.Day()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	var weekDays []int
	for i := 0; i < 7; i++ {
		day := today - (weekday - 1) + i
		if day > 0 {
			weekDays = append(weekDays, day)
		}
	}

	list, err := dao.GetMergedPlanByWeek(year, month, weekDays)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result.Success(ctx, list)
}

func GetEquipmentTypeMaintainStats(ctx *gin.Context) {
	now := time.Now()
	year, month := now.Year(), int(now.Month())

	stats, err := dao.GetEquipmentTypeMaintainStatsByMonth(year, month)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, stats)
}

func AddDailyRepairTask(ctx *gin.Context) {
	var task models.DailyRepairTask
	if err := ctx.ShouldBindJSON(&task); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 设置当前时间字符串作为报修时间
	task.ReportTime = time.Now().Format("2006-01-02 15:04:05")

	if err := dao.AddDailyRepairTask(task); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "添加成功")
}

func UpdateDailyRepairTask(ctx *gin.Context) {
	var task models.DailyRepairTask
	if err := ctx.ShouldBindJSON(&task); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := dao.UpdateDailyRepairTask(task); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "更新成功")
}

func DeleteDailyRepairTask(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, _ := strconv.Atoi(idstr)
	if err := dao.DeleteDailyRepairTask(id); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "删除成功")
}

func GetAllDailyRepairTasks(ctx *gin.Context) {
	tasks, err := dao.GetAllDailyRepairTasks()
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, tasks)
}

func GetTodayDailyRepairTasks(ctx *gin.Context) {
	tasks, err := dao.GetTodayDailyRepairTasks()
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, tasks)
}

func AddMonthlyFaultFrequency(ctx *gin.Context) {
	var m models.MonthlyFaultFrequency
	if err := ctx.ShouldBindJSON(&m); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	m.Year = now.Year()
	m.Month = int(now.Month())

	if err := dao.AddMonthlyFaultFrequency(m); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "新增成功")
}

func UpdateMonthlyFaultFrequency(ctx *gin.Context) {
	var m models.MonthlyFaultFrequency
	if err := ctx.ShouldBindJSON(&m); err != nil {
		result.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	m.Year = now.Year()
	m.Month = int(now.Month())

	if err := dao.UpdateMonthlyFaultFrequency(m); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", "")
	result.Success(ctx, "更新成功")
}

func DeleteMonthlyFaultFrequency(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Fail(ctx, http.StatusBadRequest, "无效ID")
		return
	}

	if err := dao.DeleteMonthlyFaultFrequency(id); err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	socket.NotifyAllClients("Change", id)

	result.Success(ctx, "删除成功")
}

func GetMonthlyFaultFrequencyList(ctx *gin.Context) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	list, err := dao.GetMonthlyFaultFrequencyList(year, month)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result.Success(ctx, list)
}

func GetMergedMonthlyFaultFrequencies(ctx *gin.Context) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	list, err := dao.GetMergedMonthlyFaultFrequencies(year, month)
	if err != nil {
		result.Fail(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}
	result.Success(ctx, list)
}
