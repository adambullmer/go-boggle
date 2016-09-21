package GameBoard

import (
    "errors"
    "fmt"
    "log"
    "sort"

    "../Lexicon"
)

type GameBoard struct {
    Width int
    Height int
    Board [][]Cell
    ValidWords *Lexicon.Lexicon
}

func (g *GameBoard) OutputBoard() {
    for y := range g.Board {
        str := "[ "
        for x := range g.Board[y] {
            str += g.Board[y][x].String() + " "
        }
        str += "]"
        fmt.Println(str)
    }
}

func (g *GameBoard) String() string {
    str := "[ "
    for y := range g.Board {
        str += "[ "
        for x := range g.Board[y] {
            str += g.Board[y][x].String() + " "
        }
        str += "] "
    }

    str += "]"

    return str
}

func (g *GameBoard) Init(boardMap []string, lexLocation string) {
    g.Board = make([][]Cell, g.Height)
    for x := range g.Board {
        g.Board[x] = make([]Cell, g.Width)
    }

    g.SetBoard(boardMap)

    lexicon := new(Lexicon.Lexicon)
    lexicon.LoadLexicon(lexLocation)
    g.ValidWords = lexicon
}

func (g *GameBoard) SetBoard(boardMap []string) {
    for x := 0; x < len(boardMap); x++ {
        posX := x % g.Width
        posY := x / g.Height
        g.Board[posY][posX] = Cell{ Character: boardMap[x] }
    }
}

func (g *GameBoard) CheckBoard() []string {
    wip := new(WordInProgress)

    for y, row := range g.Board {
        for x, _ := range row {
            g.CheckNeighbors(wip, x, y)
        }
    }

    // Use map to record duplicates as we find them.
    encountered := map[string]bool{}
    result := []string{}

    for v := range wip.Words {
        if encountered[wip.Words[v]] == true {
        } else {
            encountered[wip.Words[v]] = true
            result = append(result, wip.Words[v])
        }
    }

    sort.Strings(result)
    log.Println(result)

    return result
}

func (g *GameBoard) CheckNeighbors(wip *WordInProgress, posX int, posY int) error {
    char := g.Board[posY][posX]
    if char.InUse == true {
        return nil
    }

    wip.Letters = wip.Push(char.Character)

    if len(wip.Letters) >= 3 {
        // check dictionary
        word := wip.Compress()
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
    posYup    := posY - 1
    posYdown  := posY + 1
    posXleft  := posX - 1
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
        g.CheckNeighbors(wip, posXright, posY)
    }

    if pos2 {
        g.CheckNeighbors(wip, posXright, posYdown)
    }

    if pos3 {
        g.CheckNeighbors(wip, posX, posYdown)
    }

    if pos4 {
        g.CheckNeighbors(wip, posXleft, posYdown)
    }

    if pos5 {
        g.CheckNeighbors(wip, posXleft, posY)
    }

    if pos6 {
        g.CheckNeighbors(wip, posXleft, posYup)
    }

    if pos7 {
        g.CheckNeighbors(wip, posX, posYup)
    }

    if pos8 {
        g.CheckNeighbors(wip, posXright, posYup)
    }

    char.InUse = false
    g.Board[posY][posX] = char
    wip.Pop()

    return nil
}

// Game Board Individual Cell
type Cell struct {
    Character string
    InUse bool
}

func (c *Cell) String() string {
    return c.Character
}

// Word In Progress
type WordInProgress struct {
    Letters []string
    Words []string
}

func (w *WordInProgress) Compress() string {
    word := ""
    for _, letter := range w.Letters {
        word += letter
    }

    return word
}

func (w *WordInProgress) Push(letter string) []string {
    w.Letters = append(w.Letters, letter)
    return w.Letters
}

func (w *WordInProgress) Pop() string {
    index     := len(w.Letters) - 1
    letter    := w.Letters[index]
    w.Letters  = w.Letters[:index]
    return letter
}

func (w *WordInProgress) AddWord(word string) {
    w.Words = append(w.Words, word)
}
