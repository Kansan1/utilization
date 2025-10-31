package dao

import (
	"database/sql"
	"errors"
	"log"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
)

// 通用 SQL 执行函数（带日志）
func execSQL(query string, args ...interface{}) error {
	_, err := config.DB.Exec(query, args...)
	if err != nil {
		log.Printf("SQL执行失败: %s，参数: %+v，错误: %v", query, args, err)
	}
	return err
}

// 添加设备
func AddEquipment(e models.Equipment) error {
	query := `INSERT INTO Equipment (name, type, qty,line) VALUES (@p2, @p3, @p4,@p5)`
	return execSQL(query, e.ID, e.Name, e.Type, e.Qty, e.Line)
}

// 修改设备
func UpdateEquipment(e models.Equipment) error {
	query := `UPDATE Equipment SET name = @p1, type = @p2, qty = @p3,line = @p5 WHERE id = @p4`
	return execSQL(query, e.Name, e.Type, e.Qty, e.ID, e.Line)
}

// 删除设备（根据ID）
func DeleteEquipmentByID(id int) error {
	query := `DELETE FROM Equipment WHERE id = @p1`
	return execSQL(query, id)
}

// 批量删除设备（传多个ID）
func DeleteEquipmentsByIDs(ids []int) error {
	query := `DELETE FROM Equipment WHERE id = @p1`
	for _, id := range ids {
		if err := execSQL(query, id); err != nil {
			return err
		}
	}
	return nil
}

// 查询所有设备
func GetAllEquipments() ([]models.Equipment, error) {
	var list []models.Equipment
	query := `SELECT id, name, type, qty, ISNULL(Line, '') AS Line
FROM Equipment
ORDER BY type DESC , name`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Printf("查询所有设备失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Equipment
		err = rows.Scan(&e.ID, &e.Name, &e.Type, &e.Qty, &e.Line)
		if err != nil {
			log.Printf("扫描设备数据失败: %v", err)
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

// 根据 ID 查询设备
func GetEquipmentByID(id int) (*models.Equipment, error) {
	var e models.Equipment
	query := `SELECT id, name, type, qty,line FROM Equipment WHERE id = ?`
	err := config.DB.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Type, &e.Qty, &e.Line)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 查不到不是错误
		}
		log.Printf("查询设备 ID=%d 失败: %v", id, err)
		return nil, err
	}
	return &e, nil
}
