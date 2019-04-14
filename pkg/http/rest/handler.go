package rest

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/config"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"github.com/ottotech/riskmanagement/pkg/utils"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type App struct {
	List *List
	Add  *Add
	Get  *Get
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

func (h *Add) Handler(s adding.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := config.TPL.ExecuteTemplate(w, "add.gohtml", nil)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}

		p := r.PostFormValue("project")
		rm := adding.RiskMatrix{Project: p}
		s.AddRiskMatrix(rm)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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

			rm, err := s.GetRiskMatrix(id)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			utils.RenderTemplate(w, "detail.gohtml", rm)
			return
		}

	})
}
