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

func BenchmarkCheckWordShort(b *testing.B) {
	lex := Lexicon{}
	lex.LoadLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		lex.CheckWord("sin")
	}
}

func BenchmarkCheckWordMedium(b *testing.B) {
	lex := Lexicon{}
	lex.LoadLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		lex.CheckWord("since")
	}
}

func BenchmarkCheckWordLong(b *testing.B) {
	lex := Lexicon{}
	lex.LoadLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		lex.CheckWord("science")
	}
}

func BenchmarkCheckWordExtraLong(b *testing.B) {
	lex := Lexicon{}
	lex.LoadLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		lex.CheckWord("scientifically")
	}
}

func BenchmarkCheckWordMissing(b *testing.B) {
	lex := Lexicon{}
	lex.LoadLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		lex.CheckWord("foobarbaz")
	}
}
