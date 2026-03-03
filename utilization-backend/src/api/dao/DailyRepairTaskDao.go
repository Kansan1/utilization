package dao

import (
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

func AddDailyRepairTask(task models.DailyRepairTask) error {
	query := `
		INSERT INTO Repair_Report (报修班组, 型号, 故障现象, 报修时间, 状态, 维修人)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6)
	`
	_, err := config.DB2.Exec(query, task.Line, task.DeviceName, task.Fault, task.ReportTime, task.State, task.Repairer)
	return err
}

func UpdateDailyRepairTask(task models.DailyRepairTask) error {
	query := `
		UPDATE Repair_Report
		SET 报修班组 = @p1, 型号 = @p2, 故障现象 = @p3, 报修时间 = @p4, 状态 = @p5, 维修人 = @p6
		WHERE ID = @p7
	`
	_, err := config.DB2.Exec(query, task.Line, task.DeviceName, task.Fault, task.ReportTime, task.State, task.Repairer, task.ID)
	return err
}

func DeleteDailyRepairTask(id int) error {
	query := `DELETE FROM Repair_Report WHERE ID = @p1`
	_, err := config.DB2.Exec(query, id)
	return err
}

func GetAllDailyRepairTasks() ([]models.DailyRepairTask, error) {
	query := `
		SELECT 
			ID, 
			COALESCE(报修班组, '') as 报修班组, 
			COALESCE(型号, '') as 型号, 
			COALESCE(故障现象, '') as 故障现象, 
			报修时间, 
			COALESCE(状态, '') as 状态, 
			COALESCE(维修人, '') as 维修人 
		FROM Repair_Report
	`
	rows, err := config.DB2.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DailyRepairTask
	for rows.Next() {
		var task models.DailyRepairTask
		var reportTime *time.Time
		if err := rows.Scan(&task.ID, &task.Line, &task.DeviceName, &task.Fault, &reportTime, &task.State, &task.Repairer); err != nil {
			return nil, err
		}
		if reportTime != nil {
			task.ReportTime = reportTime.Format("2006-01-02 15:04:05")
		} else {
			task.ReportTime = ""
		}
		list = append(list, task)
	}
	return list, nil
}

func GetTodayDailyRepairTasks() ([]models.DailyRepairTask, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT 
			ID, 
			COALESCE(报修班组, '') as 报修班组, 
			COALESCE(型号, '') as 型号, 
			COALESCE(故障现象, '') as 故障现象, 
			报修时间, 
			COALESCE(状态, '') as 状态, 
			COALESCE(维修人, '') as 维修人
		FROM Repair_Report
		WHERE CONVERT(date, 报修时间) = @p1
	`

	rows, err := config.DB2.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DailyRepairTask
	for rows.Next() {
		var task models.DailyRepairTask
		var reportTime *time.Time
		if err := rows.Scan(&task.ID, &task.Line, &task.DeviceName, &task.Fault, &reportTime, &task.State, &task.Repairer); err != nil {
			return nil, err
		}
		if reportTime != nil {
			task.ReportTime = reportTime.Format("2006-01-02 15:04:05")
		} else {
			task.ReportTime = ""
		}
		list = append(list, task)
	}
	return list, nil
}
