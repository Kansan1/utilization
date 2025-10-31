package models

type MonthlyFaultFrequency struct {
	ID            int    `json:"id"`
	Line          string `json:"line"`
	EquipmentName string `json:"equipment_name"`
	Content       string `json:"content"`
	FaultCount    int    `json:"fault_count"`
	State         string `json:"state"`
	Defender      string `json:"defender"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
}
