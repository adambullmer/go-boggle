package gameboard

import (
	"bytes"
	"errors"
	"sort"
	"text/template"

	"github.com/adambullmer/go-boggle/internal/lexicon"
	log "github.com/sirupsen/logrus"
)

// The GameBoard type holds the state of all the characters on the gamebaord
type GameBoard struct {
	Width  int
	Height int
	Board  [][]Cell
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
func (g *GameBoard) CheckBoard(l lexicon.WordPrefixGroup) Groups {
	words := make(Words)

	for y, row := range g.Board {
		for x := range row {
			wip := Letters{}
			g.checkNeighbors(l, wip, &words, x, y)
		}
	}

	var keys []int
	groups := make(Groups)

	for k := range words {
		index := len(k)
		if _, ok := groups[index]; ok == false {
			groups[index] = make([]string, 0)
			keys = append(keys, index)
		}
		groups[index] = append(groups[index], k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		if key < 3 {
			continue
		}

		words := groups[key]

		sort.Strings(words)
		groups[key] = words
	}

	return groups
}

func (g *GameBoard) checkNeighbors(l lexicon.WordPrefixGroup, wip Letters, words *Words, posX int, posY int) error {
	char := g.Board[posY][posX]
	if char.InUse == true {
		return nil
	}

	wip = append(wip, char.Character)

	if len(wip) >= 3 {
		// check dictionary
		word := wip.String()
		valid, err := lexicon.CheckWord(l, word)

		// early return if prefix isn't in lexicon
		if err != nil {
			wip = wip[:len(wip)-1]
			return errors.New("")
		}

		if valid == true {
			(*words)[word] = struct{}{}
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
		g.checkNeighbors(l, wip, words, posXright, posY)
	}

	if pos2 {
		g.checkNeighbors(l, wip, words, posXright, posYdown)
	}

	if pos3 {
		g.checkNeighbors(l, wip, words, posX, posYdown)
	}

	if pos4 {
		g.checkNeighbors(l, wip, words, posXleft, posYdown)
	}

	if pos5 {
		g.checkNeighbors(l, wip, words, posXleft, posY)
	}

	if pos6 {
		g.checkNeighbors(l, wip, words, posXleft, posYup)
	}

	if pos7 {
		g.checkNeighbors(l, wip, words, posX, posYup)
	}

	if pos8 {
		g.checkNeighbors(l, wip, words, posXright, posYup)
	}

	char.InUse = false
	g.Board[posY][posX] = char
	if len(wip) > 0 {
		wip = wip[:len(wip)-1]
	}

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
