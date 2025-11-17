package dao

import (
	"fmt"
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

type EquipmentInspectionStatus struct {
	EquipCode     string `db:"equip_code"`
	EquipName     string `db:"equip_name"`
	EquipLocation string `db:"equip_location"`
	EquipType     string `db:"equip_type"`
	InspectedAm   bool   `db:"inspected_am"`
	InspectedPm   bool   `db:"inspected_pm"`
}

// GetEquipmentInspectionStatus 获取所有设备及当日点检状态
func GetEquipmentInspectionStatus() ([]models.EquipmentInspection, error) {
	var results []models.EquipmentInspection

	// 获取今天的日期，格式为 YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// SQL 查询：
	// 1. 从 equipment_inspection 表中选择所有设备。
	// 2. LEFT JOIN 当天的 Inspection 记录。
	// 3. 使用 CASE 语句判断当天是否有点检记录 (i.type IS NOT NULL)，并生成 inspected 状态 (1 或 0)。
	query := `
		SELECT 
			ei.seq, 
			ei.equip_name, 
			ei.equip_code, 
			ei.equip_location, 
			ei.equip_type,
			CASE 
				WHEN i.type IS NOT NULL THEN 1 
				ELSE 0 
			END AS inspected
		FROM 
			equipment_inspection ei
		LEFT JOIN 
			Inspection i ON ei.equip_code = i.type AND CONVERT(date, i.time) = @p1
		ORDER BY
			ei.seq`

	rows, err := config.DB.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("查询设备点检状态失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.EquipmentInspection
		if err := rows.Scan(&item.Seq, &item.EquipName, &item.EquipCode, &item.EquipLocation, &item.EquipType, &item.Inspected); err != nil {
			return nil, fmt.Errorf("扫描数据失败: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}
