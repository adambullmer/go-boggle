package lexicon

import (
	"bufio"
	"errors"
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
	file, err := os.Open(filename)
	if err != nil {
		// This should probably end the app
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Each line is a new word to be added into the lexicon
		word := scanner.Text()

		prefix, ok := wordPrefix(word)
		if ok == false {
			continue
		}

		wordPrefix, ok := lexicon[prefix]
		if ok == false {
			isWord := prefix == word
			wordPrefix = NewWordPrefix(prefix, isWord)
		}

		group, _ := createWordGroups(wordPrefix, word, len(prefix))
		wordPrefix.Words[group.Prefix] = group

		lexicon[prefix] = wordPrefix
	}

	if err := scanner.Err(); err != nil {
		// This might also need to end the app
	}

	return lexicon, nil
}

func createWordGroups(parentGroup WordPrefix, word string, index int) (WordPrefix, error) {
	if index == len(word) {
		return WordPrefix{}, errors.New("Out of characters")
	}

	char := string(word[index])
	group, ok := parentGroup.Words[char]
	if ok == false {
		group = NewWordPrefix(char, false)
	}

	nextGroups, err := createWordGroups(group, word, index+1)

	if err == nil {
		if _, ok := group.Words[nextGroups.Prefix]; ok == false {
			group.Words[nextGroups.Prefix] = nextGroups
		}
	} else {
		group.IsWord = true
	}

	return group, nil
}
