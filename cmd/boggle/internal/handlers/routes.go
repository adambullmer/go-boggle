package handlers

import (
	"net/http"

	"github.com/adambullmer/go-boggle/internal/gameboard"
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

func solverHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// Handle error
	}

	tmpl, err := web.RenderTemplate(w, r, "results.html")
	if err != nil {
		return
	}

	boardMap := r.Form["board[]"]
	boardWidth := 5
	boardHeight := 5

	gameBoard := gameboard.GameBoard{Height: boardHeight, Width: boardWidth}
	gameBoard.Init(boardMap, "./dictionaries/sowpods.txt")
	log.Println(gameBoard)
	words := gameBoard.CheckBoard()
	wordCount := 0

	for _, wordList := range words {
		wordCount += len(wordList)
	}

	data := struct {
		Words     map[int][]string
		WordCount int
	}{
		words,
		wordCount,
	}

	if err := tmpl.ExecuteTemplate(w, "layout", &data); err != nil {
		log.Info(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}