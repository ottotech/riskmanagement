package rest

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
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
	List  *List
	Add   *Add
	Get   *Get
	Media *Media
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

		p := r.PostFormValue("project")
		if p == "" {
			utils.RenderTemplate(w, "add.gohtml", "Error: You need to specify the project name.")
			return
		}

		t := time.Now().Format("02_01_2006_03_04_05")
		filename := fmt.Sprintf("%v.png", t)
		rm := adding.RiskMatrix{Project: p, Path: filename}
		a.AddRiskMatrix(rm)
		newRm, _ := l.GetRiskMatrixByPath(filename)
		err := draw.DrawRiskMatrix(filename, newRm)
		if err != nil {
			utils.RenderTemplate(w, "add.gohtml", err.Error())
			return
		}

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
