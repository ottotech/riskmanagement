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
	ListHandler *ListHandler
	AddHandler  *AddHandler
}

type ListHandler struct {
}

func (h *ListHandler) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		list := s.GetAllRiskMatrix()

		err := config.TPL.ExecuteTemplate(w, "list.gohtml", list)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	})
}


type AddHandler struct {
}

func (h *AddHandler) Handler(s adding.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := config.TPL.ExecuteTemplate(w, "add.gohtml", nil)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}
	})
}

