package listing

import "image/color"

type RiskMatrix struct {
	ID              int
	Path            string
	Project         string
	MatImgWidth     int
	MatImgHeight    int
	MatNrRows       int
	MatNrCols       int
	MatSize         int
	BorderWidth     int
	Multiple        int
	WordHeight      int
	WordWidth       int
	HighRiskColor   *color.RGBA
	MediumRiskColor *color.RGBA
	LowRiskColor    *color.RGBA
	RiskLabelColor  *color.RGBA
	BorderColor     *color.RGBA
}
