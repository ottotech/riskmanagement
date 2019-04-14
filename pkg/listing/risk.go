package listing

type Risk struct {
	ID           string
	RiskMatrixID int
	Name         string
	Probability  int
	Impact       int
	Strategy     string
	ResponsePlan string
}
