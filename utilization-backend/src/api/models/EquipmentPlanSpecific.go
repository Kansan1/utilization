package models

type EquipmentPlanSpecific struct {
	ID          int    `json:"id"`           // 自增主键
	Line        string `json:"line"`         // 产线
	EquipmentID int    `json:"equipment_id"` // 设备ID
	Content     string `json:"content"`      // 内容
	State       string `json:"state"`        // 状态（例如：待处理、已完成等）
	Defender    string `json:"defender"`     // 责任人
	Day         int    `json:"day"`
	Month       int    `json:"month"`
	Year        int    `json:"year"`
	Shu         int    `json:"shu"`
}
