package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsufficientMaterial(t *testing.T) {
	ts := []struct {
		name   string
		fen    string
		isDraw bool
	}{
		{
			name:   "K vs K",
			fen:    "8/8/4k3/8/8/4K3/8/8 w - - 0 1",
			isDraw: true,
		},
		{
			name:   "KB vs K (white bishop)",
			fen:    "8/8/4k3/8/8/2B1K3/8/8 w - - 0 1",
			isDraw: true,
		},
		{
			name:   "KB vs K (black bishop)",
			fen:    "8/8/2b1k3/8/8/4K3/8/8 w - - 0 1",
			isDraw: true,
		},
		{
			name:   "KN vs K (white knight)",
			fen:    "8/8/4k3/8/8/2N1K3/8/8 w - - 0 1",
			isDraw: true,
		},
		{
			name:   "KN vs K (black knight)",
			fen:    "8/8/2n1k3/8/8/4K3/8/8 w - - 0 1",
			isDraw: true,
		},
		{
			name: "KB vs KB with same-color bishops",
			// c1 and f4 are both dark squares
			fen:    "8/8/4k3/8/5b2/4K3/8/2B5 w - - 0 1",
			isDraw: true,
		},
		{
			name: "KB vs KB with opposite-color bishops is not a draw",
			// c1 dark, e4 light
			fen:    "8/8/4k3/8/4b3/4K3/8/2B5 w - - 0 1",
			isDraw: false,
		},
		{
			name:   "KNN vs K is not automatically a draw",
			fen:    "8/8/4k3/8/8/1NN1K3/8/8 w - - 0 1",
			isDraw: false,
		},
		{
			name:   "K vs K with a pawn is not a draw",
			fen:    "8/8/4k3/8/8/4K3/4P3/8 w - - 0 1",
			isDraw: false,
		},
		{
			name:   "K vs K with a rook is not a draw",
			fen:    "8/8/4k3/8/8/4K3/8/6R1 w - - 0 1",
			isDraw: false,
		},
		{
			name:   "K vs K with a queen is not a draw",
			fen:    "8/8/4k3/8/8/4K3/8/6Q1 w - - 0 1",
			isDraw: false,
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			assert.Equal(t, tc.isDraw, g.IsDraw)
			assert.Equal(t, tc.isDraw, g.IsGameOver)
		})
	}
}

func TestFiftyMoveRule(t *testing.T) {
	t.Run("clock at 99 is not claimable", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/4k3/8/8/4K3/8/6R1 w - - 99 80")
		require.NoError(t, err)
		assert.False(t, g.CanClaimDraw)
		assert.False(t, g.IsDraw)
	})
	t.Run("clock at 100 is claimable but not automatic", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/4k3/8/8/4K3/8/6R1 w - - 100 80")
		require.NoError(t, err)
		assert.True(t, g.CanClaimDraw)
		assert.False(t, g.IsDraw)
		assert.False(t, g.IsGameOver)
	})
	t.Run("clock above 100 is still claimable (no exact-equality bug)", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/4k3/8/8/4K3/8/6R1 w - - 120 90")
		require.NoError(t, err)
		assert.True(t, g.CanClaimDraw)
		assert.False(t, g.IsDraw)
	})
	t.Run("clock at 150 is an automatic draw (75-move rule)", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/4k3/8/8/4K3/8/6R1 w - - 150 100")
		require.NoError(t, err)
		assert.True(t, g.IsDraw)
		assert.True(t, g.IsGameOver)
	})
	t.Run("clock reached via DoAction", func(t *testing.T) {
		g, err := NewGameFromFEN("8/8/4k3/8/8/4K3/8/6R1 w - - 99 80")
		require.NoError(t, err)
		var action Action
		for _, a := range g.Actions {
			if a.FromPiece.PieceType == PieceRook && !a.IsCapture {
				action = a
				break
			}
		}
		require.NotEqual(t, Action{}, action)
		newGame := g.DoAction(action)
		assert.Equal(t, 100, newGame.HalfMoveClock)
		assert.True(t, newGame.CanClaimDraw)
		assert.False(t, newGame.IsDraw)
	})
}

// doMoves applies a sequence of from→to moves in algebraic-square pairs.
func doMoves(t *testing.T, g Game, moves ...string) Game {
	t.Helper()
	for i := 0; i < len(moves); i += 2 {
		from := XY{int(moves[i][0] - 'a'), int('8' - moves[i][1])}
		to := XY{int(moves[i+1][0] - 'a'), int('8' - moves[i+1][1])}
		var action Action
		for _, a := range g.Actions {
			if !a.IsResign && a.FromPiece.XY == from && a.ToXY == to {
				action = a
				break
			}
		}
		require.NotEqual(t, Action{}, action, "move %s%s not found", moves[i], moves[i+1])
		g = g.DoAction(action)
	}
	return g
}

func TestThreefoldRepetition(t *testing.T) {
	t.Run("shuffling knights back and forth reaches threefold", func(t *testing.T) {
		g := NewDefaultGame()

		// Ng1-f3 Ng8-f6 Nf3-g1 Nf6-g8 repeats the starting position; twice over reaches
		// the third occurrence of the initial position.
		g = doMoves(t, g,
			"g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8", // initial position x2
			"g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8", // initial position x3
		)
		assert.True(t, g.CanClaimDraw, "third occurrence should be claimable")
		assert.False(t, g.IsDraw, "threefold is claimable, not automatic")
		assert.False(t, g.IsGameOver)
	})

	t.Run("two occurrences is not claimable", func(t *testing.T) {
		g := NewDefaultGame()
		g = doMoves(t, g, "g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8")
		assert.False(t, g.CanClaimDraw)
	})

	t.Run("fivefold repetition is an automatic draw", func(t *testing.T) {
		g := NewDefaultGame()
		for i := 0; i < 4; i++ {
			g = doMoves(t, g, "g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8")
		}
		assert.True(t, g.IsDraw)
		assert.True(t, g.IsGameOver)
	})

	t.Run("pawn move resets repetition tracking", func(t *testing.T) {
		g := NewDefaultGame()
		g = doMoves(t, g,
			"g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8",
			"e2", "e4", "e7", "e5", // irreversible: history restarts
			"g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8",
		)
		assert.False(t, g.CanClaimDraw, "repetitions before a pawn move don't count")
	})

	t.Run("history crosses stateless boundary via WithPositionHistory", func(t *testing.T) {
		g := NewDefaultGame()
		g = doMoves(t, g, "g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8")
		history := g.PositionHistory()

		// Reconstruct the same position from FEN (stateless boundary) and restore history.
		restored, err := NewGameFromFEN(g.ToFEN())
		require.NoError(t, err)
		restored = restored.WithPositionHistory(history)

		restored = doMoves(t, restored, "g1", "f3", "g8", "f6", "f3", "g1", "f6", "g8")
		assert.True(t, restored.CanClaimDraw, "restored history should enable threefold detection")
	})
}

func TestStalemateIsNotAffectedByDrawRules(t *testing.T) {
	// Classic stalemate: black king on a8, white queen on b6, white king on c6; black to move.
	g, err := NewGameFromFEN("k7/8/1QK5/8/8/8/8/8 b - - 0 1")
	require.NoError(t, err)
	assert.True(t, g.IsStalemate)
	assert.True(t, g.IsGameOver)
	assert.False(t, g.IsCheckmate)
}
