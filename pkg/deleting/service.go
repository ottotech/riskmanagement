package deleting

// Service provides risk matrix and risk deleting operations
type Service interface {
	DeleteRisk(...string) error
	DeleteRiskMatrix(...int) error
}

// Repository provide access to RiskMatrix repository
type Repository interface {
	DeleteRisk(string) error
	DeleteRiskMatrix(int) error
}

type service struct {
	r Repository // rmR = risk matrix repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// DeleteRisk deletes the risk with specified ID
func (s *service) DeleteRisk(riskIDs ...string) error {
	for _, id := range riskIDs {
		_ = s.r.DeleteRisk(id)
	}
	return nil
}

// DeleteRiskMatrix delete the risk matrix with the specified ID
func (s *service) DeleteRiskMatrix(riskMatrixIDs ...int) error {
	for _, id := range riskMatrixIDs {
		_ = s.r.DeleteRiskMatrix(id)
	}
	return nil
}
