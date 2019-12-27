package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBishopActions(t *testing.T) {
	ts := []struct {
		name    string
		board   board
		actions []action
		color   color
		xy      xy
	}{
		{
			name: "black bishop: no actions because it's trapped",
			board: board{
				board: []string{
					"♜♞ ♛♚♝♞♜",
					"♟    ♟♟♟",
					"  ♟ ♟   ",
					"   ♝    ",
					"  ♟ ♟   ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: 'b',
			},
			color:   colorBlack,
			xy:      xy{3, 3},
			actions: []action{},
		},
		{
			name: "white bishop: no actions because it's trapped",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"  ♙ ♙   ",
					"   ♗    ",
					"  ♙ ♙   ",
					"♙♙    ♙♙",
					"♖♘ ♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorWhite,
			xy:      xy{3, 4},
			actions: []action{},
		},
		{
			name: "white bishop: all actions including 2 captures",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"   ♗    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘ ♕♔♗♘♖",
				},
				turn: 'w',
			},
			color: colorWhite,
			xy:    xy{3, 4},
			actions: []action{
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{2, 5}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{4, 5}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{2, 3}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{1, 2}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{4, 3}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{5, 2}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{0, 1}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{0, 1}}},
				{fromPiece: piece{pieceType: pieceBishop, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{6, 1}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{6, 1}}},
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
