package core

import (
	"testing"
)

// Representative positions for benchmarks
var (
	fenStart     = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	fenMiddle    = "r1bqkb1r/pppp1ppp/2n2n2/4p2Q/2B1P3/8/PPPP1PPP/RNB1K1NR w KQkq - 4 3"
	fenKiwipete  = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	fenEndgame   = "8/8/4k3/8/8/4K3/4P3/8 w - - 0 1"
)

func BenchmarkCalculateAllActions_Start(b *testing.B) {
	g, _ := NewGameFromFEN(fenStart)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.calculateAllActions()
	}
}

func BenchmarkCalculateAllActions_Kiwipete(b *testing.B) {
	g, _ := NewGameFromFEN(fenKiwipete)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.calculateAllActions()
	}
}

func BenchmarkDoAction_Start(b *testing.B) {
	g, _ := NewGameFromFEN(fenStart)
	// Pick the first non-resign action
	action := g.Actions[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.DoAction(action)
	}
}

func BenchmarkDoAction_Kiwipete(b *testing.B) {
	g, _ := NewGameFromFEN(fenKiwipete)
	action := g.Actions[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.DoAction(action)
	}
}

func BenchmarkNewGameFromFEN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewGameFromFEN(fenKiwipete)
	}
}

func BenchmarkPerft3_Start(b *testing.B) {
	g, _ := NewGameFromFEN(fenStart)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		perft(g, 3)
	}
}

func perft(g Game, depth int) int {
	if depth == 0 {
		return 1
	}
	nodes := 0
	for _, action := range g.Actions {
		if action.IsResign || action.IsDraw {
			continue
		}
		nodes += perft(g.DoAction(action), depth-1)
	}
	return nodes
}
