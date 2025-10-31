package dao

import (
	"log"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

// 添加设备维护计划
func AddEquipmentPlan(plan models.EquipmentPlan) error {
	query := `INSERT INTO EquipmentPlan (equipment_id, day) VALUES (@p1, @p2)`
	_, err := config.DB.Exec(query, plan.EquipmentID, plan.Day)
	if err != nil {
		log.Printf("添加设备维护计划失败: %v", err)
	}
	return err
}

// 根据ID更新设备维护计划
func UpdateEquipmentPlan(plan models.EquipmentPlan) error {
	query := `UPDATE EquipmentPlan SET equipment_id = @p1, day = @p2 WHERE id = @p4`
	_, err := config.DB.Exec(query, plan.EquipmentID, plan.Day, plan.ID)
	if err != nil {
		log.Printf("更新设备维护计划失败: %v", err)
	}
	return err
}

// 根据ID删除设备维护计划
func DeleteEquipmentPlanByID(id int) error {
	query := `DELETE FROM EquipmentPlan WHERE id = @p1`
	_, err := config.DB.Exec(query, id)
	if err != nil {
		log.Printf("删除设备维护计划失败: %v", err)
	}
	return err
}

func CheckEquipmentPlanExists(equipmentID int, day int) (bool, error) {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(*) FROM EquipmentPlan WHERE equipment_id = @p1 AND day = @p2
	`, equipmentID, day).Scan(&count)

	if err != nil {
		log.Printf("查询计划是否存在失败: %v", err)
		return false, err
	}
	return count > 0, nil
}

func DeleteEquipmentPlanByEquipmentIDAndDay(equipmentID, day int) error {
	query := `DELETE FROM EquipmentPlan WHERE equipment_id = @p1 AND day = @p2`
	_, err := config.DB.Exec(query, equipmentID, day)
	if err != nil {
		log.Printf("删除设备维护计划失败: %v", err)
	}
	return err
}

func GetEquipmentPlanVoList() ([]models.EquipmentPlanVo, error) {
	query := `
        SELECT 
            p.id, p.equipment_id, p.day,
            e.name, e.type,ISNULL(e.Line, '')
        FROM EquipmentPlan p
        LEFT JOIN Equipment e ON p.equipment_id = e.id
        ORDER BY p.equipment_id, p.day
    `

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EquipmentPlanVo
	for rows.Next() {
		var plan models.EquipmentPlanVo
		err = rows.Scan(&plan.ID, &plan.EquipmentID, &plan.Day, &plan.Name, &plan.Type, &plan.Line) // 👈 直接用 bool
		if err != nil {
			return nil, err
		}
		list = append(list, plan)
	}
	return list, nil
}

// 查询当天需要维护的设备（关联设备表，查当天的 day）
func GetEquipmentPlanVoByDay(day int) ([]models.EquipmentPlanVo, error) {
	query := `
        SELECT ep.id, ep.equipment_id, ep.day, e.name, e.type,ISNULL(e.Line, '')
        FROM EquipmentPlan ep
        JOIN Equipment e ON ep.equipment_id = e.id
        WHERE ep.day = @p1
    `
	rows, err := config.DB.Query(query, day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.EquipmentPlanVo
	for rows.Next() {
		var vo models.EquipmentPlanVo
		err := rows.Scan(&vo.ID, &vo.EquipmentID, &vo.Day, &vo.Name, &vo.Type, &vo.Line)
		if err != nil {
			return nil, err
		}
		list = append(list, vo)
	}
	return list, nil
}
