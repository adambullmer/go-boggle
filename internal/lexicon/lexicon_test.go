package lexicon

import (
	"os"
	"testing"
)

const testDictionary = "../../dictionaries/test.txt"
const largeDictionary = "../../dictionaries/sowpods.txt"

func TestNewLexicon(t *testing.T) {
	lex, err := NewLexicon(testDictionary)

	if err != nil {
		t.Errorf("Unable to find dictionary %v", testDictionary)
		return
	}
	const expectedRootNodes = 1
	rootLength := len(lex)
	if rootLength != expectedRootNodes {
		t.Errorf("Total root nodes should be %d, got %d", expectedRootNodes, rootLength)
	}

	const expectedRootNode = "tes"
	rootNode, ok := lex[expectedRootNode]
	if ok != true {
		t.Errorf("Expected root node '%v' was not found", expectedRootNode)
	}
	if rootNode.IsWord == true {
		t.Errorf("Expected node not to be a valid word")
	}

	const expectedBaseNode = "t"
	baseNode, ok := rootNode.Words[expectedBaseNode]
	if ok != true {
		t.Errorf("Expected node '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode)
	}
	if baseNode.IsWord == false {
		t.Errorf("Expected node to be a valid word")
	}

	const expectedEndNode1 = "s"
	endNode, ok := baseNode.Words[expectedEndNode1]
	if ok != true {
		t.Errorf("Expected node  '%v' -- '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode, expectedEndNode1)
	}
	if endNode.IsWord == false {
		t.Errorf("Expected node to be a valid word")
	}

	const expectedEndNode2 = "y"
	endNode, ok = baseNode.Words[expectedEndNode2]
	if ok != true {
		t.Errorf("Expected node  '%v' -- '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode, expectedEndNode2)
	}
	if endNode.IsWord == false {
		t.Errorf("Expected node to be a valid word")
	}

	const expectedEndNode3 = "i"
	endNode, ok = baseNode.Words[expectedEndNode3]
	if ok != true {
		t.Errorf("Expected node  '%v' -- '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode, expectedEndNode3)
	}
	if endNode.IsWord == true {
		t.Errorf("Expected node not to be a valid word")
	}

	const expectedEndNode4 = "n"
	endNode, ok = endNode.Words[expectedEndNode4]
	if ok != true {
		t.Errorf("Expected node  '%v' -- '%v' -- '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode, expectedEndNode3, expectedEndNode4)
	}
	if endNode.IsWord == true {
		t.Errorf("Expected node not to be a valid word")
	}

	const expectedEndNode5 = "g"
	endNode, ok = endNode.Words[expectedEndNode5]
	if ok != true {
		t.Errorf("Expected node  '%v' -- '%v' -- '%v' -- '%v' -- '%v' was not found", expectedRootNode, expectedBaseNode, expectedEndNode3, expectedEndNode4, expectedEndNode5)
	}
	if endNode.IsWord == false {
		t.Errorf("Expected node to be a valid word")
	}
}

func TestLexiconNewLexiconUnknown(t *testing.T) {
	_, err := NewLexicon("./dictionaries/definitely-non-existent-dictionary.txt")
	if err == nil {
		t.Errorf("Expected to return an error on a missing dictionary")
	}
}

// ========================================================

func BenchmarkLexiconSmall(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewLexicon(testDictionary)
	}
}

func BenchmarkLexiconLarge(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		NewLexicon(largeDictionary)
	}
}

func BenchmarkReadSmall(b *testing.B) {
	const dictionary = testDictionary
	h := func(w []byte) {}
	file, _ := os.Open(dictionary)
	defer file.Close()
	for i := 0; i < b.N; i++ {
		read(file, h)
	}
}

func BenchmarkReadLarge(b *testing.B) {
	h := func(w []byte) {}
	file, _ := os.Open(largeDictionary)
	defer file.Close()
	for i := 0; i < b.N; i++ {
		read(file, h)
	}
}
