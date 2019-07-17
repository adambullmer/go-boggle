package lexicon

import (
	"bufio"
	"io"
	"os"
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

	h := func(word []byte) {
		if len(word) < PrefixLength {
			return
		}
		prefix := string(word[:PrefixLength])
		rest := word[PrefixLength:]

		wordPrefix, ok := lexicon[prefix]
		if ok == false {
			isWord := len(word) == PrefixLength
			wordPrefix = NewWordPrefix(prefix, isWord)
			lexicon[prefix] = wordPrefix
		}

		for i, char := range rest {
			isWord := i+1 == len(rest)
			str := string(char)

			nextWord, ok := wordPrefix.Words[str]
			if ok == false {
				nextWord = NewWordPrefix(str, isWord)
				wordPrefix.Words[str] = nextWord
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

func read(file *os.File, h func(word []byte)) error {
	file.Seek(0, io.SeekStart)
	r := bufio.NewReader(file)
	var err error

	for {
		word, _, err := r.ReadLine()
		if err != nil {
			break
		}

		h(word)
	}

	if err == io.EOF {
		return nil
	}

	return err
}
