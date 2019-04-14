package config

import "html/template"

var TPL *template.Template

func init() {
	// be aware of duplicated template names, this is done this way for simplicity
	templateList := []string{
		"templates/add.gohtml",
		"templates/list.gohtml",
	}
	TPL = template.Must(template.ParseFiles(templateList...))
}
