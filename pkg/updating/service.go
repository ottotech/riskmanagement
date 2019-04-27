package updating

// Service provides risk matrix and risk updating operations
type Service interface {
	UpdateRiskMatrixSize(riskMatrixID, newImageWidth int) error
}

// Repository provides access to RiskMatrix repository
type Repository interface {
	UpdateRiskMatrixSize(riskMatrixID, newImageWidth int) error
}

type service struct {
	r Repository // rmR = risk matrix repository
}

// NewService creates an updating service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// UpdateRiskMatrixSize updates the risk matrix size of a given risk matrix
// with the specified ID in the repository
func (s *service) UpdateRiskMatrixSize(riskMatrixID, newImageWidth int) error {
	err := s.r.UpdateRiskMatrixSize(riskMatrixID, newImageWidth)
	return err
}