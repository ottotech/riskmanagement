package adding

type Risk struct {
	RiskMatrixID   int    `json:"risk_matrix_id"`
	Name           string `json:"name"`
	Probability    int    `json:"probability"`
	Impact         int    `json:"impact"`
	Classification string `json:"classification"`
	Strategy       string `json:"strategy"`
}
