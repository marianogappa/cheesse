package ai

import (
	"math/rand"
	"testing"
	"time"

	"github.com/marianogappa/cheesse/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const hardBudget = time.Second

// tacticalPuzzles is a fixed suite of positions that Hard must solve: the
// engine's move must match the expected best move (from-to in algebraic).
var tacticalPuzzles = []struct {
	name     string
	fen      string
	from, to string
}{
	{
		name: "back-rank mate in 1",
		fen:  "6k1/5ppp/8/8/8/8/8/3R2K1 w - - 0 1",
		from: "d1", to: "d8",
	},
	{
		name: "smothered mate in 1",
		fen:  "6rk/6pp/7N/8/8/8/8/R5K1 w - - 0 1",
		from: "h6", to: "f7",
	},
	{
		name: "royal fork wins the rook",
		fen:  "r3k3/8/8/1N6/8/8/8/6K1 w - - 0 1",
		from: "b5", to: "c7",
	},
	{
		name: "capture the undefended rook",
		fen:  "6k1/8/8/8/8/8/r5PP/1r2R1K1 w - - 0 1",
		from: "e1", to: "b1",
	},
	{
		name: "promote and win",
		fen:  "7k/5P2/8/8/8/8/8/1K6 w - - 0 1",
		from: "f7", to: "f8",
	},
}

func TestHard_SolvesTacticalSuite(t *testing.T) {
	for _, puzzle := range tacticalPuzzles {
		t.Run(puzzle.name, func(t *testing.T) {
			g, err := core.NewGameFromFEN(puzzle.fen)
			require.NoError(t, err)
			action, _, ok := TimedAIAction(g, hardBudget)
			require.True(t, ok)
			got := action.FromPiece.XY.ToAlgebraic() + action.ToXY.ToAlgebraic()
			assert.Equal(t, puzzle.from+puzzle.to, got,
				"expected %s%s, got %s", puzzle.from, puzzle.to, got)
		})
	}
}

func TestHard_FindsMateIn2(t *testing.T) {
	// White Qa3, Rh1, Kg1 vs lone Black king: 1.Qa7 Kf8 2.Rh8#.
	g, err := core.NewGameFromFEN("6k1/8/8/8/8/Q7/8/6KR w - - 0 1")
	require.NoError(t, err)

	for ply := 0; ply < 4 && !g.IsGameOver; ply++ {
		_, newGame, ok := TimedAIAction(g, hardBudget)
		require.True(t, ok)
		g = newGame
	}
	assert.True(t, g.IsCheckmate, "Hard should mate within 2 moves; final FEN: %v", g.ToFEN())
}

func TestHard_PrefersShorterMate(t *testing.T) {
	// White Qa3, Rh1, Kg1 vs lone Black king. Mate-distance scoring must make
	// the engine actually finish the game rather than shuffle indefinitely in a
	// trivially won position.
	g, err := core.NewGameFromFEN("6k1/8/8/8/8/Q7/8/6KR w - - 0 1")
	require.NoError(t, err)

	for ply := 0; ply < 20 && !g.IsGameOver; ply++ {
		_, newGame, ok := TimedAIAction(g, hardBudget)
		require.True(t, ok)
		g = newGame
	}
	assert.True(t, g.IsCheckmate, "Hard should convert Q+R vs K quickly; final FEN: %v", g.ToFEN())
}

func TestHard_ConvertsKQKEndgame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow endgame conversion test in -short mode")
	}
	// K+Q vs K is a textbook win; the engine must checkmate within 30 moves.
	g, err := core.NewGameFromFEN("8/8/4k3/8/8/3Q4/4K3/8 w - - 0 1")
	require.NoError(t, err)

	for ply := 0; ply < 60 && !g.IsGameOver; ply++ {
		_, newGame, ok := TimedAIAction(g, hardBudget)
		require.True(t, ok)
		g = newGame
	}
	assert.True(t, g.IsCheckmate, "Hard should convert KQ vs K; final FEN: %v", g.ToFEN())
}

func TestHard_BeatsRandom(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow Hard-vs-Random match in -short mode")
	}
	// Hard must win >= 95% of games vs Random. 20 games, alternating colors,
	// with a reduced budget to keep the test fast — strength vs Random saturates
	// well below the full Hard budget.
	const games = 20
	rng := rand.New(rand.NewSource(20260809))
	wins := 0

	for i := 0; i < games; i++ {
		g := core.NewDefaultGame()
		hardPlaysWhite := i%2 == 0
		for ply := 0; ply < 300 && !g.IsGameOver; ply++ {
			var ok bool
			if (g.Turn() == core.ColorWhite) == hardPlaysWhite {
				_, g, ok = TimedAIAction(g, 100*time.Millisecond)
			} else {
				_, g, ok = RandomAction(g, rng)
			}
			require.True(t, ok)
		}
		hardColor := core.Color(core.ColorBlack)
		if hardPlaysWhite {
			hardColor = core.ColorWhite
		}
		if g.IsCheckmate && g.GameOverWinner == hardColor {
			wins++
		}
	}

	winRate := float64(wins) / float64(games)
	assert.GreaterOrEqual(t, winRate, 0.95, "Hard won %d/%d games", wins, games)
}

func TestTimedAIAction_RespectsBudget(t *testing.T) {
	g := core.NewDefaultGame()
	start := time.Now()
	_, _, ok := TimedAIAction(g, 100*time.Millisecond)
	elapsed := time.Since(start)
	require.True(t, ok)
	// Allow generous slack: the final in-flight iteration is aborted via
	// periodic clock checks, not preemption.
	assert.Less(t, elapsed, 2*time.Second, "search should stop near its budget")
}

func TestTimedAIAction_GameOverReturnsNoMove(t *testing.T) {
	g, err := core.NewGameFromFEN("r1bqkb1r/pppp1Qpp/2n2n2/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 4")
	require.NoError(t, err)
	_, _, ok := TimedAIAction(g, 50*time.Millisecond)
	assert.False(t, ok)
}

func TestEvaluate_StartPositionIsBalanced(t *testing.T) {
	g := core.NewDefaultGame()
	assert.Equal(t, 0, evaluate(g), "the starting position should evaluate to 0")
}

func TestEvaluate_MaterialAdvantage(t *testing.T) {
	// White up a queen: evaluation from White's perspective must be strongly positive.
	g, err := core.NewGameFromFEN("6k1/8/8/8/8/8/8/3Q2K1 w - - 0 1")
	require.NoError(t, err)
	assert.Greater(t, evaluate(g), 800, "White up a queen should evaluate near +900")

	// Same position with Black to move: negamax convention flips the sign.
	g2, err := core.NewGameFromFEN("6k1/8/8/8/8/8/8/3Q2K1 b - - 0 1")
	require.NoError(t, err)
	assert.Less(t, evaluate(g2), -800, "side to move down a queen should evaluate near -900")
}
