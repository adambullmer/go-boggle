package lexicon

import (
	"bufio"
	"io"
	"os"
	"strings"
)

/*
NewLexicon reads in a lexicon from a file and holds it in memory while the current solver is running.
It makes a tree map of the words so that when searching, it is easily possible to tell the moment input deviates
from possible words. It Starts with a key of the first 3 letters since that is the minimum lengthfor valid boggle
words.

On Efficiency:
The worst cases of both memory and CPU depend on the length of the longest word added to the map. For the purposes
of Boggle, the worst case word can only be 25 characters long. And Since the first 3 letters are in the root nodes,
The maximum is 2 less than the length of the word.

Memory Efficiency:
Theoreritical O(26^(N - 2)) where N is the length of the longest word in the dictionary.
Actual        O(26^(25 - 2))

CPU Efficiency:
Theoreritical: O((N-2)) where N is the length of the longest word in the dictionary.
Actual         O(25 - 2)

DIAGRAM
Input: test, tests, testy, testing

                    false    false    true
                    [ i ] -- [ n ] -- [ g ]
                  /
false      true  /  true
[ tes ] -- [ t ] -- [ s ]
                 \
                  \ true
                    [ y ]
*/
func NewLexicon(filename string) (WordPrefixGroup, error) {
	lexicon := make(WordPrefixGroup)

	h := func(word string) {
		prefix, rest := SplitWord(word)
		if prefix == "" {
			// ignore words that are too short
			return
		}

		wordPrefix, ok := lexicon[prefix]
		if ok == false {
			isWord := prefix == word
			wordPrefix = NewWordPrefix(prefix, isWord)
			lexicon[prefix] = wordPrefix
		}

		for i, c := range strings.Split(rest, "") {
			isWord := i+1 == len(rest)

			nextWord, ok := wordPrefix.Words[c]
			if ok == false {
				nextWord = NewWordPrefix(c, isWord)
				wordPrefix.Words[c] = nextWord
			}
			wordPrefix = nextWord
		}
	}

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		// This should probably end the app
		return nil, err
	}

	return lexicon, read(file, h)
}

func read(file *os.File, h func(word string)) error {
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024), 1024)
	for scanner.Scan() {
		// Each line is a new word to be added into the lexicon
		word := scanner.Text()
		h(word)
	}

	return scanner.Err()
}
