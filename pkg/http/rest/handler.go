package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/deleting"
	"github.com/ottotech/riskmanagement/pkg/draw"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"github.com/ottotech/riskmanagement/pkg/updating"
	"github.com/ottotech/riskmanagement/pkg/utils"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type App struct {
	List             *List
	Add              *Add
	Get              *Get
	Media            *Media
	AddRisk          *AddRisk
	DeleteRisk       *DeleteRisk
	DeleteRiskMatrix *DeleteRiskMatrix
}

type List struct {
}

func (h *List) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		list := s.GetAllRiskMatrix()
		utils.RenderTemplate(w, "list.gohtml", list)
		return
	})
}

type Add struct {
}

func (h *Add) Handler(a adding.Service, l listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			utils.RenderTemplate(w, "add.gohtml", nil)
			return
		}
		name := r.PostFormValue("project")
		if name == "" {
			utils.RenderTemplate(w, "add.gohtml", "Error: You need to specify the project name.")
			return
		}
		t := time.Now().Format("02_01_2006_03_04_05")
		filename := fmt.Sprintf("%v.png", t)
		rm := adding.RiskMatrix{Project: name, Path: filename}
		_ = a.AddRiskMatrix(rm)
		newRm, _ := l.GetRiskMatrixByPath(filename)
		err := draw.RiskMatrixDrawer(filename, newRm, []adding.Risk{})
		if err != nil {
			utils.RenderTemplate(w, "add.gohtml", err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	})
}

type DeleteRiskMatrix struct {
}

func (h *DeleteRiskMatrix) Handler(d deleting.Service, l listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get riskID from request
		id, _ := strconv.Atoi(r.PostFormValue("risk_matrix_id"))

		// before deleting the matrix, let's instantiate it
		riskMatrix, _ := l.GetRiskMatrix(id)

		// delete the risk matrix
		_ = d.DeleteRiskMatrix(id)

		// delete all risks from matrix
		risks := l.GetAllRisks(riskMatrix.ID)
		risksIDs := make([]string, 0)
		for _, r := range risks {
			risksIDs = append(risksIDs, r.ID)
		}
		_ = d.DeleteRisk(risksIDs...)

		// if all goes well we return response 200
		w.WriteHeader(http.StatusOK)
	})
}

type AddRisk struct {
}

func (h *AddRisk) Handler(a adding.Service, l listing.Service, u updating.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get all risks from request
		var risks []adding.Risk
		postData := r.PostFormValue("data")
		err := json.Unmarshal([]byte(postData), &risks)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// anonymous func that returns the risk classification
		riskClassifier := func(r adding.Risk) (string, bool) {
			if r.Probability == 3 && r.Impact == 1 {
				return "medium", true
			}
			if r.Probability == 3 && r.Impact == 2 {
				return "high", true
			}
			if r.Probability == 3 && r.Impact == 3 {
				return "high", true
			}
			if r.Probability == 2 && r.Impact == 1 {
				return "low", true
			}
			if r.Probability == 2 && r.Impact == 2 {
				return "medium", true
			}
			if r.Probability == 2 && r.Impact == 3 {
				return "high", true
			}
			if r.Probability == 1 && r.Impact == 1 {
				return "low", true
			}
			if r.Probability == 1 && r.Impact == 2 {
				return "low", true
			}
			if r.Probability == 1 && r.Impact == 3 {
				return "medium", true
			}
			return "", false
		}

		// we need to set the risk classification for each risk
		for i, r := range risks {
			c, ok := riskClassifier(r)
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				msg := fmt.Sprintf("Probability and impact numbers should be a number between 1 and 0 on risk %v", r.Name)
				_, _ = w.Write([]byte(msg)) // ignoring error for simplicity
				return
			}
			risks[i].Classification = c
		}

		// classify new risks
		var newRisks []adding.Risk
		riskMatrixID := risks[0].RiskMatrixID
		preexistingRisks := l.GetAllRisks(riskMatrixID)
		for _, r := range risks {
			exists := false
			for _, pr := range preexistingRisks {
				if r.Name == pr.Name {
					exists = true
				}
			}
			if !exists {
				newRisks = append(newRisks, r)
			}
		}

		// if there are no new risks we send a response with status code 200
		if len(newRisks) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		// adding risks
		_ = a.AddRisk(newRisks...)

		// get risk matrix
		riskMatrix, _ := l.GetRiskMatrix(riskMatrixID)

		// this anonymous func will count how many risks are per block in the risk matrix.
		// then it will get the number of risks from the block that has the most number of risks.
		// finally based on this max number this func will resize the matrix in order to show the
		// risks accordingly.
		riskMatrixResize := func(rm listing.RiskMatrix, r []adding.Risk) {
			var rb1, rb2, rb3, rb4, rb5, rb6, rb7, rb8, rb9 int // rb = risk matrix block
			for _, r := range risks {
				if r.Probability == 3 && r.Impact == 1 {
					rb1 += 1
				}
				if r.Probability == 3 && r.Impact == 2 {
					rb2 += 1
				}
				if r.Probability == 3 && r.Impact == 3 {
					rb3 += 1
				}
				if r.Probability == 2 && r.Impact == 1 {
					rb4 += 1
				}
				if r.Probability == 2 && r.Impact == 2 {
					rb5 += 1
				}
				if r.Probability == 2 && r.Impact == 3 {
					rb6 += 1
				}
				if r.Probability == 1 && r.Impact == 1 {
					rb7 += 1
				}
				if r.Probability == 1 && r.Impact == 2 {
					rb8 += 1
				}
				if r.Probability == 1 && r.Impact == 3 {
					rb9 += 1
				}
			}
			var max int
			for _, v := range []int{rb1, rb2, rb3, rb4, rb5, rb6, rb7, rb8, rb9} {
				if v > max {
					max = v
				}
			}
			// block size:   200
			// border width: 3*2 (top and bottom)
			// word height:  13
			// line spacing: word height + 2
			blockWritableSize := 200 - 30 // 30 is a dummy value to add more space while adding risk labels on image
			actualSize := 0
			for i := 1; i < max; i++ {
				actualSize += 15
			}
			// if the size required is greater than the block writable size
			// we need to find the right size where the risks can fit.
			// note that the size of the blocks should be a reminder of a number
			// that is a multiple of 3 because the risk matrix has 3 rows and 3 columns
			if actualSize > blockWritableSize {
				actualSize = actualSize + 30 // 30 is a dummy value to add more space while adding risk labels on image
				newBlockWidth := 0
				multiple := riskMatrix.Multiple
				for {
					x := multiple + (3 - multiple%3)
					if x > actualSize {
						newBlockWidth = x
						break
					}
					multiple = x
				}
				newImgWidth := newBlockWidth * 3
				riskMatrix.MatImgWidth = newImgWidth
				riskMatrix.MatImgHeight = newImgWidth
				riskMatrix.Multiple = newImgWidth / riskMatrix.MatNrCols

				// finally we need to update the risk matrix sizes in the storage
				err := u.UpdateRiskMatrixSize(riskMatrixID, newImgWidth)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write([]byte(fmt.Sprintf("While updating the size of the risk "+
						"matrix we encounter this err: %v", err.Error())))
					return
				}
			}
		}

		// resize risk matrix image if necessary
		riskMatrixResize(riskMatrix, risks)

		// draw risk matrix again in order to add the new risks
		_ = draw.RiskMatrixDrawer(riskMatrix.Path, riskMatrix, risks)

		// if all goes well we send a status 200
		w.WriteHeader(http.StatusOK)
	})
}

