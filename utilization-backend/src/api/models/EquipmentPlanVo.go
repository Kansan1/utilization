package models

type EquipmentPlanVo struct {
	ID          int    `json:"id"`
	Line        string `json:"line"`
	EquipmentID int    `json:"equipment_id"`
	Day         int    `json:"day"`
	Name        string `json:"name"` // 设备名称
	Type        string `json:"type"` // 设备类型

}
