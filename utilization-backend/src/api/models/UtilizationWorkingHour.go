package models

type UtilizationWorkingHour struct {
	Time         string `db:"time"`
	Type         string `db:"type"`
	WorkingHours string `db:"working_hours"`
	Line         string `db:"line"`
}