type DeleteRisk struct {
}

func (h *DeleteRisk) Handler(d deleting.Service, l listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only allow POST method
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get riskID from request
		id := r.PostFormValue("risk_id")

		// we get the risk instance before deleting it to use it later
		risk, _ := l.GetRisk(id)

		// delete the risk
		_ = d.DeleteRisk(id)

		// get riskMatrix of the deleted risk
		riskMatrix, _ := l.GetRiskMatrix(risk.RiskMatrixID)

		// get all risks of the matrix
		var risks []adding.Risk
		for _, risk := range l.GetAllRisks(riskMatrix.ID) {
			r := adding.Risk{
				RiskMatrixID:   riskMatrix.ID,
				Name:           risk.Name,
				Probability:    risk.Probability,
				Impact:         risk.Impact,
				Classification: risk.Classification,
				Strategy:       risk.Strategy,
			}
			risks = append(risks, r)
		}

		// draw risk matrix again in order to not show the deleted risks
		_ = draw.RiskMatrixDrawer(riskMatrix.Path, riskMatrix, risks)

		// if all goes well we return response 200
		w.WriteHeader(http.StatusOK)
	})
}

type Get struct {
}

func (h *Get) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// clean path to get RiskMatrix ID
			p := path.Clean("/" + r.URL.Path)
			i := strings.Index(p[1:], "/") + 1
			tail := p[i+1:]

			// check if id is valid
			id, err := strconv.Atoi(tail)
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("This Risk Matrix ID is not valid: %v.", tail), http.StatusBadRequest)
				return
			}

			// get risk matrix
			riskMatrix, err := s.GetRiskMatrix(id)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// get all risks
			risks := s.GetAllRisks(riskMatrix.ID)

			// build context
			ctx := struct {
				RiskMatrix listing.RiskMatrix
				Risks      []listing.Risk
			}{
				riskMatrix,
				risks,
			}

			// render tpl
			utils.RenderTemplate(w, "detail.gohtml", ctx)
			return
		}
	})
}

type Media struct {
}

func (h *Media) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := "image/png"
		fullPath := r.URL.Path
		w.Header().Add("Content-Type", ct)

		f, err := os.Open(strings.TrimLeft(fullPath, "/"))
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		fi, err := f.Stat()
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, f.Name(), fi.ModTime(), f)
	})
}
