package main

import (
	"errors"
	"sort"

	log "github.com/sirupsen/logrus"
)

type GameBoard struct {
	Width      int
	Height     int
	Board      [][]Cell
	ValidWords *Lexicon
}

/**
 * Initializes the GameBoard. Triggers loading the dictionary, assigns the user input to the gameboard
 *
 * @param  []string   boardMap    Flat string of characters in the board.
 * @param  string     lexLocation which lexicon to be used when checking for words
 */
func (g *GameBoard) Init(boardMap []string, lexLocation string) {
	g.initBoard()
	g.setBoard(boardMap)

	lexicon := new(Lexicon)
	lexicon.LoadLexicon(lexLocation)
	g.ValidWords = lexicon
}

func (g *GameBoard) initBoard() {
	g.Board = make([][]Cell, g.Height)
	for x := range g.Board {
		g.Board[x] = make([]Cell, g.Width)
	}
}

/**
 * Overrides the default stringification of this struct.
 */
func (g GameBoard) String() string {
	str := "[\n"
	for y := range g.Board {
		str += "  [ "
		for x := range g.Board[y] {
			str += g.Board[y][x].String() + " "
		}
		str += "]\n"
	}

	str += "]"

	return str
}

func (g *GameBoard) setBoard(boardMap []string) {
	for x := 0; x < len(boardMap); x++ {
		posX := x % g.Width
		posY := x / g.Height
		g.Board[posY][posX] = Cell{Character: boardMap[x]}
	}
}

/**
 * This one is the main bit. Iterates over the entire board, looking for words.
 * Words can be strung together from any direction, so long as a letter has not
 * been used more than once in the same word.
 *
 * @returns  map[int][]string  a map of lists of words, sorted alphabetically,
 *                             grouped by word length.
 */
func (g *GameBoard) CheckBoard() map[int][]string {
	wip := new(WordInProgress)
	wip.Words = make(map[int][]string)

	for y, row := range g.Board {
		for x, _ := range row {
			g.checkNeighbors(wip, x, y)
		}
	}

	keys := []int{}
	for r := range wip.Words {
		keys = append(keys, r)
	}
	sort.Ints(keys)

	wordLists := new(WordInProgress)
	wordLists.Words = make(map[int][]string)
	for _, key := range keys {
		if key < 3 {
			continue
		}

		words := wip.Words[key]

		sort.Strings(words)
		wordLists.Words[key] = words
	}

	log.Info(wordLists.Words)

	return wordLists.Words
}

func (g *GameBoard) checkNeighbors(wip *WordInProgress, posX int, posY int) error {
	char := g.Board[posY][posX]
	if char.InUse == true {
		return nil
	}

	wip.Letters = wip.Push(char.Character)

	if len(wip.Letters) >= 3 {
		// check dictionary
		word := wip.String()
		valid, err := g.ValidWords.CheckWord(word)

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
		g.checkNeighbors(wip, posXright, posY)
	}

	if pos2 {
		g.checkNeighbors(wip, posXright, posYdown)
	}

	if pos3 {
		g.checkNeighbors(wip, posX, posYdown)
	}

	if pos4 {
		g.checkNeighbors(wip, posXleft, posYdown)
	}

	if pos5 {
		g.checkNeighbors(wip, posXleft, posY)
	}

	if pos6 {
		g.checkNeighbors(wip, posXleft, posYup)
	}

	if pos7 {
		g.checkNeighbors(wip, posX, posYup)
	}

	if pos8 {
		g.checkNeighbors(wip, posXright, posYup)
	}

	char.InUse = false
	g.Board[posY][posX] = char
	wip.Pop()

	return nil
}

// Game Board Individual Cell
type Cell struct {
	Character string
	InUse     bool
}

func (c Cell) String() string {
	return c.Character
}

// Word In Progress
type WordInProgress struct {
	Letters []string
	Words   map[int][]string
}

func (w WordInProgress) String() string {
	word := ""
	for _, letter := range w.Letters {
		word += letter
	}

	return word
}

/**
 * Pushes a single character to the working word
 *
 * @param  string  letter  letter to be added to the end of the working word
 *
 * @returns  []string  list of all the letters in the current word
 */
func (w *WordInProgress) Push(letter string) []string {
	w.Letters = append(w.Letters, letter)
	return w.Letters
}

/**
 * Pops the last character of the working word stack
 *
 * @returns  string  The letter that was removed from the working word
 */
func (w *WordInProgress) Pop() string {
	index := len(w.Letters) - 1
	letter := w.Letters[index]
	w.Letters = w.Letters[:index]
	return letter
}

/**
 * Adds the passed word to the map of total matched words
 *
 * @param  string  word  The word to to added to the map
 */
func (w *WordInProgress) AddWord(word string) {
	key, list := w.getWordKey(word)

	// dedupe list
	for _, existingWord := range list {
		if word == existingWord {
			return
		}
	}

	w.Words[key] = append(list, word)
}

func (w *WordInProgress) getWordKey(word string) (int, []string) {
	key := len(word)
	list, ok := w.Words[key]

	if ok == false {
		list = []string{}
	}

	return key, list
}
