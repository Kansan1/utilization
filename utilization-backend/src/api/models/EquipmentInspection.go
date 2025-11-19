package models

// EquipmentInspection 对应 equipment_inspection 表
type EquipmentInspection struct {
	Seq           int    `json:"seq"`
	EquipName     string `json:"equip_name"`
	EquipCode     string `json:"equip_code"`
	EquipLocation string `json:"equip_location"`
	EquipType     string `json:"equip_type"`
	InspectedAm   bool   `json:"inspected_am"`
	InspectedPm   bool   `json:"inspected_pm"`
}

// TableName 指定表名
func (EquipmentInspection) TableName() string {
	return "equipment_inspection"
}
