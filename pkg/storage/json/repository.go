package json

import (
	"encoding/json"
	"errors"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"image/color"
	"path"
	"runtime"
	"strconv"
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

const (
	// dir defines the name of the directory where the files are stored
	dir = "/data/"
	// CollectionRisk identifier for the JSON collection of risks
	CollectionRisk = "risks"
	// CollectionMatrix identifier for the JSON collection of matrix
	CollectionMatrix = "matrix"
)

// Memory storage keeps data in memory
type Storage struct {
	db *scribble.Driver
}

// NewStorage returns a new JSON  storage
func NewStorage() (*Storage, error) {
	var err error
	s := new(Storage)
	_, filename, _, _ := runtime.Caller(0)
	p := path.Dir(filename)
	s.db, err = scribble.New(p+dir, nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// AddRiskMatrix saves the given risk matrix in repository
func (s *Storage) AddRiskMatrix(rm adding.RiskMatrix) error {
	existingRiskMatrix := s.GetAllRiskMatrix()
	newRM := RiskMatrix{
		ID:              len(existingRiskMatrix) + 1,
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
	resource := strconv.Itoa(newRM.ID)
	if err := s.db.Write(CollectionMatrix, resource, newRM); err != nil {
		return err
	}
	return nil
}

// UpdateRiskMatrixSize updates the risk matrix size of a given risk matrix in the repository
func (s *Storage) UpdateRiskMatrixSize(riskMatrixID, newSize int) error {
	var riskMatrix RiskMatrix
	err := s.db.Read(CollectionMatrix, strconv.Itoa(riskMatrixID), riskMatrix)
	if err != nil {
		return err
	}
	riskMatrix.MatImgWidth = newSize
	riskMatrix.MatImgHeight = newSize
	riskMatrix.Multiple = newSize / riskMatrix.MatNrCols
	resource := strconv.Itoa(riskMatrix.ID)
	if err := s.db.Write(CollectionMatrix, resource, riskMatrix); err != nil {
		return err
	}
	return nil
}

// AddRisk saves the given risk in the repository
func (s *Storage) AddRisk(r adding.Risk) error {
	var riskMatrix RiskMatrix
	err := s.db.Read(CollectionMatrix, strconv.Itoa(r.RiskMatrixID), riskMatrix)
	if err != nil {
		return err
	}
	dateCreated := time.Now()
	id := fmt.Sprintf("%d_%d", r.RiskMatrixID, dateCreated.Unix())
	newR := Risk{
		ID:             id,
		RiskMatrixID:   r.RiskMatrixID,
		Name:           r.Name,
		Probability:    r.Probability,
		Impact:         r.Impact,
		Classification: r.Classification,
		Strategy:       r.Strategy,
	}
	resource := newR.ID
	if err := s.db.Write(CollectionRisk, resource, newR); err != nil {
		return err
	}
	return nil
}

// GetRiskMatrix returns a risk matrix with the specified ID
func (s *Storage) GetRiskMatrix(id int) (listing.RiskMatrix, error) {
	var riskMatrix RiskMatrix
	var listingRiskMatrix listing.RiskMatrix
	err := s.db.Read(CollectionMatrix, strconv.Itoa(id), riskMatrix)
	if err != nil {
		return listingRiskMatrix, err
	}
	listingRiskMatrix.ID = riskMatrix.ID
	listingRiskMatrix.Path = riskMatrix.Path
	listingRiskMatrix.Project = riskMatrix.Project
	listingRiskMatrix.MatImgWidth = riskMatrix.MatImgWidth
	listingRiskMatrix.MatImgHeight = riskMatrix.MatImgHeight
	listingRiskMatrix.MatNrRows = riskMatrix.MatNrRows
	listingRiskMatrix.MatNrCols = riskMatrix.MatNrCols
	listingRiskMatrix.MatSize = riskMatrix.MatSize
	listingRiskMatrix.BorderWidth = riskMatrix.BorderWidth
	listingRiskMatrix.Multiple = riskMatrix.Multiple
	listingRiskMatrix.WordHeight = riskMatrix.WordHeight
	listingRiskMatrix.WordWidth = riskMatrix.WordWidth
	listingRiskMatrix.HighRiskColor = riskMatrix.HighRiskColor
	listingRiskMatrix.MediumRiskColor = riskMatrix.MediumRiskColor
	listingRiskMatrix.LowRiskColor = riskMatrix.LowRiskColor
	listingRiskMatrix.RiskLabelColor = riskMatrix.RiskLabelColor
	listingRiskMatrix.BorderColor = riskMatrix.BorderColor
	return listingRiskMatrix, nil
}

// GetRiskMatrixByPath returns a risk matrix with the specified image path
func (s *Storage) GetRiskMatrixByPath(p string) (listing.RiskMatrix, error) {
	var riskMatrix listing.RiskMatrix
	records := s.GetAllRiskMatrix()
	for i := range records {
		if records[i].Path == p {
			riskMatrix.ID = records[i].ID
			riskMatrix.Path = records[i].Path
			riskMatrix.Project = records[i].Project
			riskMatrix.MatImgWidth = records[i].MatImgWidth
			riskMatrix.MatImgHeight = records[i].MatImgHeight
			riskMatrix.MatNrRows = records[i].MatNrRows
			riskMatrix.MatNrCols = records[i].MatNrCols
			riskMatrix.MatSize = records[i].MatSize
			riskMatrix.BorderWidth = records[i].BorderWidth
			riskMatrix.Multiple = records[i].Multiple
			riskMatrix.WordHeight = records[i].WordHeight
			riskMatrix.WordWidth = records[i].WordWidth
			riskMatrix.HighRiskColor = records[i].HighRiskColor
			riskMatrix.MediumRiskColor = records[i].MediumRiskColor
			riskMatrix.LowRiskColor = records[i].LowRiskColor
			riskMatrix.RiskLabelColor = records[i].RiskLabelColor
			riskMatrix.BorderColor = records[i].BorderColor
			return riskMatrix, nil
		}
	}
	return riskMatrix, errors.New(fmt.Sprintf("risk matrix not found by the given path: %v", p))
}

// GetAllRiskMatrix returns all the risk matrix stored in the database
func (s *Storage) GetAllRiskMatrix() []listing.RiskMatrix {
	var list []listing.RiskMatrix
	records, err := s.db.ReadAll(CollectionMatrix)
	if err != nil {
		return list
	}
	for _, r := range records {
		var riskMatrix RiskMatrix
		var listingRiskMatrix listing.RiskMatrix
		if err := json.Unmarshal([]byte(r), &riskMatrix); err != nil {
			// err handling omitted for simplicity
			return list
		}
		listingRiskMatrix.ID = riskMatrix.ID
		listingRiskMatrix.Path = riskMatrix.Path
		listingRiskMatrix.Project = riskMatrix.Project
		listingRiskMatrix.MatImgWidth = riskMatrix.MatImgWidth
		listingRiskMatrix.MatImgHeight = riskMatrix.MatImgHeight
		listingRiskMatrix.MatNrRows = riskMatrix.MatNrRows
		listingRiskMatrix.MatNrCols = riskMatrix.MatNrCols
		listingRiskMatrix.MatSize = riskMatrix.MatSize
		listingRiskMatrix.BorderWidth = riskMatrix.BorderWidth
		listingRiskMatrix.Multiple = riskMatrix.Multiple
		listingRiskMatrix.WordHeight = riskMatrix.WordHeight
		listingRiskMatrix.WordWidth = riskMatrix.WordWidth
		listingRiskMatrix.HighRiskColor = riskMatrix.HighRiskColor
		listingRiskMatrix.MediumRiskColor = riskMatrix.MediumRiskColor
		listingRiskMatrix.LowRiskColor = riskMatrix.LowRiskColor
		listingRiskMatrix.RiskLabelColor = riskMatrix.RiskLabelColor
		listingRiskMatrix.BorderColor = riskMatrix.BorderColor
		list = append(list, listingRiskMatrix)
	}
	return list
}

// GetAllRisks returns all the risks for a given risk matrix
func (s *Storage) GetAllRisks(riskMatrixID int) []listing.Risk {
	var list []listing.Risk
	records, err := s.db.ReadAll(CollectionRisk)
	if err != nil {
		return list
	}
	for _, r := range records {
		var risk Risk
		var listingRisk listing.Risk
		if err := json.Unmarshal([]byte(r), &risk); err != nil {
			// err handling omitted for simplicity
			return list
		}
		listingRisk.ID = risk.ID
		listingRisk.RiskMatrixID = risk.RiskMatrixID
		listingRisk.Name = risk.Name
		listingRisk.Probability = risk.Probability
		listingRisk.Impact = risk.Impact
		listingRisk.Classification = risk.Classification
		listingRisk.Strategy = risk.Strategy
		list = append(list, listingRisk)
	}
	return list
}

// GetRisk returns a risk with the given ID
func (s *Storage) GetRisk(riskID string) (listing.Risk, error) {
	var risk Risk
	var listingRisk listing.Risk
	err := s.db.Read(CollectionRisk, riskID, risk)
	if err != nil {
		return listingRisk, err
	}
	listingRisk.ID = risk.ID
	listingRisk.RiskMatrixID = risk.RiskMatrixID
	listingRisk.Name = risk.Name
	listingRisk.Probability = risk.Probability
	listingRisk.Impact = risk.Impact
	listingRisk.Classification = risk.Classification
	listingRisk.Strategy = risk.Strategy
	return listingRisk, nil
}

// DeleteRisk deletes a risk with the specified ID
func (s *Storage) DeleteRisk(riskID string) error {
	err := s.db.Delete(CollectionRisk, riskID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMatrix deletes a risk matrix with the specified ID
func (s *Storage) DeleteRiskMatrix(riskMatrixID int) error {
	var err error
	risks := s.GetAllRisks(riskMatrixID)
	for _, r := range risks {
		err = s.DeleteRisk(r.ID)
		if err != nil {
			return err
		}
	}
	err = s.db.Delete(CollectionMatrix, strconv.Itoa(riskMatrixID))
	if err != nil {
		return err
	}
	return nil
}
