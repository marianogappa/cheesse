package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnightActions(t *testing.T) {
	ts := []struct {
		name    string
		board   board
		actions []action
		color   color
		xy      xy
	}{
		{
			name: "black knight: no actions because it's trapped",
			board: board{
				board: []string{
					"♜ ♝♛♚♝♞♜",
					"  ♟ ♟   ",
					" ♟   ♟  ",
					"   ♞    ",
					" ♟   ♟  ",
					"  ♟ ♟   ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{3, 3},
			actions: []action{},
		},
		{
			name: "black knight: all actions including 2 captures",
			board: board{
				board: []string{
					"♜ ♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"   ♞    ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{3, 4},
			actions: []action{
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{2, 2}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{4, 2}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{1, 3}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{5, 3}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{1, 5}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{5, 5}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{2, 6}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 6}}},
				{fromPiece: piece{pieceType: pieceKnight, owner: colorBlack, xy: xy{3, 4}}, toXY: xy{4, 6}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{4, 6}}},
			},
		},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			g, err := newGameFromBoard(tc.board)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.actions, g.pieces[tc.color][tc.xy].calculateAllActions(g))
		})
	}
}
