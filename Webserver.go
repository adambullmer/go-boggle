package main

import (
	"errors"
	"net/http"
	"os"
	"path"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var AppName = "Online Boggle Solver"

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string) (*template.Template, error) {
	layout := path.Join("templates", "layout.html")
	page := path.Join("templates", templateName)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(page)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return new(template.Template), errors.New("File does not exist")
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return new(template.Template), errors.New("TemplateName is a directory")
	}

	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"loop": func(n int) []struct{} {
			return make([]struct{}, n)
		},
	}).ParseFiles(layout, page)

	if err != nil {
		// Log the detailed error
		log.Info(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return new(template.Template), errors.New("Error parsing templates")
	}

	return tmpl, nil
}
