package models

// EquipmentInspection 对应 equipment_inspection 表
type EquipmentInspection struct {
	Seq           int    `json:"seq"`
	EquipName     string `json:"equip_name"`
	EquipCode     string `json:"equip_code"`
	EquipLocation string `json:"equip_location"`
	EquipType     string `json:"equip_type"`
	Inspected     bool   `json:"inspected"` // 这个字段将通过 SQL 查询动态填充，代表是否已点检
}

// TableName 指定表名
func (EquipmentInspection) TableName() string {
	return "equipment_inspection"
}
