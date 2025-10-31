package dao

import (
	"fmt"
	"log"
	"utilization-backend/config"
	"utilization-backend/src/api/models"
	"utilization-backend/src/api/vo"
)

// AddEquipmentRepairCompletionRate 添加设备维护完成率数据
func AddEquipmentRepairCompletionRate(rate models.EquipmentRepairCompletionRate, time string) error {

	// 使用参数化查询防止SQL注入
	_, err := config.DB.Exec(`
			INSERT INTO equipment_repair_completion_rate (name, value, create_time)
			VALUES (@p1, CAST(@p2 AS DECIMAL(10,4)), @p3)`,
		rate.Name, rate.Value, time)

	if err != nil {
		log.Printf("添加设备维护完成率数据失败: %v", err)
		return err

	}
	return nil
}

func DeleteEquipmentRepairCompletionRate(names []string) error {
	for _, name := range names {
		_, err := config.DB.Exec(`
		DELETE FROM equipment_repair_completion_rate 
		WHERE name = @p1`,
			name)
		if err != nil {
			log.Printf("删除设备维护完成率数据失败: %v", err)
			return err
		}
	}
	return nil
}

func DeleteEquipmentRepairCompletionRateEachMonth() error {
	_, err := config.DB.Exec(`
		TRUNCATE TABLE equipment_repair_completion_rate`)
	if err != nil {
		log.Printf("清空设备维护完成率数据失败: %v", err)
		return err
	}
	return nil
}

func UpdateEquipmentRepairCompletionRate(rate models.EquipmentRepairCompletionRate, time string) error {
	_, err := config.DB.Exec(`
		UPDATE equipment_repair_completion_rate 
		SET value = CAST(@p2 AS DECIMAL(10,4)), create_time = @p3 
		WHERE name = @p1`,
		rate.Name, rate.Value, time)
	if err != nil {
		log.Printf("更新设备维护完成率数据失败: %v", err)
		return err
	}
	return nil
}

func AddEquipmentRepairCompletionTimes(rate models.EquipmentRepairCompletionTimes, time string) error {
	_, err := config.DB.Exec(`
		INSERT INTO equipment_repair_completion_times (name, value, create_time)
		VALUES (@p1, @p2, @p3)`,
		rate.Name, rate.Value, time)
	if err != nil {
		log.Printf("添加设备维护完成次数数据失败: %v", err)
		return err
	}
	return nil
}

func UpdateEquipmentRepairCompletionTimes(rate models.EquipmentRepairCompletionTimes, time string) error {
	_, err := config.DB.Exec(`
		UPDATE equipment_repair_completion_times 
		SET value = @p2, create_time = @p3 
		WHERE name = @p1`,
		rate.Name, rate.Value, time)
	if err != nil {
		log.Printf("更新设备维护完成次数数据失败: %v", err)
		return err
	}
	return nil
}

func GetEquipmentRepairCompletionTimes(yearAndMonth string) ([]vo.EquipmentRepairCompletionTimesVo, error) {
	var times []vo.EquipmentRepairCompletionTimesVo
	rows, err := config.DB.Query("select name,value,create_time from equipment_repair_completion_times where create_time like @p1", yearAndMonth+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var time vo.EquipmentRepairCompletionTimesVo
		err = rows.Scan(&time.Name, &time.Value, &time.CreateTime)
		if err != nil {
			return nil, err
		}
		times = append(times, time)
	}
	return times, nil
}

func GetEquipmentRepairCompletionRateList(yearAndMonth string) ([]vo.EquipmentRepairCompletionRateVo, error) {
	var rates []vo.EquipmentRepairCompletionRateVo
	fmt.Println("yearAndMonth", yearAndMonth)
	rows, err := config.DB.Query("select name,value,create_time from equipment_repair_completion_rate where create_time like @p1", yearAndMonth+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rate vo.EquipmentRepairCompletionRateVo
		err = rows.Scan(&rate.Name, &rate.Value, &rate.CreateTime)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}
