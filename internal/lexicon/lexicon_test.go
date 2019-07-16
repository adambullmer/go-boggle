package lexicon

import (
	"testing"
)

func BenchmarkLexicon(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		lex := Lexicon{}
		lex.LoadLexicon("../../dictionaries/sowpods.txt")
	}
}
