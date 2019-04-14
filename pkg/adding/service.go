package adding

// Service provides risk matrix and risk adding operations
type Service interface {
	AddRiskMatrix(...RiskMatrix)
	AddRisk(...Risk)
}

// Repository provides access to RiskMatrix repository.
type Repository interface {
	AddRiskMatrix(RiskMatrix) error
	AddRisk(Risk) error
}

type service struct {
	rmR Repository //rmR = risk matrix repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddRiskMatrix can add the give risk matrix to the database
func (s *service) AddRiskMatrix(rm ...RiskMatrix) {

	// Any validation can be done here

	for _, matrix := range rm {
		_ = s.rmR.AddRiskMatrix(matrix)
	}
}

// AddRiskMatrix can add the give risk matrix to the database
func (s *service) AddRisk(r ...Risk) {

	// Any validation can be done here

	for _, risk := range r {
		_ = s.rmR.AddRisk(risk)
	}
}
