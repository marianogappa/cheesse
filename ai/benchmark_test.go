package ai

import (
	"testing"

	"github.com/marianogappa/cheesse/core"
)

func BenchmarkAIDepth0_Start(b *testing.B) {
	g := core.NewDefaultGame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BasicAIAction(g, 0)
	}
}

func BenchmarkAIDepth1_Start(b *testing.B) {
	g := core.NewDefaultGame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BasicAIAction(g, 1)
	}
}

func BenchmarkAIDepth2_Endgame(b *testing.B) {
	g, _ := core.NewGameFromFEN("7k/5P2/8/8/3n4/1B6/8/K7 w - - 0 1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BasicAIAction(g, 2)
	}
}
