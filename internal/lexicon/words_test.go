package lexicon

import (
	"testing"
)

func BenchmarkCheckWordShort(b *testing.B) {
	lex, _ := NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		CheckWord(lex, "sin")
	}
}

func BenchmarkCheckWordMedium(b *testing.B) {
	lex, _ := NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		CheckWord(lex, "since")
	}
}

func BenchmarkCheckWordLong(b *testing.B) {
	lex, _ := NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		CheckWord(lex, "science")
	}
}

func BenchmarkCheckWordExtraLong(b *testing.B) {
	lex, _ := NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		CheckWord(lex, "scientifically")
	}
}

func BenchmarkCheckWordMissing(b *testing.B) {
	lex, _ := NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		CheckWord(lex, "foobarbaz")
	}
}
