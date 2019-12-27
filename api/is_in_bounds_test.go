package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInBounds(t *testing.T) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, isInBounds(xy{x, y}))
		}
	}
	assert.False(t, isInBounds(xy{-1, 0}))
	assert.False(t, isInBounds(xy{0, -1}))
	assert.False(t, isInBounds(xy{8, 0}))
	assert.False(t, isInBounds(xy{0, 8}))
}

func TestPieceIsInBounds(t *testing.T) {
	for _, c := range []color{colorBlack, colorWhite} {
		for _, pt := range []pieceType{pieceQueen, pieceKing, pieceBishop, pieceKnight, pieceRook} {
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					assert.True(t, piece{pieceType: pt, owner: c, xy: xy{x, y}}.isInBounds(xy{x, y}))
				}
			}
		}
	}
	for y := 1; y < 8; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, piece{pieceType: piecePawn, owner: colorBlack, xy: xy{x, y}}.isInBounds(xy{x, y}))
		}
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			assert.True(t, piece{pieceType: piecePawn, owner: colorWhite, xy: xy{x, y}}.isInBounds(xy{x, y}))
		}
	}
	for x := 0; x < 8; x++ {
		assert.False(t, piece{pieceType: piecePawn, owner: colorBlack, xy: xy{x, 0}}.isInBounds(xy{x, 0}))
	}
	for y := 0; y < 7; y++ {
		for x := 0; x < 8; x++ {
			assert.False(t, piece{pieceType: piecePawn, owner: colorWhite, xy: xy{x, 7}}.isInBounds(xy{x, 7}))
		}
	}
}
