package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnightActions(t *testing.T) {
	ts := []struct {
		name    string
		board   Board
		actions []Action
		color   color
		xy      XY
	}{
		{
			name: "black knight: no actions because it's trapped",
			board: Board{
				Board: []string{
					"♜ ♝♛♚♝♞♜",
					"  ♟ ♟   ",
					" ♟   ♟  ",
					"   ♞    ",
					" ♟   ♟  ",
					"  ♟ ♟   ",
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
			name: "black knight: all actions including 2 captures",
			board: Board{
				Board: []string{
					"♜ ♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"   ♞    ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{3, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{2, 2}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{4, 2}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{1, 3}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{5, 3}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{1, 5}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{5, 5}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{2, 6}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 6}}},
				{FromPiece: Piece{PieceType: PieceKnight, Owner: ColorBlack, XY: XY{3, 4}}, ToXY: XY{4, 6}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{4, 6}}},
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
