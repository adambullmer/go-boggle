package main

import (
    "errors"
    "html/template"
    "log"
    "net/http"
    "os"
    "path"

    "./GameBoard"
)

var AppName = "Online Boggle Solver"


func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string) (*template.Template, error) {
    layout := path.Join("templates", "layout.html")
    page   := path.Join("templates", templateName)

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

    tmpl, err := template.ParseFiles(layout, page)

    if err != nil {
        // Log the detailed error
        log.Println(err.Error())
        // Return a generic "Internal Server Error" message
        http.Error(w, http.StatusText(500), 500)
        return new(template.Template), errors.New("Error parsing templates")
    }

    return tmpl, nil
}


func viewHandler(w http.ResponseWriter, r *http.Request) {
    pageName := r.URL.Path
    if pageName == "" || pageName[:1] == "/" {
        pageName += "index"
    }
    pageName += ".html"

    tmpl, err := renderTemplate(w, r, pageName)
    if err != nil {
        return
    }

    if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
        log.Println(err.Error())
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

    tmpl, err := renderTemplate(w, r, "results.html")
    if err != nil {
        return
    }

    boardMap := r.Form["board[]"]
    boardWidth := 5
    boardHeight := 5

    gameBoard := GameBoard.GameBoard{Height: boardHeight, Width: boardWidth}
    gameBoard.Init(boardMap, "./dictionary.txt")
    gameBoard.OutputBoard()
    words := gameBoard.CheckBoard()

    data := struct {
        Words []string
    } {
        words,
    }

    if err := tmpl.ExecuteTemplate(w, "layout", &data); err != nil {
        log.Println(err.Error())
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

    log.Println("Listening... (Press ctrl + c to quit)")
    http.ListenAndServe(":8080", nil)
}
