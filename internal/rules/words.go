package rules

/*
WordInProgress is the data structure for both a word in progress,
and for tracking the collection of words found.
*/
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

/*
Push puts the provided letter onto the Letters stack
*/
func (w *WordInProgress) Push(letter string) []string {
	w.Letters = append(w.Letters, letter)
	return w.Letters
}

/*
Pop removes the last letter in the Letters stack
*/
func (w *WordInProgress) Pop() string {
	index := len(w.Letters) - 1
	letter := w.Letters[index]
	w.Letters = w.Letters[:index]
	return letter
}

/*
AddWord adds a word to the found words data structure
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
