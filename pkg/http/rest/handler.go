package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/deleting"
	"github.com/ottotech/riskmanagement/pkg/draw"
	"github.com/ottotech/riskmanagement/pkg/listing"
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
		err := draw.RiskMatrixDrawer(filename, newRm)
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

func (h *DeleteRiskMatrix) Handler(d deleting.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get riskID from request
		id, _ := strconv.Atoi(r.PostFormValue("risk_matrix_id"))

		// delete the risk
		_ = d.DeleteRiskMatrix(id)

		// if all goes well we return response 200
		w.WriteHeader(http.StatusOK)

	})
}

type AddRisk struct {
}

func (h *AddRisk) Handler(a adding.Service, l listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get all risks from request
		var risks []adding.Risk
		data := r.PostFormValue("data")
		err := json.Unmarshal([]byte(data), &risks)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// remove risks that already exist
		riskMatrixID := risks[0].RiskMatrixID
		preexistingRisks := l.GetAllRisks(riskMatrixID)
		for i := 0; i < len(risks); i++ {
			for j := 0; j < len(preexistingRisks); j++ {
				if risks[i].Name == preexistingRisks[j].Name {
					risks[i] = risks[len(risks)-1]
					risks = risks[:len(risks)-1]
					i--
				}
			}
		}

		// if there are no new risks we send a response with status code 200
		if len(risks) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		// adding risks
		_ = a.AddRisk(risks...)

		// get risk matrix
		riskMatrix, _ := l.GetRiskMatrix(riskMatrixID)

		// draw risk matrix again
		_ = draw.RiskMatrixDrawer(riskMatrix.Path, riskMatrix)

		// if all goes well we send a status 200
		w.WriteHeader(http.StatusOK)
	})
}

type DeleteRisk struct {
}

func (h *DeleteRisk) Handler(d deleting.Service) http.Handler {
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

		// delete the risk
		_ = d.DeleteRisk(id)

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
