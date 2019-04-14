package listing

// Service provides risk matrix and risk listing operations
type Service interface {
	GetRiskMatrix(int) (RiskMatrix, error)
	GetAllRisks(int) []Risk
}

type Repository interface {
	// GetRiskMatrix returns the risk matrix with the given ID
	GetRiskMatrix(int) (RiskMatrix, error)
	// GetAllRisks returns a list of all risks for a given RiskMatrix ID
	GetAllRisks(int) []Risk
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetRiskMatrix returns a beer
func (s *service) GetRiskMatrix(id int) (RiskMatrix, error) {
	return s.r.GetRiskMatrix(id)
}

// GetAllRisks returns all risks specified in a RiskMatrix
func (s *service) GetAllRisks(riskMatrixID int) []Risk {
	return s.r.GetAllRisks(riskMatrixID)
}
