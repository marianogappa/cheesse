package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIDefaultGame(t *testing.T) {
	expected := OutputGame{
		FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		Board: Board{
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
			CanWhiteKingsideCastle:  true,
			CanWhiteQueensideCastle: true,
			CanBlackKingsideCastle:  true,
			CanBlackQueensideCastle: true,
			HalfMoveClock:           0,
			FullMoveNumber:          1,
			EnPassantTargetSquare:   "",
			Turn:                    "White",
		},
		Actions:                 []OutputAction{},
		CanWhiteCastle:          true,
		CanWhiteKingsideCastle:  true,
		CanWhiteQueensideCastle: true,
		CanBlackCastle:          true,
		CanBlackKingsideCastle:  true,
		CanBlackQueensideCastle: true,
		HalfMoveClock:           0,
		FullMoveNumber:          1,
		IsLastMoveEnPassant:     false,
		EnPassantTargetSquare:   "",
		MoveNumber:              0,
		BlackPieces: map[string]string{
			"a7": "Pawn",
			"a8": "Rook",
			"b7": "Pawn",
			"b8": "Knight",
			"c7": "Pawn",
			"c8": "Bishop",
			"d7": "Pawn",
			"d8": "Queen",
			"e7": "Pawn",
			"e8": "King",
			"f7": "Pawn",
			"f8": "Bishop",
			"g7": "Pawn",
			"g8": "Knight",
			"h7": "Pawn",
			"h8": "Rook",
		},
		WhitePieces: map[string]string{
			"a1": "Rook",
			"a2": "Pawn",
			"b1": "Knight",
			"b2": "Pawn",
			"c1": "Bishop",
			"c2": "Pawn",
			"d1": "Queen",
			"d2": "Pawn",
			"e1": "King",
			"e2": "Pawn",
			"f1": "Bishop",
			"f2": "Pawn",
			"g1": "Knight",
			"g2": "Pawn",
			"h1": "Rook",
			"h2": "Pawn",
		},
		BlackKing:      "e8",
		WhiteKing:      "e1",
		IsCheck:        false,
		IsCheckmate:    false,
		IsStalemate:    false,
		IsDraw:         false,
		IsGameOver:     false,
		GameOverWinner: "Unknown",
		InCheckBy:      []string{},
	}
	actual := New().DefaultGame()
	actual.Actions = []OutputAction{} // Not testing every single action on this test
	assert.Equal(t, expected, actual)
}

func TestParseGame(t *testing.T) {
	testCases := []struct {
		name       string
		inputGame  InputGame
		outputGame OutputGame
		err        error
	}{
		{
			name:      "parses the default game",
			inputGame: InputGame{},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
		},
		{
			name:      "parses the default game by FEN string",
			inputGame: InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
		},
		{
			name: "parses the default game by Board",
			inputGame: InputGame{Board: Board{
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
				CanWhiteKingsideCastle:  true,
				CanWhiteQueensideCastle: true,
				CanBlackKingsideCastle:  true,
				CanBlackQueensideCastle: true,
				HalfMoveClock:           0,
				FullMoveNumber:          1,
				EnPassantTargetSquare:   "",
				Turn:                    "White",
			}},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
		},
		{
			name: "If no FEN string nor board supplied to ParseGame, defaultGame assumed so no error",
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutputGame, err := New().ParseGame(tc.inputGame)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.outputGame.Board.Board, actualOutputGame.Board.Board)
		})
	}
}

