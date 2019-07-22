package handlers

import (
	"net/http"

	"github.com/adambullmer/go-boggle/internal/platform/web"
	log "github.com/sirupsen/logrus"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	pageName := r.URL.Path
	if pageName == "" || pageName[:1] == "/" {
		pageName += "index"
	}
	pageName += ".html"

	tmpl, err := web.RenderTemplate(w, r, pageName)
	if err != nil {
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Info(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func routes() {
	// Static file serving
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/assets/", fileServer)

	http.HandleFunc("/solve", solverHandler)
	http.HandleFunc("/", viewHandler)

	http.HandleFunc("/healthcheck", health)
	http.HandleFunc("/health", health)
}
