package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

func GetAllEquipmentPlanSpecific() ([]models.EquipmentPlanSpecific, error) {
	query := `SELECT id, year, month, line, name, content, state, defender,shu FROM EquipmentPlanSpecific ORDER BY year DESC, month DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EquipmentPlanSpecific
	for rows.Next() {
		var p models.EquipmentPlanSpecific
		if err := rows.Scan(&p.ID, &p.Year, &p.Month, &p.Line, &p.EquipmentID, &p.Content, &p.State, &p.Defender, &p.Shu); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
func AddEquipmentPlanSpecific(p models.EquipmentPlanSpecific) error {
	// 获取当前时间
	now := time.Now()
	p.Year = now.Year()
	p.Month = int(now.Month())
	p.Day = now.Day()

	query := `
		INSERT INTO EquipmentPlanSpecific (year, month, day, line, equipment_id, content, state, defender,shu)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8,@p9)
	`
	_, err := config.DB.Exec(query, p.Year, p.Month, p.Day, p.Line, p.EquipmentID, p.Content, p.State, p.Defender, p.Shu)
	return err
}

func GetEquipmentPlanSpecificByID(id int) (*models.EquipmentPlanSpecific, error) {
	query := `SELECT id, year, month, day, line, equipment_id, content, state, defender,shu FROM EquipmentPlanSpecific WHERE id = @p1`
	row := config.DB.QueryRow(query, id)
	var p models.EquipmentPlanSpecific
	err := row.Scan(&p.ID, &p.Year, &p.Month, &p.Day, &p.Line, &p.EquipmentID, &p.Content, &p.State, &p.Defender, &p.Shu)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateEquipmentPlanSpecific(p models.EquipmentPlanSpecific) error {
	query := `
		UPDATE EquipmentPlanSpecific 
		SET year = @p1, month = @p2, day = @p3, line = @p4, equipment_id = @p5, content = @p6, state = @p7, defender = @p8 , shu = @p9 
		WHERE id = @p10
	`
	_, err := config.DB.Exec(query, p.Year, p.Month, p.Day, p.Line, p.EquipmentID, p.Content, p.State, p.Defender, p.Shu, p.ID)
	return err
}

func DeleteEquipmentPlanSpecific(id int) error {
	query := `DELETE FROM EquipmentPlanSpecific WHERE id = @p1`
	_, err := config.DB.Exec(query, id)
	return err
}

func GetEquipmentPlanSpecificByDay(year, month, day int) ([]models.EquipmentPlanSpecific, error) {
	query := `
		SELECT id, year, month, day, line, equipment_id, content, state, defender,shu
		FROM EquipmentPlanSpecific 
		WHERE year = @p1 AND month = @p2 AND day = @p3
	`
	rows, err := config.DB.Query(query, year, month, day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EquipmentPlanSpecific
	for rows.Next() {
		var p models.EquipmentPlanSpecific
		err := rows.Scan(&p.ID, &p.Year, &p.Month, &p.Day, &p.Line, &p.EquipmentID, &p.Content, &p.State, &p.Defender, &p.Shu)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func GetEquipmentPlanSpecificByWeek(year, month int, days []int) ([]models.EquipmentPlanSpecific, error) {
	if len(days) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(days))
	args := make([]interface{}, len(days)+2)
	args[0], args[1] = year, month
	for i, d := range days {
		placeholders[i] = "@p" + strconv.Itoa(i+3)
		args[i+2] = d
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT eps.id, eps.year, eps.month, eps.day, eps.line, eps.equipment_id, eps.content, eps.state, eps.defender
		FROM EquipmentPlanSpecific eps
		JOIN EquipmentPlan ep ON eps.equipment_id = ep.equipment_id
		WHERE eps.year = @p1 AND eps.month = @p2 AND ep.day IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EquipmentPlanSpecific
	for rows.Next() {
		var p models.EquipmentPlanSpecific
		err := rows.Scan(&p.ID, &p.Year, &p.Month, &p.Day, &p.Line, &p.EquipmentID, &p.Content, &p.State, &p.Defender)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

func GetMergedPlanByDay(year, month, day int) ([]models.EquipmentPlanSpecificVO, error) {
	// 基础计划（来自 EquipmentPlan + Equipment）
	basePlans, err := GetEquipmentPlanVoByDay(day)
	if err != nil {
		return nil, err
	}

	// 详细计划（来自 EquipmentPlanSpecific）
	specificPlans, err := GetEquipmentPlanSpecificByDay(year, month, day)
	if err != nil {
		return nil, err
	}

	// 转 map，方便查找覆盖
	specificMap := make(map[int]models.EquipmentPlanSpecific)
	for _, sp := range specificPlans {
		specificMap[sp.EquipmentID] = sp
	}

	var merged []models.EquipmentPlanSpecificVO

	for _, base := range basePlans {
		if detail, ok := specificMap[base.EquipmentID]; ok {
			// 有详细计划，用详细的
			merged = append(merged, models.EquipmentPlanSpecificVO{
				ID:            detail.ID,
				Line:          detail.Line,
				EquipmentID:   detail.EquipmentID,
				EquipmentName: base.Name,
				Content:       detail.Content,
				State:         detail.State,
				Defender:      detail.Defender,
				Day:           detail.Day,
				Month:         detail.Month,
				Year:          detail.Year,
				Shu:           detail.Shu,
			})
		} else {
			// 没有详细计划，保留基础的
			merged = append(merged, models.EquipmentPlanSpecificVO{
				ID:            0,
				Line:          base.Line,
				EquipmentID:   base.EquipmentID,
				EquipmentName: base.Name,
				Content:       "",
				State:         "待维护",
				Defender:      "",
				Day:           day,
				Month:         month,
				Year:          year,
				Shu:           0,
			})
		}
	}

	return merged, nil
}

func GetMergedPlanByWeek(year, month int, days []int) ([]models.EquipmentPlanSpecificVO, error) {
	resultMap := make(map[int]models.EquipmentPlanSpecificVO)

	for _, day := range days {
		// 查询基础计划
		basePlans, err := GetEquipmentPlanVoByDay(day)
		if err != nil {
			return nil, err
		}

		// 查询详细计划
		specificPlans, err := GetEquipmentPlanSpecificByDay(year, month, day)
		if err != nil {
			return nil, err
		}

		// 将详细计划转为 map，方便查找
		specificMap := make(map[int]models.EquipmentPlanSpecific)
		for _, sp := range specificPlans {
			specificMap[sp.EquipmentID] = sp
		}

		for _, base := range basePlans {
			// 如果已经存在且是详细计划（ID != 0），则跳过
			existing, exists := resultMap[base.EquipmentID]
			if exists && existing.ID != 0 {
				continue
			}

			if detail, ok := specificMap[base.EquipmentID]; ok {
				// 有详细计划，优先使用
				resultMap[base.EquipmentID] = models.EquipmentPlanSpecificVO{
					ID:            detail.ID,
					Line:          detail.Line,
					EquipmentID:   detail.EquipmentID,
					EquipmentName: base.Name,
					Content:       detail.Content,
					State:         detail.State,
					Defender:      detail.Defender,
					Day:           detail.Day,
					Month:         detail.Month,
					Year:          detail.Year,
				}
			} else if !exists {
				// 没有任何数据，使用基础计划作为占位
				resultMap[base.EquipmentID] = models.EquipmentPlanSpecificVO{
					ID:            0,
					Line:          base.Line,
					EquipmentID:   base.EquipmentID,
					EquipmentName: base.Name,
					Content:       "",
					State:         "待维护",
					Defender:      "",
					Day:           day,
					Month:         month,
					Year:          year,
				}
			}
		}
	}

	// 转换为切片返回
	var merged []models.EquipmentPlanSpecificVO
	for _, v := range resultMap {
		merged = append(merged, v)
	}

	return merged, nil
}

func GetEquipmentTypeMaintainStatsByMonth(year, month int) ([]models.EquipmentTypeMaintainStatVO, error) {
	query := `
        SELECT 
            e.type AS equipment_type,
            SUM(e.qty) AS total_count,
            ISNULL(SUM(eps.shu), 0) AS completed_count
        FROM Equipment e
        LEFT JOIN EquipmentPlanSpecific eps 
            ON e.id = eps.equipment_id 
            AND eps.year = @p1 
            AND eps.month = @p2 
            AND eps.state = '完成'
        GROUP BY e.type
    `
	rows, err := config.DB.Query(query, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EquipmentTypeMaintainStatVO
	for rows.Next() {
		var stat models.EquipmentTypeMaintainStatVO
		err := rows.Scan(&stat.EquipmentType, &stat.TotalCount, &stat.CompletedCount)
		if err != nil {
			return nil, err
		}
		results = append(results, stat)
	}
	return results, nil
}
