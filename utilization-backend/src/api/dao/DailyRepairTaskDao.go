package dao

import (
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

func AddDailyRepairTask(task models.DailyRepairTask) error {
	query := `
		INSERT INTO DailyRepairTask (line, device_name, fault, report_time, state, repairer)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6)
	`
	_, err := config.DB.Exec(query, task.Line, task.DeviceName, task.Fault, task.ReportTime, task.State, task.Repairer)
	return err
}

func UpdateDailyRepairTask(task models.DailyRepairTask) error {
	query := `
		UPDATE DailyRepairTask
		SET line = @p1, device_name = @p2, fault = @p3, report_time = @p4, state = @p5, repairer = @p6
		WHERE id = @p7
	`
	_, err := config.DB.Exec(query, task.Line, task.DeviceName, task.Fault, task.ReportTime, task.State, task.Repairer, task.ID)
	return err
}

func DeleteDailyRepairTask(id int) error {
	query := `DELETE FROM DailyRepairTask WHERE id = @p1`
	_, err := config.DB.Exec(query, id)
	return err
}

func GetAllDailyRepairTasks() ([]models.DailyRepairTask, error) {
	query := `SELECT id, line, device_name, fault, report_time, state, repairer FROM DailyRepairTask`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DailyRepairTask
	for rows.Next() {
		var task models.DailyRepairTask
		if err := rows.Scan(&task.ID, &task.Line, &task.DeviceName, &task.Fault, &task.ReportTime, &task.State, &task.Repairer); err != nil {
			return nil, err
		}
		list = append(list, task)
	}
	return list, nil
}

func GetTodayDailyRepairTasks() ([]models.DailyRepairTask, error) {
	today := time.Now().Format("2006-01-02") // 获取今天的日期字符串，形如 "2025-06-25"

	query := `
		SELECT id, line, device_name, fault, report_time, state, repairer
		FROM DailyRepairTask
		WHERE CONVERT(date, report_time) = @p1
	`

	rows, err := config.DB.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.DailyRepairTask
	for rows.Next() {
		var task models.DailyRepairTask
		if err := rows.Scan(&task.ID, &task.Line, &task.DeviceName, &task.Fault, &task.ReportTime, &task.State, &task.Repairer); err != nil {
			return nil, err
		}
		list = append(list, task)
	}
	return list, nil
}
