package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoActionWithActionString(t *testing.T) {
	testCases := []struct {
		name         string
		actionString string
		fromSquare   string
		toSquare     string
	}{
		{"algebraic piece move", "Nf3", "g1", "f3"},
		{"algebraic pawn move", "e4", "e2", "e4"},
		{"figurine", "♘f3", "g1", "f3"},
		{"coordinate", "g1-f3", "g1", "f3"},
		{"smith", "g1f3", "g1", "f3"},
		{"descriptive", "N-KB3", "g1", "f3"},
		{"ICCF", "7163", "g1", "f3"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outputGame, outputAction, err := New().DoAction(InputGame{}, InputAction{ActionString: tc.actionString})
			require.NoError(t, err)
			assert.Equal(t, tc.fromSquare, outputAction.FromPieceSquare)
			assert.Equal(t, tc.toSquare, outputAction.ToSquare)
			assert.Equal(t, "Black", outputGame.Board.Turn)
		})
	}
}

func TestDoActionWithActionStringPromotion(t *testing.T) {
	outputGame, outputAction, err := New().DoAction(
		InputGame{FENString: "8/5P1k/8/8/8/8/8/K7 w - - 0 1"},
		InputAction{ActionString: "f8=N"},
	)
	require.NoError(t, err)
	assert.True(t, outputAction.IsPromotion)
	assert.Equal(t, "Knight", outputAction.PromotionPieceType)
	assert.Equal(t, "Knight", outputGame.WhitePieces["f8"])
}

func TestDoActionWithActionStringCastle(t *testing.T) {
	_, outputAction, err := New().DoAction(
		InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1"},
		InputAction{ActionString: "O-O"},
	)
	require.NoError(t, err)
	assert.True(t, outputAction.IsKingsideCastle)
}

func TestDoActionWithActionStringErrors(t *testing.T) {
	t.Run("illegal move", func(t *testing.T) {
		_, _, err := New().DoAction(InputGame{}, InputAction{ActionString: "Qh5"})
		assert.Equal(t, errInvalidActionForGivenGame, err)
	})
	t.Run("garbage input", func(t *testing.T) {
		_, _, err := New().DoAction(InputGame{}, InputAction{ActionString: "xyzzy"})
		assert.Equal(t, errInvalidActionForGivenGame, err)
	})
	t.Run("multiple moves are rejected", func(t *testing.T) {
		_, _, err := New().DoAction(InputGame{}, InputAction{ActionString: "1. e4 e5"})
		assert.Equal(t, errInvalidActionForGivenGame, err)
	})
	t.Run("ambiguous move is rejected", func(t *testing.T) {
		// Two knights can reach d2: Nbd2 or Nfd2 required.
		_, _, err := New().DoAction(
			InputGame{FENString: "4k3/8/8/8/8/5N2/8/RN2K3 w - - 0 1"},
			InputAction{ActionString: "Nd2"},
		)
		assert.Error(t, err)
	})
}

func TestDoActionWithActionStringTakesPrecedenceOverSquares(t *testing.T) {
	// When both actionString and squares are supplied, actionString wins.
	_, outputAction, err := New().DoAction(InputGame{}, InputAction{
		ActionString: "d4",
		FromSquare:   "e2",
		ToSquare:     "e4",
	})
	require.NoError(t, err)
	assert.Equal(t, "d2", outputAction.FromPieceSquare)
	assert.Equal(t, "d4", outputAction.ToSquare)
}
