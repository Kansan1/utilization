package dao

import (
	"fmt"
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

// GetEquipmentInspectionStatus 获取所有设备及当日点检状态
func GetEquipmentInspectionStatus() ([]models.EquipmentInspection, error) {
	var results []models.EquipmentInspection

	// 获取今天的日期，格式为 YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// SQL 查询更新：
	// 1. 增加了 e.seq 的查询
	// 2. 为 BYQ 设备分别检查 _AM 和 _PM 后缀的记录
	// 3. 为非 BYQ 设备检查无后缀的记录到 inspected_am 字段
	query := `
SELECT
    e.seq,
    e.equip_name,
    e.equip_code,
    e.equip_location,
    e.equip_type,
    CAST(
        CASE
            WHEN e.equip_code LIKE 'BYQ%' THEN
                CASE WHEN EXISTS (SELECT 1 FROM Inspection i WHERE i.type = e.equip_code + '_AM' AND CONVERT(date, i.time) = @p1) THEN 1 ELSE 0 END
            ELSE
                CASE WHEN EXISTS (SELECT 1 FROM Inspection i WHERE i.type = e.equip_code AND CONVERT(date, i.time) = @p1) THEN 1 ELSE 0 END
        END
    AS BIT) AS inspected_am,
    CAST(
        CASE
            WHEN e.equip_code LIKE 'BYQ%' THEN
                CASE WHEN EXISTS (SELECT 1 FROM Inspection i WHERE i.type = e.equip_code + '_PM' AND CONVERT(date, i.time) = @p1) THEN 1 ELSE 0 END
            ELSE
                0
        END
    AS BIT) AS inspected_pm
FROM
    equipment_inspection e
ORDER BY
    e.seq;
`

	rows, err := config.DB.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("查询设备点检状态失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.EquipmentInspection
		if err := rows.Scan(&item.Seq, &item.EquipName, &item.EquipCode, &item.EquipLocation, &item.EquipType, &item.InspectedAm, &item.InspectedPm); err != nil {
			return nil, fmt.Errorf("扫描数据失败: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}
