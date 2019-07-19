package gameboard

import (
	"bytes"
	"errors"
	"sort"
	"text/template"

	"github.com/adambullmer/go-boggle/internal/lexicon"
	"github.com/adambullmer/go-boggle/internal/rules"
	log "github.com/sirupsen/logrus"
)

// The GameBoard type holds the state of all the characters on the gamebaord
type GameBoard struct {
	Width      int
	Height     int
	Board      [][]Cell
	ValidWords lexicon.WordPrefixGroup
}

/*
NewGameBoard is a factory for the gameboard data. It will leave the board in its
zero state to keep the factory lean and efficient.
*/
func NewGameBoard(boardMap []string) GameBoard {
	const height = 5
	const width = 5

	board := make([][]Cell, height)
	for x, char := range boardMap {
		posX := x % width
		posY := x / height
		if x%height == 0 {
			board[posY] = make([]Cell, width)
		}
		board[posY][posX] = Cell{Character: char}
	}

	// Not populating t
	return GameBoard{
		Height: height,
		Width:  width,
		Board:  board,
	}
}

const strTemplate = `
┌───────────────────────────────────────────────┐
{{- range .Board }}
│						│
│	{{ range . }}{{ . }}	{{end}}│
│						│
{{- end }}
└───────────────────────────────────────────────┘
`

func (g GameBoard) String() string {
	t, err := template.New("Gameboard String").Parse(strTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return ""
	}
	var str bytes.Buffer
	if err := t.Execute(&str, g); err != nil {
		log.Fatal("Execute: ", err)
		return ""
	}

	return str.String()
}

/*
CheckBoard starts the recursive breadth-first search on the gameboard.
This one is the main bit. Iterates over the entire board, looking for words.
Words can be strung together from any direction, so long as a letter has not
been used more than once in the same word.
*/
func (g *GameBoard) CheckBoard(l lexicon.WordPrefixGroup) map[int][]string {
	wip := new(rules.WordInProgress)
	wip.Words = make(map[int][]string)

	for y, row := range g.Board {
		for x := range row {
			g.checkNeighbors(l, wip, x, y)
		}
	}

	keys := []int{}
	for r := range wip.Words {
		keys = append(keys, r)
	}
	sort.Ints(keys)

	wordLists := new(rules.WordInProgress)
	wordLists.Words = make(map[int][]string)
	for _, key := range keys {
		if key < 3 {
			continue
		}

		words := wip.Words[key]

		sort.Strings(words)
		wordLists.Words[key] = words
	}

	return wordLists.Words
}

func (g *GameBoard) checkNeighbors(l lexicon.WordPrefixGroup, wip *rules.WordInProgress, posX int, posY int) error {
	char := g.Board[posY][posX]
	if char.InUse == true {
		return nil
	}

	wip.Letters = wip.Push(char.Character)

	if len(wip.Letters) >= 3 {
		// check dictionary
		word := wip.String()
		valid, err := lexicon.CheckWord(l, word)

		// early return if prefix isn't in lexicon
		if err != nil {
			wip.Pop()
			return errors.New("")
		}

		if valid == true {
			wip.AddWord(word)
		}
	}

	// Mark this letter as in use
	char.InUse = true
	g.Board[posY][posX] = char

	maxY := len(g.Board) - 1
	maxX := len(g.Board[0]) - 1
	posYup := posY - 1
	posYdown := posY + 1
	posXleft := posX - 1
	posXright := posX + 1

	/**
	 * Order to check surrounding spaces
	 *
	 * 6 7 8
	 * 5   1
	 * 4 3 2
	 */

	pos1 := posXright <= maxX
	pos3 := posYdown <= maxY
	pos5 := posXleft >= 0
	pos7 := posYup >= 0

	pos2 := pos1 && pos3
	pos4 := pos3 && pos5
	pos6 := pos5 && pos7
	pos8 := pos7 && pos1

	if pos1 {
		g.checkNeighbors(l, wip, posXright, posY)
	}

	if pos2 {
		g.checkNeighbors(l, wip, posXright, posYdown)
	}

	if pos3 {
		g.checkNeighbors(l, wip, posX, posYdown)
	}

	if pos4 {
		g.checkNeighbors(l, wip, posXleft, posYdown)
	}

	if pos5 {
		g.checkNeighbors(l, wip, posXleft, posY)
	}

	if pos6 {
		g.checkNeighbors(l, wip, posXleft, posYup)
	}

	if pos7 {
		g.checkNeighbors(l, wip, posX, posYup)
	}

	if pos8 {
		g.checkNeighbors(l, wip, posXright, posYup)
	}

	char.InUse = false
	g.Board[posY][posX] = char
	wip.Pop()

	return nil
}

// Cell is a Game Board individual space
type Cell struct {
	Character string
	InUse     bool
}

func (c Cell) String() string {
	return c.Character
}
