package memory

type Risk struct {
	ID             string
	RiskMatrixID   int
	Name           string
	Probability    int
	Impact         int
	Classification string
	Strategy       string
}
