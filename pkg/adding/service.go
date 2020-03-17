package adding

import (
	"github.com/ottotech/riskmanagement/pkg/config"
)

// Service provides risk matrix and risk adding operations
type Service interface {
	AddRiskMatrix(...RiskMatrix) error
	AddRisk(...Risk) error
	SaveMediaPath(path string) error
}

// Repository provides access to RiskMatrix repository.
type Repository interface {
	AddRiskMatrix(RiskMatrix) error
	AddRisk(Risk) error
	SaveMediaPath(path string) error
}

type service struct {
	rmR Repository //rmR = risk matrix repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddRiskMatrix can add the given risk matrix to the database
func (s *service) AddRiskMatrix(rm ...RiskMatrix) error {
	for _, matrix := range rm {
		err := s.rmR.AddRiskMatrix(matrix)
		if err != nil {
			config.Logger.Println(err)
			return err
		}
	}
	return nil
}

// AddRisk can add the given risks to the database
func (s *service) AddRisk(r ...Risk) error {
	for _, risk := range r {
		err := s.rmR.AddRisk(risk)
		if err != nil {
			config.Logger.Println(err)
			return err
		}
	}
	return nil
}

// SaveMediaPath saves the media path were all the matrix will be stored
func (s *service) SaveMediaPath(path string) error {
	err := s.rmR.SaveMediaPath(path)
	return err
}
