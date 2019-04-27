package memory

import (
	"errors"
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"image/color"
	"os"
	"path/filepath"
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
	red    = color.RGBA{R: 0xff, A: 0xff}                   // rgb(255, 0, 0) high risk
	yellow = color.RGBA{R: 0xff, G: 0xff, A: 0xff}          // rgb(255, 255, 0) medium risk
	green  = color.RGBA{R: 0x90, G: 0xee, B: 0x90, A: 0xff} // rgb(144, 238, 144) low risk
	white  = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // rgb(255, 255, 255) border color
	black  = color.RGBA{A: 0xff}                            // rgb(0, 0, 0) label color
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

// UpdateRiskMatrixSize updates the risk matrix size of a given risk matrix in the repository
func (m *Storage) UpdateRiskMatrixSize(riskMatrixID, newImageWidth int ) error {
	for i := range m.riskMatrixSlice {
		if m.riskMatrixSlice[i].ID == riskMatrixID {
			m.riskMatrixSlice[i].MatImgWidth = newImageWidth
			m.riskMatrixSlice[i].MatImgHeight = newImageWidth
			m.riskMatrixSlice[i].Multiple = newImageWidth / m.riskMatrixSlice[i].MatNrCols
			return nil
		}
	}
	return errors.New(fmt.Sprintf("risk matrix not found with ID: %v", riskMatrixID))
}

// AddRisk saves the given risk in the repository
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
			ID:             id,
			RiskMatrixID:   r.RiskMatrixID,
			Name:           r.Name,
			Probability:    r.Probability,
			Impact:         r.Impact,
			Classification: r.Classification,
			Strategy:       r.Strategy,
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

// GetRiskMatrixByPath returns a risk matrix with the specified image path
func (m *Storage) GetRiskMatrixByPath(p string) (listing.RiskMatrix, error) {
	var riskMatrix listing.RiskMatrix

	for i := range m.riskMatrixSlice {

		if m.riskMatrixSlice[i].Path == p {
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

	return riskMatrix, errors.New(fmt.Sprintf("risk matrix not found by the given path: %v", p))
}

// GetAllRiskMatrix returns all the risk matrix stored in the database
func (m *Storage) GetAllRiskMatrix() []listing.RiskMatrix {
	var list []listing.RiskMatrix

	for i := range m.riskMatrixSlice {
		riskMatrix := listing.RiskMatrix{
			ID:              m.riskMatrixSlice[i].ID,
			Path:            m.riskMatrixSlice[i].Path,
			Project:         m.riskMatrixSlice[i].Project,
			MatImgWidth:     m.riskMatrixSlice[i].MatImgWidth,
			MatImgHeight:    m.riskMatrixSlice[i].MatImgHeight,
			MatNrRows:       m.riskMatrixSlice[i].MatNrRows,
			MatNrCols:       m.riskMatrixSlice[i].MatNrCols,
			MatSize:         m.riskMatrixSlice[i].MatSize,
			BorderWidth:     m.riskMatrixSlice[i].BorderWidth,
			Multiple:        m.riskMatrixSlice[i].Multiple,
			WordHeight:      m.riskMatrixSlice[i].WordHeight,
			WordWidth:       m.riskMatrixSlice[i].WordWidth,
			HighRiskColor:   m.riskMatrixSlice[i].HighRiskColor,
			MediumRiskColor: m.riskMatrixSlice[i].MediumRiskColor,
			LowRiskColor:    m.riskMatrixSlice[i].LowRiskColor,
			RiskLabelColor:  m.riskMatrixSlice[i].RiskLabelColor,
			BorderColor:     m.riskMatrixSlice[i].BorderColor,
		}
		list = append(list, riskMatrix)
	}
	return list
}

// GetAll returns all the risks for a given risk matrix
func (m *Storage) GetAllRisks(riskMatrixID int) []listing.Risk {
	var list []listing.Risk

	for i := range m.risks {
		if m.risks[i].RiskMatrixID == riskMatrixID {
			r := listing.Risk{
				ID:             m.risks[i].ID,
				RiskMatrixID:   m.risks[i].RiskMatrixID,
				Name:           m.risks[i].Name,
				Probability:    m.risks[i].Probability,
				Impact:         m.risks[i].Impact,
				Classification: m.risks[i].Classification,
				Strategy:       m.risks[i].Strategy,
			}

			list = append(list, r)
		}
	}
	return list
}

// GetRisk returns a risk with the given ID
func (m *Storage) GetRisk(riskID string) (listing.Risk, error) {
	for i := range m.risks {
		if m.risks[i].ID == riskID {
			r := listing.Risk{
				ID:             m.risks[i].ID,
				RiskMatrixID:   m.risks[i].RiskMatrixID,
				Name:           m.risks[i].Name,
				Probability:    m.risks[i].Probability,
				Impact:         m.risks[i].Impact,
				Classification: m.risks[i].Classification,
				Strategy:       m.risks[i].Strategy,
			}
			return r, nil
		}
	}
	return listing.Risk{}, errors.New(fmt.Sprintf("risk not found by the given ID: %v.", riskID))
}

// DeleteRisk deletes a risk with the specified ID
func (m *Storage) DeleteRisk(riskID string) error {
	for i := range m.risks {
		if m.risks[i].ID == riskID {
			m.risks[i] = m.risks[len(m.risks)-1]
			m.risks = m.risks[:len(m.risks)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("risk not found by the given ID: %v.", riskID))
}

// DeleteMatrix deletes a risk matrix with the specified ID
func (m *Storage) DeleteRiskMatrix(riskMatrixID int) error {
	for i := range m.riskMatrixSlice {
		if m.riskMatrixSlice[i].ID == riskMatrixID {
			// we remove the RiskMatrix image
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			path := filepath.Join(wd, "media", m.riskMatrixSlice[i].Path)
			_ = os.Remove(path)
			// we remove the data of the matrix stored in memory
			m.riskMatrixSlice[i] = m.riskMatrixSlice[len(m.riskMatrixSlice)-1]
			m.riskMatrixSlice = m.riskMatrixSlice[:len(m.riskMatrixSlice)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("risk matrix not found by the give ID: %v.", riskMatrixID))

}
