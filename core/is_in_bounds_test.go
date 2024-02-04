package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInBounds(t *testing.T) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, isInBounds(XY{x, y}))
		}
	}
	assert.False(t, isInBounds(XY{-1, 0}))
	assert.False(t, isInBounds(XY{0, -1}))
	assert.False(t, isInBounds(XY{8, 0}))
	assert.False(t, isInBounds(XY{0, 8}))
}

func TestPieceIsInBounds(t *testing.T) {
	for _, c := range []color{ColorBlack, ColorWhite} {
		for _, pt := range []PieceType{PieceQueen, PieceKing, PieceBishop, PieceKnight, PieceRook} {
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					assert.True(t, Piece{PieceType: pt, Owner: c, XY: XY{x, y}}.isInBounds(XY{x, y}))
				}
			}
		}
	}
	for y := 1; y < 8; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{x, y}}.isInBounds(XY{x, y}))
		}
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{x, y}}.isInBounds(XY{x, y}))
		}
	}
	for x := 0; x < 8; x++ {
		assert.False(t, Piece{PieceType: PiecePawn, Owner: ColorBlack, XY: XY{x, 0}}.isInBounds(XY{x, 0}))
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			assert.False(t, Piece{PieceType: PiecePawn, Owner: ColorWhite, XY: XY{x, 7}}.isInBounds(XY{x, 7}))
		}
	}
}
