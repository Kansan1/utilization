package models

type Equipment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Qty  int    `json:"qty"`
	Line string `json:"line"`
}
