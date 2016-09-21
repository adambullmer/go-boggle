package Lexicon

import (
    "bufio"
    "errors"
    "fmt"
    "os"
)

func WordPrefix(word string) (string, bool) {
    if len(word) < 3 {
        return "", false
    }

    return word[:3], true
}

type WordPrefixGroup struct {
    Prefix string                    // Character(s) in this set (redundant in use as the previous map will have this same prefix)
    IsWord bool                      // Bool if the prefix to this letter is also a word
    Words map[string]WordPrefixGroup  // Map of the next list of words
}


type Lexicon struct {
    Words map[string]WordPrefixGroup
}

func (l *Lexicon) String() string {
    return fmt.Sprintln(l.Words)
}

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

        prefix, ok := WordPrefix(word)
        if ok == false {
            continue
        }

        wordPrefix, ok := l.Words[prefix]
        if ok {
        } else {
            isWord := prefix == word
            wordPrefix = WordPrefixGroup{Prefix: prefix, IsWord: isWord, Words: make(map[string]WordPrefixGroup)}
        }

        group, _ := CreateWordGroups(wordPrefix, word, len(prefix))
        wordPrefix.Words[group.Prefix] = group

        l.Words[prefix] = wordPrefix
    }

    if err := scanner.Err(); err != nil {
        // This might also need to end the app
    }
}

func CreateWordGroups(parentGroup WordPrefixGroup, word string, index int) (WordPrefixGroup, error)  {
    if index == len(word) {
        return WordPrefixGroup{}, errors.New("Out of characters")
    }
    char := string(word[index])
    group, ok := parentGroup.Words[char]
    if ok {
    } else {
        group = WordPrefixGroup{Prefix: char, Words: make(map[string]WordPrefixGroup)}
    }

    nextGroups, err := CreateWordGroups(group, word, index + 1)

    if err == nil {
        if _, ok := group.Words[nextGroups.Prefix]; ok {
        } else {
            group.Words[nextGroups.Prefix] = nextGroups
        }
    } else {
        group.IsWord = true
    }

    return group, nil
}


func (l *Lexicon) CheckWord(word string) (bool, error) {
    prefix, ok := WordPrefix(word)
    if ok == false {
        return false, errors.New("Word not Long enough for a prefix")
    }

    if wordPrefix, ok := l.Words[prefix]; ok {
        if len(word) == len(prefix) {
            if wordPrefix.IsWord == true {
                return true, nil
            }

            return false, nil
        }

        return recursiveCheckWord(wordPrefix, word, len(prefix))
    } else {
        return false, errors.New("Prefix Not Found")
    }

    // Word not found, no specific error necessary
    return false, nil
}

func recursiveCheckWord(group WordPrefixGroup, word string, index int) (bool, error) {
    // No more letters left
    if len(word) == index {
        return false, nil
    }

    if nextGroup, ok := group.Words[string(word[index])]; ok {
        isLastLetter := len(word) - 1 == index
        if isLastLetter && nextGroup.IsWord {
            return true, nil
        } else {
            return recursiveCheckWord(nextGroup, word, index + 1)
        }
    } else {
        return false, errors.New("Deviated from mapped lexicon words")
    }

    return false, nil
}
