package models

type DailyRepairTask struct {
	ID         int    `json:"id"`
	Line       string `json:"line"`
	DeviceName string `json:"device_name"`
	Fault      string `json:"fault"`
	ReportTime string `json:"report_time"` // 格式为字符串，例如 "2025-06-23 10:00:00"
	State      string `json:"state"`
	Repairer   string `json:"repairer"`
}
