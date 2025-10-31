package models

type EquipmentPlan struct {
	ID          int `json:"id"`
	EquipmentID int `json:"equipment_id"`
	Day         int `json:"day"`
}
