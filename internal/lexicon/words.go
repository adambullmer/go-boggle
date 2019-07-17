package lexicon

import (
	"errors"
)

// A WordPrefixGroup is a group of word prefix nodes
type WordPrefixGroup map[string]WordPrefix

/*
A WordPrefix is a node in the lexicon tree. This node can represent the last letter of a word,
and can link to many more nodes in the tree. Due to the nature of
*/
type WordPrefix struct {
	Prefix string          // Character(s) in this set (redundant in use as the previous map will have this same prefix)
	IsWord bool            // Bool if the prefix to this letter is also a word
	Words  WordPrefixGroup // Map of the next list of words
}

// NewWordPrefix is a factory function for creating new word prefix nodes
func NewWordPrefix(prefix string, isWord bool) WordPrefix {
	return WordPrefix{
		Prefix: prefix,
		IsWord: isWord,
		Words:  make(WordPrefixGroup),
	}
}

// PrefixLength is a package constant for consistent splitting on the word prefix optimization
const PrefixLength = 3

// SplitWord is a helper function for splitting words up by the prefix and their remainder
func SplitWord(word string) (string, string) {
	if len(word) < PrefixLength {
		return "", ""
	}

	return word[:PrefixLength], word[PrefixLength:]
}

func wordPrefix(word string) (string, bool) {
	if len(word) < 3 {
		return "", false
	}

	return word[:3], true
}

/*
CheckWord checks if the provided word exists in the lexicon.
The current implementation does a Breadth-first, like search since the word
*/
func CheckWord(lexicon WordPrefixGroup, word string) (bool, error) {
	prefix, ok := wordPrefix(word)
	if ok == false {
		return false, errors.New("Word not Long enough for a prefix")
	}

	wordPrefix, ok := lexicon[prefix]
	if ok == false {
		return false, errors.New("Prefix Not Found")
	}

	return recursiveCheckWord(wordPrefix, word, len(prefix))
}

func recursiveCheckWord(group WordPrefix, word string, index int) (bool, error) {
	if len(word) == index {
		return group.IsWord, nil
	}

	nextGroup, ok := group.Words[string(word[index])]
	if ok == false {
		return false, errors.New("Deviated from mapped lexicon words")
	}

	return recursiveCheckWord(nextGroup, word, index+1)
}
