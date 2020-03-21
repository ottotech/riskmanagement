package utils

import (
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/assets"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	sb1, err := assets.Asset(name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
	}
	tpl, err := template.New("").Parse(string(sb1))
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	sb2, err := assets.Asset("templates/shutdown.gohtml")
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	tpl, err = tpl.Parse(string(sb2))
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
