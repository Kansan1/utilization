package dao

import (
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

// AddInspection 添加点检记录
func AddInspection(inspection models.Inspection) error {
	// 使用当前时间作为点检时间
	query := `INSERT INTO Inspection (time, type) VALUES (@p1, @p2)`
	_, err := config.DB.Exec(query, time.Now().Format("2006-01-02 15:04:05"), inspection.Type)
	return err
}
