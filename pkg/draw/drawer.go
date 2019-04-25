package draw

import (
	"errors"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func RiskMatrixDrawer(filename string, m listing.RiskMatrix, risks []adding.Risk) error {
	myImg := image.NewRGBA(image.Rect(0, 0, m.MatImgWidth, m.MatImgHeight))
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "media", filename)
	outputFile, err := os.Create(path)
	if err != nil {
		return errors.New("we couldn't create the base risk matrix")
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	riskColor := func(riskBlock int) color.RGBA {
		if riskBlock == 1 || riskBlock == 5 || riskBlock == 9 {
			return m.MediumRiskColor
		} else if riskBlock == 2 || riskBlock == 3 || riskBlock == 6 {
			return m.HighRiskColor
		} else if riskBlock == 4 || riskBlock == 7 || riskBlock == 8 {
			return m.LowRiskColor
		}
		return color.RGBA{}
	}

	for blockNbr := 1; blockNbr <= m.MatSize; blockNbr++ {

		if blockNbr == 1 {
			r := image.Rect(0, 0, m.Multiple, m.Multiple)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
			wordSpacer := 2
			for _, r := range risks {
				if r.Probability == 3 && r.Impact == 1 {
					addLabel(myImg, r.Name, m.BorderWidth+m.WordHeight+wordSpacer, m.BorderWidth+m.WordHeight+wordSpacer, m.RiskLabelColor)
				}
				wordSpacer += 2
			}
		}
		if blockNbr == 2 {
			r := image.Rect(m.Multiple, 0, m.Multiple*2, m.Multiple)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 3 {
			r := image.Rect(m.Multiple*2, 0, m.Multiple*3, m.Multiple)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 4 {
			r := image.Rect(0, m.Multiple, m.Multiple, m.Multiple*2)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 5 {
			r := image.Rect(m.Multiple, m.Multiple, m.Multiple*2, m.Multiple*2)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 6 {
			r := image.Rect(m.Multiple*2, m.Multiple, m.Multiple*3, m.Multiple*2)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 7 {
			r := image.Rect(0, m.Multiple*2, m.Multiple, m.Multiple*3)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 8 {
			r := image.Rect(m.Multiple, m.Multiple*2, m.Multiple*2, m.Multiple*3)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
		if blockNbr == 9 {
			r := image.Rect(m.Multiple*2, m.Multiple*2, m.Multiple*3, m.Multiple*3)
			draw.Draw(myImg, r, &image.Uniform{C: riskColor(blockNbr)}, image.ZP, draw.Src)
		}
	}

	drawRiskMatrixBorders(myImg, m)

	// create image
	if err = png.Encode(outputFile, myImg); err != nil {
		return errors.New("we couldn't save the base risk matrix")
	}

	return nil
}

func drawRiskMatrixBorders(im *image.RGBA, m listing.RiskMatrix) {
	// draw borders
	topBorder := image.Rect(0, 0, m.MatImgWidth, 3)
	rightBorder := image.Rect(m.MatImgWidth-m.BorderWidth, 0, m.MatImgWidth, m.MatImgHeight)
	bottomBorder := image.Rect(0, m.MatImgHeight-m.BorderWidth, m.MatImgWidth, m.MatImgHeight)
	leftBorder := image.Rect(0, 0, m.BorderWidth, m.MatImgHeight)
	v1Border := image.Rect(m.MatImgWidth/m.MatNrCols, 0, (m.MatImgHeight/m.MatNrRows)-m.BorderWidth, m.MatImgHeight)
	v2Border := image.Rect((m.MatImgWidth/m.MatNrCols)*2, 0, ((m.MatImgWidth/m.MatNrCols)*2)-m.BorderWidth, m.MatImgHeight)
	h1Border := image.Rect(0, m.MatImgHeight/m.MatNrRows, m.MatImgWidth, (m.MatImgHeight/m.MatNrRows)-m.BorderWidth)
	h2Border := image.Rect(0, (m.MatImgHeight/m.MatNrRows)*2, m.MatImgWidth, ((m.MatImgHeight/m.MatNrRows)*2)-m.BorderWidth)

	// draw borders inside image
	draw.Draw(im, topBorder, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, rightBorder, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, bottomBorder, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, leftBorder, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, v1Border, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, v2Border, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, h1Border, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
	draw.Draw(im, h2Border, &image.Uniform{C: m.BorderColor}, image.ZP, draw.Src)
}

func addLabel(img *image.RGBA, label string, x, y int, c color.RGBA) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
