package gameboard

import (
	"strings"
	"testing"

	"github.com/adambullmer/go-boggle/internal/lexicon"
)

func BenchmarkCheckBoard(b *testing.B) {
	letters := strings.Split("a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,qu,r,s,t,v,w,x,y,z", ",")
	gameboard := NewGameBoard(letters)
	lexicon, _ := lexicon.NewLexicon("../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		gameboard.CheckBoard(lexicon)
	}
}
