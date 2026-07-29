package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoubleCheck(t *testing.T) {
	// White bishop at b2, rook at e1, king at h1
	// Black king at e8. Moving Bb2-d4 doesn't give check.
	// For a true double check: Rook on e1, Bishop on a2. Move Nd4 to f5 discovers
	// the rook check on e8 AND knight gives check from f5 doesn't work either.
	// Classic double check: discovered + direct.
	// Position: White Ke1, Bc1 (no, simpler):
	// White: Kg1, Re1, Bd3. Black: Ke8 with an open e-file.
	// Move Bd3-b5+: only direct check, not double.
	// True double check: White Kg1, Qd1, Bb5 on diagonal, Re1 on file.
	// After some piece moves off e-file or diagonal to reveal both.
	//
	// Simplest: White Kh1, Re1, Nd5. Black Ke8.
	// Nf6+ gives check from knight AND discovers rook on e1.
	g, err := NewGameFromFEN("4k3/8/8/3N4/8/8/8/4R2K w - - 0 1")
	require.NoError(t, err)

	var doubleCheckAction Action
	for _, action := range g.Actions {
		// Nd5-f6 gives double check
		if action.FromPiece.PieceType == PieceKnight &&
			action.FromPiece.XY == (XY{X: 3, Y: 3}) &&
			action.ToXY == (XY{X: 5, Y: 2}) {
			doubleCheckAction = action
			break
		}
	}
	require.NotEqual(t, Action{}, doubleCheckAction, "Nd5-f6 should be a legal move")

	newGame := g.DoAction(doubleCheckAction)
	assert.True(t, newGame.IsCheck, "should be check")
	assert.True(t, newGame.IsDoubleCheck, "should be double check (knight + rook)")
	assert.True(t, newGame.IsDiscoverCheck, "should also be discovered check (rook revealed)")
	assert.Len(t, newGame.InCheckBy, 2, "two pieces give check")
}

func TestDiscoveredCheck(t *testing.T) {
	// White: Kh1, Re1, Nd2. Black: Ke8.
	// Move Nd2-f3 discovers the rook on e1 → discovered check (rook checks, not knight).
	g, err := NewGameFromFEN("4k3/8/8/8/8/8/3N4/4R2K w - - 0 1")
	require.NoError(t, err)

	var discoverAction Action
	for _, action := range g.Actions {
		if action.FromPiece.PieceType == PieceKnight &&
			action.FromPiece.XY == (XY{X: 3, Y: 6}) &&
			action.ToXY == (XY{X: 5, Y: 5}) {
			discoverAction = action
			break
		}
	}
	require.NotEqual(t, Action{}, discoverAction, "Nd2-f3 should be a legal move")

	newGame := g.DoAction(discoverAction)
	assert.True(t, newGame.IsCheck, "should be check")
	assert.False(t, newGame.IsDoubleCheck, "should not be double check")
	assert.True(t, newGame.IsDiscoverCheck, "should be discovered check (rook revealed)")
	assert.Len(t, newGame.InCheckBy, 1, "one piece gives check (rook)")
}

func TestDirectCheck(t *testing.T) {
	// Simple direct check: White Kg1, Qd1. Black Ke8.
	// Qd1-d8+ is direct check.
	g, err := NewGameFromFEN("4k3/8/8/8/8/8/8/3Q2K1 w - - 0 1")
	require.NoError(t, err)

	var directAction Action
	for _, action := range g.Actions {
		if action.FromPiece.PieceType == PieceQueen && action.ToXY == (XY{X: 3, Y: 0}) {
			directAction = action
			break
		}
	}
	require.NotEqual(t, Action{}, directAction, "Qd8+ should be a legal move")

	newGame := g.DoAction(directAction)
	assert.True(t, newGame.IsCheck, "should be check")
	assert.False(t, newGame.IsDoubleCheck, "should not be double check")
	assert.False(t, newGame.IsDiscoverCheck, "should not be discovered check")
}

func TestNoCheck(t *testing.T) {
	g := NewDefaultGame()
	// e2-e4 is not check
	var e2e4 Action
	for _, action := range g.Actions {
		if action.FromPiece.PieceType == PiecePawn &&
			action.FromPiece.XY == (XY{X: 4, Y: 6}) &&
			action.ToXY == (XY{X: 4, Y: 4}) {
			e2e4 = action
			break
		}
	}
	newGame := g.DoAction(e2e4)
	assert.False(t, newGame.IsCheck)
	assert.False(t, newGame.IsDoubleCheck)
	assert.False(t, newGame.IsDiscoverCheck)
}
