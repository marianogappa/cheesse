package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRookActions(t *testing.T) {
	ts := []struct {
		name    string
		board   board
		actions []action
		color   color
		xy      xy
	}{
		{
			name: "black rook at a8: no actions because it's white's turn",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorBlack,
			xy:      xy{0, 0},
			actions: []action{},
		},
		{
			name: "white rook at a1: no actions because it's black's turn",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: 'b',
			},
			color:   colorBlack,
			xy:      xy{0, 7},
			actions: []action{},
		},
		{
			name: "black rook: no actions because it's trapped",
			board: board{
				board: []string{
					" ♞♝♛♚♝♞♜",
					"♟    ♟♟♟",
					"   ♟    ",
					"  ♟♜♟   ",
					"   ♟    ",
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
			name: "white rook: no actions because it's trapped",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"   ♙    ",
					"  ♙♖♙   ",
					"   ♙    ",
					"♙♙    ♙♙",
					" ♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorBlack,
			xy:      xy{3, 4},
			actions: []action{},
		},
		{
			name: "black rook at a8: forwards, including capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: 'b',
			},
			color: colorBlack,
			xy:    xy{0, 0},
			actions: []action{
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 1}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 2}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 3}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 5}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{0, 0}}, toXY: xy{0, 6}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{0, 6}}},
			},
		},
		{
			name: "white rook: all directions, including capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"    ♖   ",
					"        ",
					" ♙♙♙♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color: colorWhite,
			xy:    xy{4, 4},
			actions: []action{
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{0, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{1, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{2, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{3, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{5, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{6, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{7, 4}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{4, 5}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{4, 3}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{4, 2}},
				{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{4, 1}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 1}}},
			},
		},
		{
			name: "white rook: can't move because king is threatened",
			board: board{
				board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝       ",
					"   ♙♖   ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorWhite,
			xy:      xy{4, 4},
			actions: []action{},
		},
		{
			name: "white rook: can only move to block bishop from threatening king",
			board: board{
				board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝       ",
					"    ♖   ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorWhite,
			xy:      xy{4, 4},
			actions: []action{{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 4}}, toXY: xy{1, 4}}},
		},
		{
			name: "white rook: can only capture bishop",
			board: board{
				board: []string{
					"♜♞♝♛♚ ♞♜",
					"♟♟♟♟♟♟ ♟",
					"        ",
					"♝   ♖   ",
					"        ",
					"        ",
					" ♙♙ ♙♙♙♙",
					" ♘♗♕♔♗♘♖",
				},
				turn: 'w',
			},
			color:   colorWhite,
			xy:      xy{4, 3},
			actions: []action{{fromPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{4, 3}}, toXY: xy{0, 3}, isCapture: true, capturedPiece: piece{pieceType: pieceBishop, owner: colorBlack, xy: xy{0, 3}}}},
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
