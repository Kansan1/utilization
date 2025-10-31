package models

type EquipmentTypeMaintainStatVO struct {
	EquipmentType  string `json:"equipment_type"`  // 设备种类名称
	TotalCount     int    `json:"total_count"`     // 维护任务总数
	CompletedCount int    `json:"completed_count"` // 已维护完成任务数
}
