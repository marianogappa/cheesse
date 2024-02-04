package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXYAdd(t *testing.T) {
	assert.Equal(t, XY{1, 1}, XY{0, 0}.add(XY{1, 1}))
	assert.Equal(t, XY{1, 2}, XY{0, 1}.add(XY{1, 1}))
	assert.Equal(t, XY{2, 1}, XY{1, 0}.add(XY{1, 1}))
	assert.Equal(t, XY{}, XY{}.add(XY{}))
	assert.Equal(t, XY{3, 3}, XY{1, 2}.add(XY{2, 1}))
}

func TestXYEq(t *testing.T) {
	assert.True(t, XY{}.eq(XY{}))
	assert.True(t, XY{1, 2}.eq(XY{1, 2}))
	assert.False(t, XY{1, 2}.eq(XY{3, 4}))
}

func TestXYDeltaTowards(t *testing.T) {
	assert.Equal(t, XY{0, 0}, XY{0, 0}.deltaTowards(XY{0, 0}))

	assert.Equal(t, XY{1, 0}, XY{0, 0}.deltaTowards(XY{1, 0}))
	assert.Equal(t, XY{0, 1}, XY{0, 0}.deltaTowards(XY{0, 1}))
	assert.Equal(t, XY{-1, 0}, XY{0, 0}.deltaTowards(XY{-1, 0}))
	assert.Equal(t, XY{0, -1}, XY{0, 0}.deltaTowards(XY{0, -1}))

	assert.Equal(t, XY{-1, 0}, XY{1, 0}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{0, -1}, XY{0, 1}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{1, 0}, XY{-1, 0}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{0, 1}, XY{0, -1}.deltaTowards(XY{0, 0}))

	assert.Equal(t, XY{1, 0}, XY{0, 0}.deltaTowards(XY{2, 0}))
	assert.Equal(t, XY{0, 1}, XY{0, 0}.deltaTowards(XY{0, 2}))
	assert.Equal(t, XY{-1, 0}, XY{0, 0}.deltaTowards(XY{-2, 0}))
	assert.Equal(t, XY{0, -1}, XY{0, 0}.deltaTowards(XY{0, -2}))

	assert.Equal(t, XY{-1, 0}, XY{2, 0}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{0, -1}, XY{0, 2}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{1, 0}, XY{-2, 0}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{0, 1}, XY{0, -2}.deltaTowards(XY{0, 0}))

	assert.Equal(t, XY{1, 1}, XY{0, 0}.deltaTowards(XY{1, 1}))
	assert.Equal(t, XY{1, 1}, XY{0, 0}.deltaTowards(XY{2, 2}))
	assert.Equal(t, XY{-1, -1}, XY{1, 1}.deltaTowards(XY{0, 0}))
	assert.Equal(t, XY{-1, -1}, XY{2, 2}.deltaTowards(XY{0, 0}))

	assert.Equal(t, XY{1, -1}, XY{0, 1}.deltaTowards(XY{1, 0}))
	assert.Equal(t, XY{-1, 1}, XY{1, 0}.deltaTowards(XY{0, 1}))
	assert.Equal(t, XY{1, -1}, XY{0, 2}.deltaTowards(XY{2, 0}))
	assert.Equal(t, XY{-1, 1}, XY{2, 0}.deltaTowards(XY{0, 2}))
}

func TestXYTowards(t *testing.T) {
	// Pieces that don't multi move
	assert.Equal(t, []XY{}, Piece{PieceNone, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{}, Piece{PiecePawn, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{}, Piece{PieceKing, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{}, Piece{PieceKnight, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))

	for _, pType := range []PieceType{PieceQueen, PieceBishop, PieceRook} {
		// Chebyshev distance <= 1 means there are no squares in between piece and destination square
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 0}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 0}}.xysTowards(XY{1, 0}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 1}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 0}}.xysTowards(XY{1, 1}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 0}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{1, 0}}.xysTowards(XY{0, 0}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{0, 1}}.xysTowards(XY{0, 0}))
		assert.Equal(t, []XY{}, Piece{pType, ColorBlack, XY{1, 1}}.xysTowards(XY{0, 0}))
	}

	// Impossible XYs given piece types
	assert.Equal(t, []XY{}, Piece{PieceRook, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{}, Piece{PieceBishop, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 3}))
	assert.Equal(t, []XY{}, Piece{PieceQueen, ColorBlack, XY{0, 0}}.xysTowards(XY{1, 2}))

	// Valid xysTowards
	assert.Equal(t, []XY{{1, 1}, {2, 2}}, Piece{PieceBishop, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{{0, 1}, {0, 2}}, Piece{PieceRook, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 3}))
	assert.Equal(t, []XY{{1, 0}, {2, 0}}, Piece{PieceRook, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 0}))
	assert.Equal(t, []XY{{1, 1}, {2, 2}}, Piece{PieceQueen, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 3}))
	assert.Equal(t, []XY{{0, 1}, {0, 2}}, Piece{PieceQueen, ColorBlack, XY{0, 0}}.xysTowards(XY{0, 3}))
	assert.Equal(t, []XY{{1, 0}, {2, 0}}, Piece{PieceQueen, ColorBlack, XY{0, 0}}.xysTowards(XY{3, 0}))
}

