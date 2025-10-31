package models

type UtilizationException struct {
	Time             string `db:"time"`
	Line             string `db:"line"`
	RepairPerson     string `db:"repair_person"`
	RepairTime       string `db:"repair_time"`
	ExceptionDesc    string `db:"exception_desc"`
	AbnormalStopTime string `db:"abnormal_stop_time"`
	Remark           string `db:"remark"`
	MaintenanceTime  string `db:"maintenance_time"`
}
