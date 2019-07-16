package lexicon

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func wordPrefix(word string) (string, bool) {
	if len(word) < 3 {
		return "", false
	}

	return word[:3], true
}

type WordPrefixGroup struct {
	Prefix string                     // Character(s) in this set (redundant in use as the previous map will have this same prefix)
	IsWord bool                       // Bool if the prefix to this letter is also a word
	Words  map[string]WordPrefixGroup // Map of the next list of words
}

type Lexicon struct {
	Words map[string]WordPrefixGroup
}

func (l Lexicon) String() string {
	return fmt.Sprintln(l.Words)
}

/**
 * LoadLexicon reads in a lexicon from a filename and holds it in memory while the current solver is running.
 * Makes a tree map of the words so that when searching, we can tell the moment we deviate from
 * possible words. Starts with a key of the first 3 letters since that is our minimum word length
 *
 * @example
 * test, tests, testing
 *
 * false      true     false    false    true
 * [ tes ] -- [ t ] -- [ i ] -- [ n ] -- [ g ]
 *                  \
 *                   \ true
 *                     [ s ]
 *
 * @param  string  filename  file name and location of the lexicon to read and load
 */
func (l *Lexicon) LoadLexicon(filename string) {
	l.Words = make(map[string]WordPrefixGroup)
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

		wordPrefix, ok := l.Words[prefix]
		if ok == false {
			isWord := prefix == word
			wordPrefix = WordPrefixGroup{Prefix: prefix, IsWord: isWord, Words: make(map[string]WordPrefixGroup)}
		}

		group, _ := createWordGroups(wordPrefix, word, len(prefix))
		wordPrefix.Words[group.Prefix] = group

		l.Words[prefix] = wordPrefix
	}

	if err := scanner.Err(); err != nil {
		// This might also need to end the app
	}
}

func createWordGroups(parentGroup WordPrefixGroup, word string, index int) (WordPrefixGroup, error) {
	if index == len(word) {
		return WordPrefixGroup{}, errors.New("Out of characters")
	}
	char := string(word[index])
	group, ok := parentGroup.Words[char]
	if ok == false {
		group = WordPrefixGroup{Prefix: char, Words: make(map[string]WordPrefixGroup)}
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

/**
 * Check if the provided word is in the lexicon. Traverses through the lexicon map.
 * @param  stringn  word  The word to be checked against.
 * @return bool     true if the word was found in the lexicon
 * @return error    an error if the word goes beyond a possble map. This is important to
 *                  the ability to discontinue looking for longer words efficiently.
 */
func (l *Lexicon) CheckWord(word string) (bool, error) {
	prefix, ok := wordPrefix(word)
	if ok == false {
		return false, errors.New("Word not Long enough for a prefix")
	}

	wordPrefix, ok := l.Words[prefix]
	if ok == false {
		return false, errors.New("Prefix Not Found")
	}

	return recursiveCheckWord(wordPrefix, word, len(prefix))
}

func recursiveCheckWord(group WordPrefixGroup, word string, index int) (bool, error) {
	if len(word) == index {
		return group.IsWord, nil
	}

	nextGroup, ok := group.Words[string(word[index])]
	if ok == false {
		return false, errors.New("Deviated from mapped lexicon words")
	}

	return recursiveCheckWord(nextGroup, word, index+1)
}
