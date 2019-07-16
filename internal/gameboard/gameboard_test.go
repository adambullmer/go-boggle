package gameboard

import (
	"strings"
	"testing"
)

func BenchmarkCheckBoard(b *testing.B) {
	letters := strings.Split("a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,qu,r,s,t,v,w,x,y,z", ",")
	gameboard := GameBoard{Height: 5, Width: 5}
	gameboard.Init(letters, "../../dictionaries/sowpods.txt")
	for i := 0; i < b.N; i++ {
		gameboard.CheckBoard()
	}
}
