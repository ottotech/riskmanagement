package listing

// Service provides risk matrix and risk listing operations
type Service interface {
	GetRiskMatrix(int) (RiskMatrix, error)
	GetRiskMatrixByPath(string) (RiskMatrix, error)
	GetAllRisks(int) []Risk
	GetRisk(string) (Risk, error)
	GetAllRiskMatrix() []RiskMatrix
	GetMediaPath() (string, error)
}

type Repository interface {
	// GetRiskMatrix returns the risk matrix with the given ID
	GetRiskMatrix(int) (RiskMatrix, error)
	// GetRiskMatrixByPath returns a risk matrix with the specified image path
	GetRiskMatrixByPath(string) (RiskMatrix, error)
	// GetAllRiskMatrix returns all the risk matrix stored
	GetAllRiskMatrix() []RiskMatrix
	// GetAllRisks returns a list of all risks for a given RiskMatrix ID
	GetAllRisks(int) []Risk
	// GetRisk returns a risk with the given ID
	GetRisk(string) (Risk, error)
	// GetMediaPath returns the path where the media is stored
	GetMediaPath() (string, error)
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetRiskMatrix returns a risk matrix with specified ID
func (s *service) GetRiskMatrix(id int) (RiskMatrix, error) {
	rm, err := s.r.GetRiskMatrix(id)
	return rm, err
}

// GetRiskMatrixByPath returns a risk matrix with the specified image path
func (s *service) GetRiskMatrixByPath(p string) (RiskMatrix, error) {
	rm, err := s.r.GetRiskMatrixByPath(p)
	return rm, err
}

// GetAllRiskMatrix returns all the risk matrix stored
func (s *service) GetAllRiskMatrix() []RiskMatrix {
	return s.r.GetAllRiskMatrix()
}

// GetAllRisks returns all risks specified in a RiskMatrix
func (s *service) GetAllRisks(riskMatrixID int) []Risk {
	return s.r.GetAllRisks(riskMatrixID)
}

// GetRisk returns a risk with the given ID
func (s *service) GetRisk(riskID string) (Risk, error) {
	r, err := s.r.GetRisk(riskID)
	return r, err
}

// GetMediaPath returns the path where the media is stored
func (s *service) GetMediaPath() (string, error) {
	path, err := s.r.GetMediaPath()
	return path, err
}
