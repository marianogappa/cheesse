package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Board editors send naive FENs: castling always "KQkq" regardless of piece placement,
// arbitrary positions, etc. These tests assert the lenient-parsing policy: auto-correct
// what's mechanically derivable, reject only what's truly unplayable.

func TestFENAutoCorrectsCastlingRights(t *testing.T) {
	testCases := []struct {
		name             string
		fen              string
		expectedCastling string // the castling field of the re-rendered FEN
	}{
		{
			name:             "moved white king narrows KQkq to kq",
			fen:              "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R2K3R w KQkq - 0 1",
			expectedCastling: "kq",
		},
		{
			name:             "moved black king narrows KQkq to KQ",
			fen:              "r4k1r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			expectedCastling: "KQ",
		},
		{
			name:             "missing white kingside rook narrows K",
			fen:              "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K3 w KQkq - 0 1",
			expectedCastling: "Qkq",
		},
		{
			name:             "missing white queenside rook narrows Q",
			fen:              "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/4K2R w KQkq - 0 1",
			expectedCastling: "Kkq",
		},
		{
			name:             "missing black kingside rook narrows k",
			fen:              "r3k3/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			expectedCastling: "KQq",
		},
		{
			name:             "missing black queenside rook narrows q",
			fen:              "4k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
			expectedCastling: "KQk",
		},
		{
			name:             "editor position with no castling material at all",
			fen:              "8/1k6/8/8/8/8/6K1/8 w KQkq - 0 1",
			expectedCastling: "-",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			fenParts := splitFEN(t, g.ToFEN())
			assert.Equal(t, tc.expectedCastling, fenParts[2])
		})
	}
}

func TestFENAutoCorrectsImpossibleEnPassant(t *testing.T) {
	testCases := []struct {
		name string
		fen  string
	}{
		{
			name: "e.p. square with no pawn next to it is dropped",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e6 0 1",
		},
		{
			name: "e.p. square with black pawn that didn't double-advance is dropped",
			fen:  "rnbqkbnr/ppppp1pp/5p2/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e6 0 1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromFEN(tc.fen)
			require.NoError(t, err)
			assert.False(t, g.IsLastMoveEnPassant)
			fenParts := splitFEN(t, g.ToFEN())
			assert.Equal(t, "-", fenParts[3])
		})
	}
}

func TestFENRejectsTrulyUnplayablePositions(t *testing.T) {
	testCases := []struct {
		name string
		fen  string
		err  error
	}{
		{
			name: "no white king",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1BNR w KQkq - 0 1",
			err:  errFENKingMissing,
		},
		{
			name: "two black kings",
			fen:  "rnbqkbnr/pppppppp/8/4k3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			err:  errFENDuplicateKing,
		},
		{
			name: "pawn on back rank",
			fen:  "Pnbqkbnr/1ppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			err:  errFENPawnInImpossibleRank,
		},
		{
			name: "side not to move is in check",
			fen:  "4k3/4R3/8/8/8/8/8/4K3 w - - 0 1",
			err:  errFENSideNotToMoveInCheck,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGameFromFEN(tc.fen)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestEditorPositionsAreFullyPlayable(t *testing.T) {
	// Simulates converter-demo board editor positions with the naive
	// "<placement> w KQkq - 0 1" FEN construction: parses must succeed, actions
	// must be generated, and DoAction must not panic.
	editorFENs := []string{
		"8/1k6/8/8/8/8/6K1/8 w KQkq - 0 1",             // two bare kings
		"4k3/8/8/8/8/8/4P3/4K3 w KQkq - 0 1",           // king and pawn
		"1k6/8/8/8/8/8/8/R3K2R w KQkq - 0 1",           // castling-ready white, bare black king
		"k7/pppppppp/8/8/8/8/8/7K b KQkq - 0 1",        // black pawn wall
		"4k3/8/8/3q4/8/8/8/4K3 b KQkq - 0 1",           // queen endgame, black to move
	}
	for _, fen := range editorFENs {
		t.Run(fen, func(t *testing.T) {
			g, err := NewGameFromFEN(fen)
			require.NoError(t, err)
			require.NotEmpty(t, g.Actions)
			// Play every available action; none may panic.
			for _, action := range g.Actions {
				_ = g.DoAction(action)
			}
		})
	}
}

func TestBoardAutoCorrection(t *testing.T) {
	t.Run("board with all castling flags but moved king parses with narrowed rights", func(t *testing.T) {
		b := Board{
			Board: []string{
				"♜   ♚  ♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖  ♔   ♖",
			},
			CanWhiteKingsideCastle:  true,
			CanWhiteQueensideCastle: true,
			CanBlackKingsideCastle:  true,
			CanBlackQueensideCastle: true,
			FullMoveNumber:          1,
			Turn:                    "White",
		}
		g, err := NewGameFromBoard(b)
		require.NoError(t, err)
		assert.False(t, g.CanWhiteKingsideCastle)
		assert.False(t, g.CanWhiteQueensideCastle)
		assert.False(t, g.CanWhiteCastle)
		assert.True(t, g.CanBlackKingsideCastle)
		assert.True(t, g.CanBlackQueensideCastle)
		assert.True(t, g.CanBlackCastle)
	})

	t.Run("board with impossible en passant parses with e.p. dropped", func(t *testing.T) {
		b := Board{
			Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			},
			FullMoveNumber:        1,
			EnPassantTargetSquare: "e6",
			Turn:                  "White",
		}
		g, err := NewGameFromBoard(b)
		require.NoError(t, err)
		assert.False(t, g.IsLastMoveEnPassant)
	})

	t.Run("board with side not to move in check is rejected", func(t *testing.T) {
		b := Board{
			Board: []string{
				"    ♚   ",
				"    ♖   ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"    ♔   ",
			},
			FullMoveNumber: 1,
			Turn:           "White",
		}
		_, err := NewGameFromBoard(b)
		assert.Equal(t, errBoardSideNotToMoveInCheck, err)
	})
}

// TestBoardFENCastlingConsistency verifies the CanXCastle aggregate is OR (not AND)
// of the side rights for both Board and FEN inputs.
func TestBoardFENCastlingConsistency(t *testing.T) {
	b := Board{
		Board: []string{
			"♜   ♚   ",
			"♟♟♟♟♟♟♟♟",
			"        ",
			"        ",
			"        ",
			"        ",
			"♙♙♙♙♙♙♙♙",
			"♖   ♔   ",
		},
		CanWhiteKingsideCastle:  false,
		CanWhiteQueensideCastle: true,
		CanBlackKingsideCastle:  false,
		CanBlackQueensideCastle: true,
		FullMoveNumber:          1,
		Turn:                    "White",
	}
	gBoard, err := NewGameFromBoard(b)
	require.NoError(t, err)

	gFEN, err := NewGameFromFEN("r3k3/pppppppp/8/8/8/8/PPPPPPPP/R3K3 w Qq - 0 1")
	require.NoError(t, err)

	assert.True(t, gBoard.CanWhiteCastle, "queenside-only castling must still mean CanWhiteCastle")
	assert.True(t, gBoard.CanBlackCastle)
	assert.Equal(t, gFEN.CanWhiteCastle, gBoard.CanWhiteCastle)
	assert.Equal(t, gFEN.CanBlackCastle, gBoard.CanBlackCastle)
	assert.Equal(t, gFEN.ToFEN(), gBoard.ToFEN())
}

func splitFEN(t *testing.T, fen string) []string {
	t.Helper()
	parts := strings.Split(fen, " ")
	require.Len(t, parts, 6)
	return parts
}
