package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPawnActions(t *testing.T) {
	ts := []struct {
		name    string
		board   board
		actions []action
		color   color
		xy      xy
	}{
		{
			name: "black pawn at a7: no actions because it's white's turn",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorBlack,
			xy:      xy{0, 1},
			actions: []action{},
		},
		{
			name: "white pawn at a2: no actions because it's black's turn",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{0, 6},
			actions: []action{},
		},
		{
			name: "black pawn at a7: forwards and en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{0, 1},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{0, 1}}, toXY: xy{0, 2}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{0, 1}}, toXY: xy{0, 3}, isEnPassant: true},
			},
		},
		{
			name: "black pawn at a6: forwards but no en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"♟       ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{0, 2},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{0, 2}}, toXY: xy{0, 3}},
			},
		},
		{
			name: "black pawn at a6: can't move due to friendly piece forwards",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"  ♟♟♟♟♟♟",
					"♟       ",
					"♟       ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{0, 2},
			actions: []action{},
		},
		{
			name: "black pawn at a6: can't move due to opponent piece forwards",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"  ♟♟♟♟♟♟",
					"♟       ",
					"♙       ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{0, 2},
			actions: []action{},
		},
		{
			name: "black pawn at a7: can't move due to friendly piece forwards, including en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟ ♟♟♟♟♟♟",
					"♟       ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{0, 1},
			actions: []action{},
		},
		{
			name: "black pawn at a7: can't move due to opponent piece forwards, including en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"♙       ",
					"        ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{0, 1},
			actions: []action{},
		},
		{
			name: "white pawn at a2: forwards and en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{0, 6},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{0, 6}}, toXY: xy{0, 5}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{0, 6}}, toXY: xy{0, 4}, isEnPassant: true},
			},
		},
		{
			name: "white pawn at a3: forwards but no en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♙       ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{0, 5},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{0, 5}}, toXY: xy{0, 4}},
			},
		},
		{
			name: "white pawn at a3: can't move due to friendly piece forwards",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"♙       ",
					"♙       ",
					"  ♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorWhite,
			xy:      xy{0, 5},
			actions: []action{},
		},
		{
			name: "white pawn at a3: can't move due to opponent piece forwards",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"♟       ",
					"♙       ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorWhite,
			xy:      xy{0, 5},
			actions: []action{},
		},
		{
			name: "white pawn at a2: can't move due to friendly piece forwards, including en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♙       ",
					"♙ ♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorWhite,
			xy:      xy{0, 6},
			actions: []action{},
		},
		{
			name: "white pawn at a2: can't move due to opponent piece forwards, including en passant",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♟       ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorWhite,
			xy:      xy{0, 6},
			actions: []action{},
		},
		{
			name: "white pawn at d4: forward or capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟ ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{3, 4},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{3, 3}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{4, 3}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}},
			},
		},
		{
			name: "white pawn at d4: forward, and not capture because friendly piece",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"    ♙   ",
					"   ♙    ",
					"        ",
					"♙♙♙  ♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{3, 4},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}, toXY: xy{3, 3}},
			},
		},
		{
			name: "black pawn at e5: forward or capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟ ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{4, 3},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}, toXY: xy{4, 4}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}, toXY: xy{3, 4}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}},
			},
		},
		{
			name: "black pawn at e5: forward and no capture because friendly piece",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟  ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♟    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{4, 3},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}, toXY: xy{4, 4}},
			},
		},
		{
			name: "black pawn at e5: two captures",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟  ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙♟♙  ",
					"        ",
					"♙♙♙ ♙ ♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{4, 3},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}, toXY: xy{3, 4}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{4, 3}}, toXY: xy{5, 4}, isCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{5, 4}}},
			},
		},
		{
			name: "black pawn: forwards and en passant capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"  ♟♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn:                  "Black",
				enPassantTargetSquare: "d3",
			},
			color: colorBlack,
			xy:    xy{2, 4},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 4}}, toXY: xy{2, 5}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 4}}, toXY: xy{3, 5}, isCapture: true, isEnPassantCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 4}}},
			},
		},
		{
			name: "white pawn: forwards and en passant capture",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"  ♟♙    ",
					"        ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn:                  "White",
				enPassantTargetSquare: "c6",
			},
			color: colorWhite,
			xy:    xy{3, 3},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 3}}, toXY: xy{3, 2}},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{3, 3}}, toXY: xy{2, 2}, isCapture: true, isEnPassantCapture: true, capturedPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 3}}},
			},
		},
		{
			name: "black pawn: can't promote",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"  ♔     ",
				},
				turn: "Black",
			},
			color:   colorBlack,
			xy:      xy{2, 6},
			actions: []action{},
		},
		{
			name: "black pawn: promotes to all pieces",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"     ♔  ",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{2, 6},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceRook},
			},
		},
		{
			name: "black pawn: promotes to all pieces while capturing",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"  ♔♖    ",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{2, 6},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceRook},
			},
		},
		{
			name: "black pawn: promotes to all pieces while capturing and while advancing",
			board: board{
				board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"   ♖  ♔ ",
				},
				turn: "Black",
			},
			color: colorBlack,
			xy:    xy{2, 6},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{2, 7}, isPromotion: true, promotionPieceType: pieceRook},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorBlack, xy: xy{2, 6}}, toXY: xy{3, 7}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorWhite, xy: xy{3, 7}}, isPromotion: true, promotionPieceType: pieceRook},
			},
		},
		{
			name: "white pawn: can't promote",
			board: board{
				board: []string{
					"  ♚     ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color:   colorWhite,
			xy:      xy{2, 1},
			actions: []action{},
		},
		{
			name: "white pawn: promotes to all pieces",
			board: board{
				board: []string{
					"     ♚  ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{2, 1},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceRook},
			},
		},
		{
			name: "white pawn: promotes to all pieces while capturing",
			board: board{
				board: []string{
					"  ♚♜    ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{2, 1},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceRook},
			},
		},
		{
			name: "white pawn: promotes to all pieces while capturing and advancing",
			board: board{
				board: []string{
					"   ♜  ♚ ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				turn: "White",
			},
			color: colorWhite,
			xy:    xy{2, 1},
			actions: []action{
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{2, 0}, isPromotion: true, promotionPieceType: pieceRook},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceQueen},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceBishop},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceKnight},
				{fromPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{2, 1}}, toXY: xy{3, 0}, isCapture: true, capturedPiece: piece{pieceType: pieceRook, owner: colorBlack, xy: xy{3, 0}}, isPromotion: true, promotionPieceType: pieceRook},
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
