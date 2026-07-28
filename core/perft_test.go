package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPerft(t *testing.T) {
	ts := []struct {
		name   string
		fen    string
		depths map[int]int
	}{
		{
			name: "Position 1: Initial position",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depths: map[int]int{
				0: 1,
				1: 20,
				2: 400,
				3: 8902,
				4: 197281,
			},
		},
		{
			name: "Position 2: Kiwipete",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depths: map[int]int{
				1: 48,
				2: 2039,
				3: 97862,
			},
		},
		{
			name: "Position 3: endgame with en passant and promotion",
			fen:  "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depths: map[int]int{
				1: 14,
				2: 191,
				3: 2812,
				4: 43238,
			},
		},
		{
			name: "Position 4: promotions and castling",
			fen:  "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depths: map[int]int{
				1: 6,
				2: 264,
				3: 9467,
				4: 422333,
			},
		},
		{
			name: "Position 5: promotion edge cases",
			fen:  "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			depths: map[int]int{
				1: 44,
				2: 1486,
				3: 62379,
			},
		},
		{
			name: "Position 6: mirrored evaluation",
			fen:  "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			depths: map[int]int{
				1: 46,
				2: 2079,
				3: 89890,
			},
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)

			for depth, expected := range tc.depths {
				t.Run(fmt.Sprintf("depth_%d", depth), func(t *testing.T) {
					if depth >= 4 && testing.Short() {
						t.Skipf("skipping depth %d in short mode", depth)
					}
					got := Perft(g, depth)
					require.Equal(t, expected, got, "depth %d: expected %d nodes, got %d", depth, expected, got)
				})
			}
		})
	}
}
