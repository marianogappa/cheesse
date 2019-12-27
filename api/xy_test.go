package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXYAdd(t *testing.T) {
	assert.Equal(t, xy{1, 1}, xy{0, 0}.add(xy{1, 1}))
	assert.Equal(t, xy{1, 2}, xy{0, 1}.add(xy{1, 1}))
	assert.Equal(t, xy{2, 1}, xy{1, 0}.add(xy{1, 1}))
	assert.Equal(t, xy{}, xy{}.add(xy{}))
	assert.Equal(t, xy{3, 3}, xy{1, 2}.add(xy{2, 1}))
}

func TestXYEq(t *testing.T) {
	assert.True(t, xy{}.eq(xy{}))
	assert.True(t, xy{1, 2}.eq(xy{1, 2}))
	assert.False(t, xy{1, 2}.eq(xy{3, 4}))
}

func TestXYDeltaTowards(t *testing.T) {
	assert.Equal(t, xy{0, 0}, xy{0, 0}.deltaTowards(xy{0, 0}))

	assert.Equal(t, xy{1, 0}, xy{0, 0}.deltaTowards(xy{1, 0}))
	assert.Equal(t, xy{0, 1}, xy{0, 0}.deltaTowards(xy{0, 1}))
	assert.Equal(t, xy{-1, 0}, xy{0, 0}.deltaTowards(xy{-1, 0}))
	assert.Equal(t, xy{0, -1}, xy{0, 0}.deltaTowards(xy{0, -1}))

	assert.Equal(t, xy{-1, 0}, xy{1, 0}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{0, -1}, xy{0, 1}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{1, 0}, xy{-1, 0}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{0, 1}, xy{0, -1}.deltaTowards(xy{0, 0}))

	assert.Equal(t, xy{1, 0}, xy{0, 0}.deltaTowards(xy{2, 0}))
	assert.Equal(t, xy{0, 1}, xy{0, 0}.deltaTowards(xy{0, 2}))
	assert.Equal(t, xy{-1, 0}, xy{0, 0}.deltaTowards(xy{-2, 0}))
	assert.Equal(t, xy{0, -1}, xy{0, 0}.deltaTowards(xy{0, -2}))

	assert.Equal(t, xy{-1, 0}, xy{2, 0}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{0, -1}, xy{0, 2}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{1, 0}, xy{-2, 0}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{0, 1}, xy{0, -2}.deltaTowards(xy{0, 0}))

	assert.Equal(t, xy{1, 1}, xy{0, 0}.deltaTowards(xy{1, 1}))
	assert.Equal(t, xy{1, 1}, xy{0, 0}.deltaTowards(xy{2, 2}))
	assert.Equal(t, xy{-1, -1}, xy{1, 1}.deltaTowards(xy{0, 0}))
	assert.Equal(t, xy{-1, -1}, xy{2, 2}.deltaTowards(xy{0, 0}))

	assert.Equal(t, xy{1, -1}, xy{0, 1}.deltaTowards(xy{1, 0}))
	assert.Equal(t, xy{-1, 1}, xy{1, 0}.deltaTowards(xy{0, 1}))
	assert.Equal(t, xy{1, -1}, xy{0, 2}.deltaTowards(xy{2, 0}))
	assert.Equal(t, xy{-1, 1}, xy{2, 0}.deltaTowards(xy{0, 2}))
}