func TestDoAction(t *testing.T) {
	testCases := []struct {
		name         string
		inputGame    InputGame
		inputAction  InputAction
		outputGame   OutputGame
		outputAction OutputAction
		err          error
	}{
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: empty FromSquare",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "", ToSquare: "e4"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: empty ToSquare",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e2"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: FromSquare out of bounds file",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "i2", ToSquare: "e4"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: FromSquare out of bounds rank",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e0", ToSquare: "e4"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: ToSquare out of bounds file",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e2", ToSquare: "i4"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "errAlgebraicSquareInvalidOrOutOfBounds: ToSquare out of bounds rank",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e2", ToSquare: "e9"},
			err:         errAlgebraicSquareInvalidOrOutOfBounds,
		},
		{
			name:        "does standard Alekhine on a DefaultGame",
			inputGame:   InputGame{},
			inputAction: InputAction{FromSquare: "e2", ToSquare: "e4"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"    ♙   ",
				"        ",
				"♙♙♙♙ ♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:  "White",
				FromPieceType:   "Pawn",
				FromPieceSquare: "e2",
				ToSquare:        "e4",
				IsEnPassant:     true,
			},
		},
		{
			name:        "does standard Alekhine",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e2", ToSquare: "e4"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"    ♙   ",
				"        ",
				"♙♙♙♙ ♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:  "White",
				FromPieceType:   "Pawn",
				FromPieceSquare: "e2",
				ToSquare:        "e4",
				IsEnPassant:     true,
			},
		},
		{
			name:        "black queenside castles",
			inputGame:   InputGame{FENString: "r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e8", ToSquare: "c8"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"  ♚♜ ♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:    "Black",
				FromPieceType:     "King",
				FromPieceSquare:   "e8",
				ToSquare:          "c8",
				IsCastle:          true,
				IsQueensideCastle: true,
			},
		},
		{
			name:        "black kingside castles",
			inputGame:   InputGame{FENString: "rnbqk2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e8", ToSquare: "g8"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛ ♜♚ ",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕♔♗♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:   "Black",
				FromPieceType:    "King",
				FromPieceSquare:  "e8",
				ToSquare:         "g8",
				IsCastle:         true,
				IsKingsideCastle: true,
			},
		},
		{
			name:        "white queenside castles",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3KBNR w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e1", ToSquare: "c1"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"  ♔♖ ♗♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:    "White",
				FromPieceType:     "King",
				FromPieceSquare:   "e1",
				ToSquare:          "c1",
				IsCastle:          true,
				IsQueensideCastle: true,
			},
		},
		{
			name:        "white kingside castles",
			inputGame:   InputGame{FENString: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 1"},
			inputAction: InputAction{FromSquare: "e1", ToSquare: "g1"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞♝♛♚♝♞♜",
				"♟♟♟♟♟♟♟♟",
				"        ",
				"        ",
				"        ",
				"        ",
				"♙♙♙♙♙♙♙♙",
				"♖♘♗♕ ♖♔ ",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:   "White",
				FromPieceType:    "King",
				FromPieceSquare:  "e1",
				ToSquare:         "g1",
				IsCastle:         true,
				IsKingsideCastle: true,
			},
		},
		{
			name:        "bishop captures bishop",
			inputGame:   InputGame{FENString: "rn1qk1nr/pppb1ppp/4p3/1B1p4/1b1PP3/2N5/PPP2PPP/R1BQK1NR w KQkq - 4 5"},
			inputAction: InputAction{FromSquare: "b5", ToSquare: "d7"},
			outputGame: OutputGame{Board: Board{Board: []string{
				"♜♞ ♛♚ ♞♜",
				"♟♟♟♗ ♟♟♟",
				"    ♟   ",
				"   ♟    ",
				" ♝ ♙♙   ",
				"  ♘     ",
				"♙♙♙  ♙♙♙",
				"♖ ♗♕♔ ♘♖",
			}}},
			outputAction: OutputAction{
				FromPieceOwner:    "White",
				FromPieceType:     "Bishop",
				FromPieceSquare:   "b5",
				ToSquare:          "d7",
				IsCapture:         true,
				CapturedPieceType: "Bishop",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutputGame, actualOutputAction, err := New().DoAction(tc.inputGame, tc.inputAction)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.outputGame.Board.Board, actualOutputGame.Board.Board)
			assert.Equal(t, tc.outputAction, actualOutputAction)
		})
	}
}
