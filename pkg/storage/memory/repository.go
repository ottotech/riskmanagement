package memory

import (
	"errors"
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"image/color"
	"time"
)

// matrix setup
const (
	imWidth      = 600
	imHeight     = imWidth
	matrixNrRows = 3
	matrixNrCols = 3
	matrixSize   = matrixNrRows * matrixNrCols
	multiple     = imWidth / matrixNrCols
	borderWidth  = 3
	wordWith     = 6
	wordHeight   = 13
)

// colors
var (
	red    = &color.RGBA{R: 0xff, A: 0xff}                   // rgb(255, 0, 0) high risk
	yellow = &color.RGBA{R: 0xff, G: 0xff, A: 0xff}          // rgb(255, 255, 0) medium risk
	green  = &color.RGBA{R: 0x90, G: 0xee, B: 0x90, A: 0xff} // rgb(144, 238, 144) low risk
	white  = &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // rgb(255, 255, 255) border color
	black  = &color.RGBA{A: 0xff}                            // rgb(0, 0, 0) label color
)

// Memory storage keeps data in memory
type Storage struct {
	riskMatrixSlice []RiskMatrix
	risks           []Risk
}

// Add saves the given risk matrix in repository
func (m *Storage) AddRiskMatrix(rm adding.RiskMatrix) error {
	newRM := RiskMatrix{
		ID:              len(m.riskMatrixSlice) + 1,
		Path:            rm.Path,
		Project:         rm.Project,
		MatImgWidth:     imWidth,
		MatImgHeight:    imHeight,
		MatNrRows:       matrixNrRows,
		MatNrCols:       matrixNrCols,
		MatSize:         matrixSize,
		BorderWidth:     borderWidth,
		Multiple:        multiple,
		WordHeight:      wordHeight,
		WordWidth:       wordWith,
		HighRiskColor:   red,
		MediumRiskColor: yellow,
		LowRiskColor:    green,
		RiskLabelColor:  black,
		BorderColor:     white,
	}

	m.riskMatrixSlice = append(m.riskMatrixSlice, newRM)
	return nil
}

// Add saves the given risk in the repository
func (m *Storage) AddRisk(r adding.Risk) error {
	found := false
	for rm := range m.riskMatrixSlice {
		if m.riskMatrixSlice[rm].ID == r.RiskMatrixID {
			found = true
		}
	}

	if found {
		created := time.Now()
		id := fmt.Sprintf("%d_%d", r.RiskMatrixID, created.Unix())
		newR := Risk{
			ID:           id,
			RiskMatrixID: r.RiskMatrixID,
			Name:         r.Name,
			Probability:  r.Probability,
			Impact:       r.Impact,
			Strategy:     r.Strategy,
			ResponsePlan: r.ResponsePlan,
		}
		m.risks = append(m.risks, newR)

	} else {
		return errors.New("risk matrix not found")
	}
	return nil
}

// Get returns a risk matrix with the specified ID
func (m *Storage) GetRiskMatrix(id int) (listing.RiskMatrix, error) {
	var riskMatrix listing.RiskMatrix

	for i := range m.riskMatrixSlice {

		if m.riskMatrixSlice[i].ID == id {
			riskMatrix.ID = m.riskMatrixSlice[i].ID
			riskMatrix.Path = m.riskMatrixSlice[i].Path
			riskMatrix.Project = m.riskMatrixSlice[i].Project
			riskMatrix.MatImgWidth = m.riskMatrixSlice[i].MatImgWidth
			riskMatrix.MatImgHeight = m.riskMatrixSlice[i].MatImgHeight
			riskMatrix.MatNrRows = m.riskMatrixSlice[i].MatNrRows
			riskMatrix.MatNrCols = m.riskMatrixSlice[i].MatNrCols
			riskMatrix.MatSize = m.riskMatrixSlice[i].MatSize
			riskMatrix.BorderWidth = m.riskMatrixSlice[i].BorderWidth
			riskMatrix.Multiple = m.riskMatrixSlice[i].Multiple
			riskMatrix.WordHeight = m.riskMatrixSlice[i].WordHeight
			riskMatrix.WordWidth = m.riskMatrixSlice[i].WordWidth
			riskMatrix.HighRiskColor = m.riskMatrixSlice[i].HighRiskColor
			riskMatrix.MediumRiskColor = m.riskMatrixSlice[i].MediumRiskColor
			riskMatrix.LowRiskColor = m.riskMatrixSlice[i].LowRiskColor
			riskMatrix.RiskLabelColor = m.riskMatrixSlice[i].RiskLabelColor
			riskMatrix.BorderColor = m.riskMatrixSlice[i].BorderColor
			return riskMatrix, nil
		}
	}

	return riskMatrix, errors.New("risk matrix not found")
}

// GetAll returns all the risks for a given risk matrix
func (m *Storage) GetAllRisks(riskMatrixID int) []listing.Risk {
	var list []listing.Risk

	for i := range m.risks {
		if m.risks[i].RiskMatrixID == riskMatrixID {
			r := listing.Risk{
				ID:           m.risks[i].ID,
				RiskMatrixID: m.risks[i].RiskMatrixID,
				Name:         m.risks[i].Name,
				Probability:  m.risks[i].Probability,
				Impact:       m.risks[i].Impact,
				Strategy:     m.risks[i].Strategy,
				ResponsePlan: m.risks[i].ResponsePlan,
			}

			list = append(list, r)
		}
	}
	return list
}