func TestXYTowards(t *testing.T) {
	// Pieces that don't multi move
	assert.Equal(t, []xy{}, piece{pieceNone, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{}, piece{piecePawn, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{}, piece{pieceKing, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{}, piece{pieceKnight, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))

	for _, pType := range []pieceType{pieceQueen, pieceBishop, pieceRook} {
		// Chebyshev distance <= 1 means there are no squares in between piece and destination square
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 0}}.xysTowards(xy{0, 0}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 0}}.xysTowards(xy{1, 0}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 0}}.xysTowards(xy{0, 1}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 0}}.xysTowards(xy{1, 1}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 0}}.xysTowards(xy{0, 0}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{1, 0}}.xysTowards(xy{0, 0}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{0, 1}}.xysTowards(xy{0, 0}))
		assert.Equal(t, []xy{}, piece{pType, colorBlack, xy{1, 1}}.xysTowards(xy{0, 0}))
	}

	// Impossible XYs given piece types
	assert.Equal(t, []xy{}, piece{pieceRook, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{}, piece{pieceBishop, colorBlack, xy{0, 0}}.xysTowards(xy{0, 3}))
	assert.Equal(t, []xy{}, piece{pieceQueen, colorBlack, xy{0, 0}}.xysTowards(xy{1, 2}))

	// Valid xysTowards
	assert.Equal(t, []xy{{1, 1}, {2, 2}}, piece{pieceBishop, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{{0, 1}, {0, 2}}, piece{pieceRook, colorBlack, xy{0, 0}}.xysTowards(xy{0, 3}))
	assert.Equal(t, []xy{{1, 0}, {2, 0}}, piece{pieceRook, colorBlack, xy{0, 0}}.xysTowards(xy{3, 0}))
	assert.Equal(t, []xy{{1, 1}, {2, 2}}, piece{pieceQueen, colorBlack, xy{0, 0}}.xysTowards(xy{3, 3}))
	assert.Equal(t, []xy{{0, 1}, {0, 2}}, piece{pieceQueen, colorBlack, xy{0, 0}}.xysTowards(xy{0, 3}))
	assert.Equal(t, []xy{{1, 0}, {2, 0}}, piece{pieceQueen, colorBlack, xy{0, 0}}.xysTowards(xy{3, 0}))
}

func TestXYToAlgebraic(t *testing.T) {
	assert.Equal(t, "a8", xy{0, 0}.toAlgebraic())
	assert.Equal(t, "b8", xy{1, 0}.toAlgebraic())
	assert.Equal(t, "c8", xy{2, 0}.toAlgebraic())
	assert.Equal(t, "d8", xy{3, 0}.toAlgebraic())
	assert.Equal(t, "e8", xy{4, 0}.toAlgebraic())
	assert.Equal(t, "f8", xy{5, 0}.toAlgebraic())
	assert.Equal(t, "g8", xy{6, 0}.toAlgebraic())
	assert.Equal(t, "h8", xy{7, 0}.toAlgebraic())
	assert.Equal(t, "a7", xy{0, 1}.toAlgebraic())
	assert.Equal(t, "b7", xy{1, 1}.toAlgebraic())
	assert.Equal(t, "c7", xy{2, 1}.toAlgebraic())
	assert.Equal(t, "d7", xy{3, 1}.toAlgebraic())
	assert.Equal(t, "e7", xy{4, 1}.toAlgebraic())
	assert.Equal(t, "f7", xy{5, 1}.toAlgebraic())
	assert.Equal(t, "g7", xy{6, 1}.toAlgebraic())
	assert.Equal(t, "h7", xy{7, 1}.toAlgebraic())
	assert.Equal(t, "a6", xy{0, 2}.toAlgebraic())
	assert.Equal(t, "b6", xy{1, 2}.toAlgebraic())
	assert.Equal(t, "c6", xy{2, 2}.toAlgebraic())
	assert.Equal(t, "d6", xy{3, 2}.toAlgebraic())
	assert.Equal(t, "e6", xy{4, 2}.toAlgebraic())
	assert.Equal(t, "f6", xy{5, 2}.toAlgebraic())
	assert.Equal(t, "g6", xy{6, 2}.toAlgebraic())
	assert.Equal(t, "h6", xy{7, 2}.toAlgebraic())
	assert.Equal(t, "a5", xy{0, 3}.toAlgebraic())
	assert.Equal(t, "b5", xy{1, 3}.toAlgebraic())
	assert.Equal(t, "c5", xy{2, 3}.toAlgebraic())
	assert.Equal(t, "d5", xy{3, 3}.toAlgebraic())
	assert.Equal(t, "e5", xy{4, 3}.toAlgebraic())
	assert.Equal(t, "f5", xy{5, 3}.toAlgebraic())
	assert.Equal(t, "g5", xy{6, 3}.toAlgebraic())
	assert.Equal(t, "h5", xy{7, 3}.toAlgebraic())
	assert.Equal(t, "a4", xy{0, 4}.toAlgebraic())
	assert.Equal(t, "b4", xy{1, 4}.toAlgebraic())
	assert.Equal(t, "c4", xy{2, 4}.toAlgebraic())
	assert.Equal(t, "d4", xy{3, 4}.toAlgebraic())
	assert.Equal(t, "e4", xy{4, 4}.toAlgebraic())
	assert.Equal(t, "f4", xy{5, 4}.toAlgebraic())
	assert.Equal(t, "g4", xy{6, 4}.toAlgebraic())
	assert.Equal(t, "h4", xy{7, 4}.toAlgebraic())
	assert.Equal(t, "a3", xy{0, 5}.toAlgebraic())
	assert.Equal(t, "b3", xy{1, 5}.toAlgebraic())
	assert.Equal(t, "c3", xy{2, 5}.toAlgebraic())
	assert.Equal(t, "d3", xy{3, 5}.toAlgebraic())
	assert.Equal(t, "e3", xy{4, 5}.toAlgebraic())
	assert.Equal(t, "f3", xy{5, 5}.toAlgebraic())
	assert.Equal(t, "g3", xy{6, 5}.toAlgebraic())
	assert.Equal(t, "h3", xy{7, 5}.toAlgebraic())
	assert.Equal(t, "a2", xy{0, 6}.toAlgebraic())
	assert.Equal(t, "b2", xy{1, 6}.toAlgebraic())
	assert.Equal(t, "c2", xy{2, 6}.toAlgebraic())
	assert.Equal(t, "d2", xy{3, 6}.toAlgebraic())
	assert.Equal(t, "e2", xy{4, 6}.toAlgebraic())
	assert.Equal(t, "f2", xy{5, 6}.toAlgebraic())
	assert.Equal(t, "g2", xy{6, 6}.toAlgebraic())
	assert.Equal(t, "h2", xy{7, 6}.toAlgebraic())
	assert.Equal(t, "a1", xy{0, 7}.toAlgebraic())
	assert.Equal(t, "b1", xy{1, 7}.toAlgebraic())
	assert.Equal(t, "c1", xy{2, 7}.toAlgebraic())
	assert.Equal(t, "d1", xy{3, 7}.toAlgebraic())
	assert.Equal(t, "e1", xy{4, 7}.toAlgebraic())
	assert.Equal(t, "f1", xy{5, 7}.toAlgebraic())
	assert.Equal(t, "g1", xy{6, 7}.toAlgebraic())
	assert.Equal(t, "h1", xy{7, 7}.toAlgebraic())
}
