package models

// Inspection 对应 Inspection 表
type Inspection struct {
	Time string `json:"time"`
	Type string `json:"type"` // 这里存储的是 equip_code
}

// TableName 指定 GORM 使用的表名
func (Inspection) TableName() string {
	return "Inspection"
}
