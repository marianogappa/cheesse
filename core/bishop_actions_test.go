package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBishopActions(t *testing.T) {
	ts := []struct {
		name    string
		board   Board
		actions []Action
		color   color
		xy      XY
	}{
		{
			name: "black bishop: no actions because it's trapped",
			board: Board{
				Board: []string{
					"♜♞ ♛♚♝♞♜",
					"♟    ♟♟♟",
					"  ♟ ♟   ",
					"   ♝    ",
					"  ♟ ♟   ",
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
			name: "white bishop: no actions because it's trapped",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"  ♙ ♙   ",
					"   ♗    ",
					"  ♙ ♙   ",
					"♙♙    ♙♙",
					"♖♘ ♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{3, 4},
			actions: []Action{},
		},
		{
			name: "white bishop: all actions including 2 captures",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"   ♗    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘ ♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{3, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{2, 5}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 5}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{2, 3}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{1, 2}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 3}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{5, 2}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{0, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{0, 1}}},
				{FromPiece: Piece{PieceType: PieceBishop, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{6, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{6, 1}}},
			},
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
