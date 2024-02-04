package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameString(t *testing.T) {
	expected := `♜♞♝♛♚♝♞♜
♟♟♟♟♟♟♟♟
........
........
........
........
♙♙♙♙♙♙♙♙
♖♘♗♕♔♗♘♖
`
	g, err := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)
	assert.Equal(t, expected, g.String())
}

func TestActionString(t *testing.T) {
	testCases := []struct {
		a Action
		s string
	}{
		{
			a: Action{
				FromPiece: Piece{PieceType: PieceQueen, Owner: ColorBlack, XY: XY{2, 2}},
				ToXY:      XY{3, 3},
			},
			s: "Black's Queen at c6 moves to d5",
		},
		{
			a: Action{
				FromPiece: Piece{Owner: ColorWhite},
				IsResign:  true,
			},
			s: "White resigns",
		},
		{
			a: Action{
				FromPiece:     Piece{PieceType: PieceBishop, Owner: ColorBlack, XY: XY{1, 1}},
				ToXY:          XY{4, 4},
				IsCapture:     true,
				CapturedPiece: Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{4, 4}},
			},
			s: "Black's Bishop at b7 captures White's Pawn at e4",
		},
		{
			a: Action{
				FromPiece:          Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{4, 1}},
				ToXY:               XY{4, 0},
				IsPromotion:        true,
				PromotionPieceType: PieceKnight,
			},
			s: "White's Pawn at e7 promotes to Knight",
		},
		{
			a: Action{
				FromPiece:   Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{5, 1}},
				ToXY:        XY{5, 3},
				IsEnPassant: true,
			},
			s: "Black's Pawn at f7 does en passant",
		},
		{
			a: Action{
				FromPiece:          Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{3, 5}},
				ToXY:               XY{4, 6},
				IsCapture:          true,
				IsEnPassantCapture: true,
				CapturedPiece:      Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{4, 5}},
			},
			s: "Black's Pawn at d3 captures White's Pawn at e3 which was doing en passant",
		},
		{
			a: Action{
				FromPiece:         Piece{PieceType: PieceKing, Owner: ColorBlack, XY: XY{4, 0}},
				ToXY:              XY{2, 0},
				IsCastle:          true,
				IsQueensideCastle: true,
			},
			s: "Black queenside castles",
		},
		{
			a: Action{
				FromPiece:        Piece{PieceType: PieceKing, Owner: ColorWhite, XY: XY{4, 0}},
				ToXY:             XY{6, 0},
				IsCastle:         true,
				IsKingsideCastle: true,
			},
			s: "White kingside castles",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.s, tc.a.String())
	}
}
func TestPieceTypeString(t *testing.T) {
	var (
		ptQueen  PieceType = PieceQueen
		ptKing   PieceType = PieceKing
		ptBishop PieceType = PieceBishop
		ptKnight PieceType = PieceKnight
		ptRook   PieceType = PieceRook
		ptPawn   PieceType = PiecePawn
	)
	assert.Equal(t, "Queen", ptQueen.String())
	assert.Equal(t, "King", ptKing.String())
	assert.Equal(t, "Bishop", ptBishop.String())
	assert.Equal(t, "Knight", ptKnight.String())
	assert.Equal(t, "Rook", ptRook.String())
	assert.Equal(t, "Pawn", ptPawn.String())
}
func TestColorString(t *testing.T) {
	var (
		cBlack color = ColorBlack
		cWhite color = ColorWhite
	)
	assert.Equal(t, "Black", cBlack.String())
	assert.Equal(t, "White", cWhite.String())
}
func TestPieceString(t *testing.T) {
	testCases := []struct {
		p Piece
		s string
	}{
		{
			p: Piece{
				Owner:     ColorBlack,
				PieceType: PieceQueen,
				XY:        XY{2, 2},
			},
			s: "Black's Queen at c6",
		},
		{
			p: Piece{
				Owner:     ColorWhite,
				PieceType: PieceKing,
				XY:        XY{3, 3},
			},
			s: "White's King at d5",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.s, tc.p.String())
	}
}
