package dao

import (
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

// 新增
func AddMonthlyFaultFrequency(m models.MonthlyFaultFrequency) error {
	query := `
        INSERT INTO MonthlyFaultFrequency (line, equipment_name, content, fault_count, state, defender, year, month)
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)
    `
	_, err := config.DB.Exec(query,
		m.Line, m.EquipmentName, m.Content, m.FaultCount, m.State, m.Defender, m.Year, m.Month)
	return err
}

// 更新（根据id）
func UpdateMonthlyFaultFrequency(m models.MonthlyFaultFrequency) error {
	query := `
        UPDATE MonthlyFaultFrequency
        SET line = @p1, equipment_name = @p2, content = @p3, fault_count = @p4, state = @p5, defender = @p6, year = @p7, month = @p8
        WHERE id = @p9
    `
	_, err := config.DB.Exec(query,
		m.Line, m.EquipmentName, m.Content, m.FaultCount, m.State, m.Defender, m.Year, m.Month, m.ID)
	return err
}

// 删除（根据id）
func DeleteMonthlyFaultFrequency(id int) error {
	query := `DELETE FROM MonthlyFaultFrequency WHERE id = @p1`
	_, err := config.DB.Exec(query, id)
	return err
}

// 查询（分页或者全部）- 这里示例查询全部
func GetMonthlyFaultFrequencyList(year, month int) ([]models.MonthlyFaultFrequency, error) {
	query := `
        SELECT id, line, equipment_name, content, fault_count, state, defender, year, month
        FROM MonthlyFaultFrequency
        WHERE year = @p1 AND month = @p2
    `
	rows, err := config.DB.Query(query, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MonthlyFaultFrequency
	for rows.Next() {
		var m models.MonthlyFaultFrequency
		err := rows.Scan(&m.ID, &m.Line, &m.EquipmentName, &m.Content, &m.FaultCount, &m.State, &m.Defender, &m.Year, &m.Month)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// GetMergedMonthlyFaultFrequencies 查询合并后的频繁故障数据
func GetMergedMonthlyFaultFrequencies(year, month int) ([]models.MonthlyFaultFrequency, error) {
	resultMap := make(map[string]models.MonthlyFaultFrequency)

	// 查询 MonthlyFaultFrequency 表中已有的数据
	existingQuery := `
        SELECT id, line, equipment_name, content, fault_count, state, defender, year, month
        FROM MonthlyFaultFrequency
        WHERE year = @p1 AND month = @p2
    `
	rows, err := config.DB.Query(existingQuery, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.MonthlyFaultFrequency
		if err := rows.Scan(&m.ID, &m.Line, &m.EquipmentName, &m.Content, &m.FaultCount, &m.State, &m.Defender, &m.Year, &m.Month); err != nil {
			return nil, err
		}
		// 用 线别+设备名称 作为key（去重标识）
		key := m.Line + "|" + m.EquipmentName
		resultMap[key] = m
	}

	// 查询 DailyRepairTask 中本月维修次数 ≥ 5 的记录
	freqQuery := `
        SELECT line, device_name, fault, COUNT(*) AS fault_count
        FROM DailyRepairTask
        WHERE YEAR(report_time) = @p1 AND MONTH(report_time) = @p2
        GROUP BY line, device_name, fault
        HAVING COUNT(*) >= 5
    `
	rows2, err := config.DB.Query(freqQuery, year, month)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var m models.MonthlyFaultFrequency
		var faultCount int
		if err := rows2.Scan(&m.Line, &m.EquipmentName, &m.Content, &faultCount); err != nil {
			return nil, err
		}
		m.FaultCount = faultCount
		m.State = "未处理"
		m.Defender = ""
		m.Year = year
		m.Month = month

		key := m.Line + "|" + m.EquipmentName
		// 如果已存在，不重复加入；不存在才加入
		if _, exists := resultMap[key]; !exists {
			resultMap[key] = m
		}
	}

	// 转换 map 为 slice
	var result []models.MonthlyFaultFrequency
	for _, v := range resultMap {
		result = append(result, v)
	}

	return result, nil
}