func TestXYToAlgebraic(t *testing.T) {
	assert.Equal(t, "a8", XY{0, 0}.ToAlgebraic())
	assert.Equal(t, "b8", XY{1, 0}.ToAlgebraic())
	assert.Equal(t, "c8", XY{2, 0}.ToAlgebraic())
	assert.Equal(t, "d8", XY{3, 0}.ToAlgebraic())
	assert.Equal(t, "e8", XY{4, 0}.ToAlgebraic())
	assert.Equal(t, "f8", XY{5, 0}.ToAlgebraic())
	assert.Equal(t, "g8", XY{6, 0}.ToAlgebraic())
	assert.Equal(t, "h8", XY{7, 0}.ToAlgebraic())
	assert.Equal(t, "a7", XY{0, 1}.ToAlgebraic())
	assert.Equal(t, "b7", XY{1, 1}.ToAlgebraic())
	assert.Equal(t, "c7", XY{2, 1}.ToAlgebraic())
	assert.Equal(t, "d7", XY{3, 1}.ToAlgebraic())
	assert.Equal(t, "e7", XY{4, 1}.ToAlgebraic())
	assert.Equal(t, "f7", XY{5, 1}.ToAlgebraic())
	assert.Equal(t, "g7", XY{6, 1}.ToAlgebraic())
	assert.Equal(t, "h7", XY{7, 1}.ToAlgebraic())
	assert.Equal(t, "a6", XY{0, 2}.ToAlgebraic())
	assert.Equal(t, "b6", XY{1, 2}.ToAlgebraic())
	assert.Equal(t, "c6", XY{2, 2}.ToAlgebraic())
	assert.Equal(t, "d6", XY{3, 2}.ToAlgebraic())
	assert.Equal(t, "e6", XY{4, 2}.ToAlgebraic())
	assert.Equal(t, "f6", XY{5, 2}.ToAlgebraic())
	assert.Equal(t, "g6", XY{6, 2}.ToAlgebraic())
	assert.Equal(t, "h6", XY{7, 2}.ToAlgebraic())
	assert.Equal(t, "a5", XY{0, 3}.ToAlgebraic())
	assert.Equal(t, "b5", XY{1, 3}.ToAlgebraic())
	assert.Equal(t, "c5", XY{2, 3}.ToAlgebraic())
	assert.Equal(t, "d5", XY{3, 3}.ToAlgebraic())
	assert.Equal(t, "e5", XY{4, 3}.ToAlgebraic())
	assert.Equal(t, "f5", XY{5, 3}.ToAlgebraic())
	assert.Equal(t, "g5", XY{6, 3}.ToAlgebraic())
	assert.Equal(t, "h5", XY{7, 3}.ToAlgebraic())
	assert.Equal(t, "a4", XY{0, 4}.ToAlgebraic())
	assert.Equal(t, "b4", XY{1, 4}.ToAlgebraic())
	assert.Equal(t, "c4", XY{2, 4}.ToAlgebraic())
	assert.Equal(t, "d4", XY{3, 4}.ToAlgebraic())
	assert.Equal(t, "e4", XY{4, 4}.ToAlgebraic())
	assert.Equal(t, "f4", XY{5, 4}.ToAlgebraic())
	assert.Equal(t, "g4", XY{6, 4}.ToAlgebraic())
	assert.Equal(t, "h4", XY{7, 4}.ToAlgebraic())
	assert.Equal(t, "a3", XY{0, 5}.ToAlgebraic())
	assert.Equal(t, "b3", XY{1, 5}.ToAlgebraic())
	assert.Equal(t, "c3", XY{2, 5}.ToAlgebraic())
	assert.Equal(t, "d3", XY{3, 5}.ToAlgebraic())
	assert.Equal(t, "e3", XY{4, 5}.ToAlgebraic())
	assert.Equal(t, "f3", XY{5, 5}.ToAlgebraic())
	assert.Equal(t, "g3", XY{6, 5}.ToAlgebraic())
	assert.Equal(t, "h3", XY{7, 5}.ToAlgebraic())
	assert.Equal(t, "a2", XY{0, 6}.ToAlgebraic())
	assert.Equal(t, "b2", XY{1, 6}.ToAlgebraic())
	assert.Equal(t, "c2", XY{2, 6}.ToAlgebraic())
	assert.Equal(t, "d2", XY{3, 6}.ToAlgebraic())
	assert.Equal(t, "e2", XY{4, 6}.ToAlgebraic())
	assert.Equal(t, "f2", XY{5, 6}.ToAlgebraic())
	assert.Equal(t, "g2", XY{6, 6}.ToAlgebraic())
	assert.Equal(t, "h2", XY{7, 6}.ToAlgebraic())
	assert.Equal(t, "a1", XY{0, 7}.ToAlgebraic())
	assert.Equal(t, "b1", XY{1, 7}.ToAlgebraic())
	assert.Equal(t, "c1", XY{2, 7}.ToAlgebraic())
	assert.Equal(t, "d1", XY{3, 7}.ToAlgebraic())
	assert.Equal(t, "e1", XY{4, 7}.ToAlgebraic())
	assert.Equal(t, "f1", XY{5, 7}.ToAlgebraic())
	assert.Equal(t, "g1", XY{6, 7}.ToAlgebraic())
	assert.Equal(t, "h1", XY{7, 7}.ToAlgebraic())
}
