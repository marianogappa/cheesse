package api

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
	g, err := newGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	require.NoError(t, err)
	assert.Equal(t, expected, g.String())
}

func TestActionString(t *testing.T) {
	testCases := []struct {
		a action
		s string
	}{
		{
			a: action{
				fromPiece: piece{pieceType: pieceQueen, owner: colorBlack, xy: xy{2, 2}},
				toXY:      xy{3, 3},
			},
			s: "Black's Queen at c6 moves to d5",
		},
		{
			a: action{
				fromPiece: piece{owner: colorWhite},
				isResign:  true,
			},
			s: "White resigns",
		},
		{
			a: action{
				fromPiece:     piece{pieceType: pieceBishop, owner: colorBlack, xy: xy{1, 1}},
				toXY:          xy{4, 4},
				isCapture:     true,
				capturedPiece: piece{pieceType: piecePawn, owner: colorWhite, xy: xy{4, 4}},
			},
			s: "Black's Bishop at b7 captures White's Pawn at e4",
		},
		{
			a: action{
				fromPiece:          piece{pieceType: piecePawn, owner: colorWhite, xy: xy{4, 1}},
				toXY:               xy{4, 0},
				isPromotion:        true,
				promotionPieceType: pieceKnight,
			},
			s: "White's Pawn at e7 promotes to Knight",
		},
		{
			a: action{
				fromPiece:   piece{pieceType: piecePawn, owner: colorBlack, xy: xy{5, 1}},
				toXY:        xy{5, 3},
				isEnPassant: true,
			},
			s: "Black's Pawn at f7 does en passant",
		},
		{
			a: action{
				fromPiece:          piece{pieceType: piecePawn, owner: colorBlack, xy: xy{3, 5}},
				toXY:               xy{4, 6},
				isCapture:          true,
				isEnPassantCapture: true,
				capturedPiece:      piece{pieceType: piecePawn, owner: colorWhite, xy: xy{4, 5}},
			},
			s: "Black's Pawn at d3 captures White's Pawn at e3 which was doing en passant",
		},
		{
			a: action{
				fromPiece:         piece{pieceType: pieceKing, owner: colorBlack, xy: xy{4, 0}},
				toXY:              xy{2, 0},
				isCastle:          true,
				isQueensideCastle: true,
			},
			s: "Black queenside castles",
		},
		{
			a: action{
				fromPiece:        piece{pieceType: pieceKing, owner: colorWhite, xy: xy{4, 0}},
				toXY:             xy{6, 0},
				isCastle:         true,
				isKingsideCastle: true,
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
		ptQueen  pieceType = pieceQueen
		ptKing   pieceType = pieceKing
		ptBishop pieceType = pieceBishop
		ptKnight pieceType = pieceKnight
		ptRook   pieceType = pieceRook
		ptPawn   pieceType = piecePawn
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
		cBlack color = colorBlack
		cWhite color = colorWhite
	)
	assert.Equal(t, "Black", cBlack.String())
	assert.Equal(t, "White", cWhite.String())
}
func TestPieceString(t *testing.T) {
	testCases := []struct {
		p piece
		s string
	}{
		{
			p: piece{
				owner:     colorBlack,
				pieceType: pieceQueen,
				xy:        xy{2, 2},
			},
			s: "Black's Queen at c6",
		},
		{
			p: piece{
				owner:     colorWhite,
				pieceType: pieceKing,
				xy:        xy{3, 3},
			},
			s: "White's King at d5",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.s, tc.p.String())
	}
}
