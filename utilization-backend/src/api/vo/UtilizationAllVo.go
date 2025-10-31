package vo

type UtilizationAllVo struct {
	Date        string  `json:"date"`
	Value       float32 `json:"value"`
	ValueAll    float32 `json:"valueAll"`
	ValueActual float32 `json:"valueActual"`
}
