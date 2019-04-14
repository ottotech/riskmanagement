package rest

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/config"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"log"
	"net/http"
)

type App struct {
	HomeHandler *HomeHandler
	AddHandler  *AddHandler
	ListHandler *ListHandler
}

type HomeHandler struct {
}

func (h *HomeHandler) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		err := config.TPL.ExecuteTemplate(w, "list.gohtml", nil)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	})
}

type CreateRiskMatrixHandler struct {
}

func (h *CreateRiskMatrixHandler) Handler(s adding.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodGet {
			err := config.TPL.ExecuteTemplate(w, "index.gohtml", nil)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}
	})
}

type AddHandler struct {
}

func (h *AddHandler) Handler(s adding.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.AddRiskMatrix()
	})
}

type ListHandler struct {
}

func (h *ListHandler) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
