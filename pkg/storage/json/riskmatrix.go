package json

import (
	"image/color"
	"time"
)

type RiskMatrix struct {
	ID              int        `json:"id"`
	Path            string     `json:"path"`
	Project         string     `json:"project"`
	DateCreated     time.Time  `json:"date_created"`
	MatImgWidth     int        `json:"mat_img_width"`
	MatImgHeight    int        `json:"mat_img_height"`
	MatNrRows       int        `json:"mat_nr_rows"`
	MatNrCols       int        `json:"mat_nr_cols"`
	MatSize         int        `json:"mat_size"`
	BorderWidth     int        `json:"border_width"`
	Multiple        int        `json:"multiple"`
	WordHeight      int        `json:"word_height"`
	WordWidth       int        `json:"word_width"`
	HighRiskColor   color.RGBA `json:"high_risk_color"`
	MediumRiskColor color.RGBA `json:"medium_risk_color"`
	LowRiskColor    color.RGBA `json:"low_risk_color"`
	RiskLabelColor  color.RGBA `json:"risk_label_color"`
	BorderColor     color.RGBA `json:"border_color"`
}
