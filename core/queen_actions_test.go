package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueenActions(t *testing.T) {
	ts := []struct {
		name    string
		board   Board
		actions []Action
		color   color
		xy      XY
	}{
		{
			name: "white queen: all actions including captures",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"  ♙♕    ",
					"        ",
					"♙♙  ♙♙♙♙",
					"♖♘♗ ♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{3, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 5}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 6}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 7}},

				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 3}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 2}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{3, 1}}},

				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 4}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{5, 4}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{6, 4}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{7, 4}},

				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 3}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{5, 2}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{6, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{6, 1}}},

				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{2, 3}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{1, 2}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{0, 1}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{0, 1}}},

				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{2, 5}},
				{FromPiece: Piece{PieceType: PieceQueen, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 5}},
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
