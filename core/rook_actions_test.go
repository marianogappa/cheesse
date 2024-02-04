package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRookActions(t *testing.T) {
	ts := []struct {
		name    string
		board   Board
		actions []Action
		color   color
		xy      XY
	}{
		{
			name: "black rook at a8: no actions because it's white's turn",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorBlack,
			xy:      XY{0, 0},
			actions: []Action{},
		},
		{
			name: "white rook at a1: no actions because it's black's turn",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 7},
			actions: []Action{},
		},
		{
			name: "black rook: no actions because it's trapped",
			board: Board{
				Board: []string{
					" ♞♝♛♚♝♞♜",
					"♟    ♟♟♟",
					"   ♟    ",
					"  ♟♜♟   ",
					"   ♟    ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{3, 3},
			actions: []Action{},
		},
		{
			name: "white rook: no actions because it's trapped",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"   ♙    ",
					"  ♙♖♙   ",
					"   ♙    ",
					"♙♙    ♙♙",
					" ♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorBlack,
			xy:      XY{3, 4},
			actions: []Action{},
		},
		{
			name: "black rook at a8: forwards, including capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{0, 0},
			actions: []Action{
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 1}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 2}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 3}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 5}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{0, 0}}, ToXY: XY{0, 6}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{0, 6}}},
			},
		},
		{
			name: "white rook: all directions, including capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"    ♖   ",
					"        ",
					" ♙♙♙♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{4, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{0, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{1, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{2, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{3, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{5, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{6, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{7, 4}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{4, 5}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{4, 3}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{4, 2}},
				{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{4, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 1}}},
			},
		},
		{
			name: "white rook: can't move because king is threatened",
			board: Board{
				Board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝       ",
					"   ♙♖   ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{4, 4},
			actions: []Action{},
		},
		{
			name: "white rook: can only move to block bishop from threatening king",
			board: Board{
				Board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝       ",
					"    ♖   ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{4, 4},
			actions: []Action{{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 4}}, ToXY: XY{1, 4}}},
		},
		{
			name: "white rook: can only capture bishop",
			board: Board{
				Board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝   ♖   ",
					"        ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{4, 3},
			actions: []Action{{FromPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{4, 3}}, ToXY: XY{0, 3}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceBishop, Owner: ColorBlack, XY: XY{0, 3}}}},
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewGameFromBoard(tc.board)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.actions, g.Pieces[tc.color][tc.xy].calculateAllActions(g))
		})
	}
}
