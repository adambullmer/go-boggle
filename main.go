package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	pageName := r.URL.Path
	if pageName == "" || pageName[:1] == "/" {
		pageName += "index"
	}
	pageName += ".html"

	tmpl, err := RenderTemplate(w, r, pageName)
	if err != nil {
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Info(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

/**
 * @TODO: have a favicon
 *
 * @param  {http.ResponseWriter} w  web write stream
 * @param  {*http.Request} r        web response information
 */
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "")
}

func solverHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// Handle error
	}

	tmpl, err := RenderTemplate(w, r, "results.html")
	if err != nil {
		return
	}

	boardMap := r.Form["board[]"]
	boardWidth := 5
	boardHeight := 5

	gameBoard := GameBoard{Height: boardHeight, Width: boardWidth}
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

func main() {
	// Static file serving
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/assets/", fileServer)
	http.HandleFunc("/favicon.ico", faviconHandler)

	http.HandleFunc("/solve", solverHandler)
	http.HandleFunc("/", viewHandler)

	log.Info("Listening on http://localhost:8080 ... (Press ctrl + c to quit)")
	http.ListenAndServe(":8080", nil)
}
