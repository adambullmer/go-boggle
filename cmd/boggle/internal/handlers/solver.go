package handlers

import (
	"net/http"

	"github.com/adambullmer/go-boggle/internal/gameboard"
	"github.com/adambullmer/go-boggle/internal/lexicon"
	"github.com/adambullmer/go-boggle/internal/platform/web"
	log "github.com/sirupsen/logrus"
)

func solverHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// Handle error
	}

	tmpl, err := web.RenderTemplate(w, r, "results.html")
	if err != nil {
		return
	}

	boardMap := r.Form["board[]"]

	gameBoard := gameboard.NewGameBoard(boardMap)
	lexicon, _ := lexicon.NewLexicon("./dictionaries/sowpods.txt")
	log.Println(gameBoard)
	words := gameBoard.CheckBoard(lexicon)
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
