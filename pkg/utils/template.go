package utils

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/config"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := config.TPL.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
	}
}
