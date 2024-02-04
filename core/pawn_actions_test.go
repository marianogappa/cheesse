package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPawnActions(t *testing.T) {
	ts := []struct {
		name    string
		board   Board
		actions []Action
		color   color
		xy      XY
	}{
		{
			name: "black pawn at a7: no actions because it's white's turn",
			board: Board{
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
				Turn: "White",
			},
			color:   ColorBlack,
			xy:      XY{0, 1},
			actions: []Action{},
		},
		{
			name: "white pawn at a2: no actions because it's black's turn",
			board: Board{
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
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 6},
			actions: []Action{},
		},
		{
			name: "black pawn at a7: forwards and en passant",
			board: Board{
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
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{0, 1},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{0, 1}}, ToXY: XY{0, 2}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{0, 1}}, ToXY: XY{0, 3}, IsEnPassant: true},
			},
		},
		{
			name: "black pawn at a6: forwards but no en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"♟       ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{0, 2},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{0, 2}}, ToXY: XY{0, 3}},
			},
		},
		{
			name: "black pawn at a6: can't move due to friendly piece forwards",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"  ♟♟♟♟♟♟",
					"♟       ",
					"♟       ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 2},
			actions: []Action{},
		},
		{
			name: "black pawn at a6: can't move due to opponent piece forwards",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"  ♟♟♟♟♟♟",
					"♟       ",
					"♙       ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 2},
			actions: []Action{},
		},
		{
			name: "black pawn at a7: can't move due to friendly piece forwards, including en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟ ♟♟♟♟♟♟",
					"♟       ",
					"        ",
					"        ",
					"        ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 1},
			actions: []Action{},
		},
		{
			name: "black pawn at a7: can't move due to opponent piece forwards, including en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"♙       ",
					"        ",
					"        ",
					"        ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{0, 1},
			actions: []Action{},
		},
		{
			name: "white pawn at a2: forwards and en passant",
			board: Board{
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
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{0, 6},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{0, 6}}, ToXY: XY{0, 5}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{0, 6}}, ToXY: XY{0, 4}, IsEnPassant: true},
			},
		},
		{
			name: "white pawn at a3: forwards but no en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♙       ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{0, 5},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{0, 5}}, ToXY: XY{0, 4}},
			},
		},
		{
			name: "white pawn at a3: can't move due to friendly piece forwards",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"♙       ",
					"♙       ",
					"  ♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{0, 5},
			actions: []Action{},
		},
		{
			name: "white pawn at a3: can't move due to opponent piece forwards",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"♟       ",
					"♙       ",
					" ♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{0, 5},
			actions: []Action{},
		},
		{
			name: "white pawn at a2: can't move due to friendly piece forwards, including en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♙       ",
					"♙ ♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{0, 6},
			actions: []Action{},
		},
		{
			name: "white pawn at a2: can't move due to opponent piece forwards, including en passant",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					" ♟♟♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"♟       ",
					"♙♙♙♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{0, 6},
			actions: []Action{},
		},
		{
			name: "white pawn at d4: forward or capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟ ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{3, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 3}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{4, 3}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}},
			},
		},
		{
			name: "white pawn at d4: forward, and not capture because friendly piece",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟♟♟♟♟",
					"        ",
					"    ♙   ",
					"   ♙    ",
					"        ",
					"♙♙♙  ♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{3, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}, ToXY: XY{3, 3}},
			},
		},
		{
			name: "black pawn at e5: forward or capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟♟ ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{4, 3},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}, ToXY: XY{4, 4}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}, ToXY: XY{3, 4}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}},
			},
		},
		{
			name: "black pawn at e5: forward and no capture because friendly piece",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟  ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♟    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{4, 3},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}, ToXY: XY{4, 4}},
			},
		},
		{
			name: "black pawn at e5: two captures",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟♟  ♟♟♟",
					"        ",
					"    ♟   ",
					"   ♙♟♙  ",
					"        ",
					"♙♙♙ ♙ ♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{4, 3},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}, ToXY: XY{3, 4}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{4, 3}}, ToXY: XY{5, 4}, IsCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{5, 4}}},
			},
		},
		{
			name: "black pawn: forwards and en passant capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"  ♟♙    ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn:                  "Black",
				EnPassantTargetSquare: "d3",
			},
			color: ColorBlack,
			xy:    XY{2, 4},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 4}}, ToXY: XY{2, 5}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 4}}, ToXY: XY{3, 5}, IsCapture: true, IsEnPassantCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 4}}},
			},
		},
		{
			name: "white pawn: forwards and en passant capture",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"  ♟♙    ",
					"        ",
					"        ",
					"♙♙♙ ♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn:                  "White",
				EnPassantTargetSquare: "c6",
			},
			color: ColorWhite,
			xy:    XY{3, 3},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 3}}, ToXY: XY{3, 2}},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{3, 3}}, ToXY: XY{2, 2}, IsCapture: true, IsEnPassantCapture: true, CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 3}}},
			},
		},
		{
			name: "black pawn: can't promote",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"  ♔     ",
				},
				Turn: "Black",
			},
			color:   ColorBlack,
			xy:      XY{2, 6},
			actions: []Action{},
		},
		{
			name: "black pawn: promotes to all pieces",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"     ♔  ",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{2, 6},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceRook},
			},
		},
		{
			name: "black pawn: promotes to all pieces while capturing",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"  ♔♖    ",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{2, 6},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceRook},
			},
		},
		{
			name: "black pawn: promotes to all pieces while capturing and while advancing",
			board: Board{
				Board: []string{
					"♜♞♝♛♚♝♞♜",
					"♟♟ ♟♟♟♟♟",
					"        ",
					"        ",
					"        ",
					"        ",
					"  ♟     ",
					"   ♖  ♔ ",
				},
				Turn: "Black",
			},
			color: ColorBlack,
			xy:    XY{2, 6},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{2, 7}, IsPromotion: true, PromotionPieceType: PieceRook},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{2, 6}}, ToXY: XY{3, 7}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorWhite, XY: XY{3, 7}}, IsPromotion: true, PromotionPieceType: PieceRook},
			},
		},
		{
			name: "white pawn: can't promote",
			board: Board{
				Board: []string{
					"  ♚     ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color:   ColorWhite,
			xy:      XY{2, 1},
			actions: []Action{},
		},
		{
			name: "white pawn: promotes to all pieces",
			board: Board{
				Board: []string{
					"     ♚  ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{2, 1},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceRook},
			},
		},
		{
			name: "white pawn: promotes to all pieces while capturing",
			board: Board{
				Board: []string{
					"  ♚♜    ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{2, 1},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceRook},
			},
		},
		{
			name: "white pawn: promotes to all pieces while capturing and advancing",
			board: Board{
				Board: []string{
					"   ♜  ♚ ",
					"  ♙     ",
					"        ",
					"        ",
					"        ",
					"        ",
					"♙♙ ♙♙♙♙♙",
					"♖♘♗♕♔♗♘♖",
				},
				Turn: "White",
			},
			color: ColorWhite,
			xy:    XY{2, 1},
			actions: []Action{
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{2, 0}, IsPromotion: true, PromotionPieceType: PieceRook},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceQueen},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceBishop},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceKnight},
				{FromPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{2, 1}}, ToXY: XY{3, 0}, IsCapture: true, CapturedPiece: Piece{PieceType: PieceRook, Owner: ColorBlack, XY: XY{3, 0}}, IsPromotion: true, PromotionPieceType: PieceRook},
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
