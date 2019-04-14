package rest

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"io"
	"net/http"
)

type App struct {
	HomeHandler *HomeHandler
	AddHandler  *AddHandler
}

type HomeHandler struct {
}

func (h *HomeHandler) Handler(s listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("run")
		fmt.Println(s.GetRiskMatrix(1))
		io.WriteString(w, "Holaa")
	})
}

type AddHandler struct {
}

func (h *AddHandler) Handler(s adding.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
