package core

import "math/bits"

// Bitboard square mapping: sq = y*8 + x, where y=0 is rank 8 (matching the XY convention).

func sqOf(xy XY) int { return xy.Y*8 + xy.X }

func xyOfSq(sq int) XY { return XY{X: sq & 7, Y: sq >> 3} }

func sqBit(sq int) uint64 { return 1 << uint(sq) }

const (
	dirN = iota
	dirNE
	dirE
	dirSE
	dirS
	dirSW
	dirW
	dirNW
)

var (
	rayAttacks    [8][64]uint64
	knightAttacks [64]uint64
	kingAttacks   [64]uint64
	// pawnCaptureAttacks[c][sq] is the set of squares a pawn of color c standing on sq attacks.
	// By symmetry, it's also the set of squares from which an opponent pawn attacks sq.
	pawnCaptureAttacks [2][64]uint64
	// betweenMasks[from][to] is the set of squares strictly between from and to, if they
	// share a rank, file or diagonal; zero otherwise.
	betweenMasks [64][64]uint64

	castleEmptyMasks      [2][2]uint64
	castleUnthreatenedSqs [2][2][3]int8
)

func init() {
	dirDeltas := [8]XY{dirN: {0, -1}, dirNE: {1, -1}, dirE: {1, 0}, dirSE: {1, 1}, dirS: {0, 1}, dirSW: {-1, 1}, dirW: {-1, 0}, dirNW: {-1, -1}}
	knightDeltas := [8]XY{{-1, -2}, {1, -2}, {-1, 2}, {1, 2}, {-2, -1}, {-2, 1}, {2, -1}, {2, 1}}
	for sq := 0; sq < 64; sq++ {
		from := xyOfSq(sq)
		for dir, delta := range dirDeltas {
			for xy := from.add(delta); isInBounds(xy); xy = xy.add(delta) {
				rayAttacks[dir][sq] |= sqBit(sqOf(xy))
			}
		}
		for _, delta := range knightDeltas {
			if xy := from.add(delta); isInBounds(xy) {
				knightAttacks[sq] |= sqBit(sqOf(xy))
			}
		}
		for _, delta := range dirDeltas {
			if xy := from.add(delta); isInBounds(xy) {
				kingAttacks[sq] |= sqBit(sqOf(xy))
			}
		}
		for _, dx := range [2]int{-1, 1} {
			if xy := from.add(XY{dx, 1}); isInBounds(xy) {
				pawnCaptureAttacks[ColorBlack][sq] |= sqBit(sqOf(xy))
			}
			if xy := from.add(XY{dx, -1}); isInBounds(xy) {
				pawnCaptureAttacks[ColorWhite][sq] |= sqBit(sqOf(xy))
			}
		}
		for _, delta := range dirDeltas {
			between := uint64(0)
			for xy := from.add(delta); isInBounds(xy); xy = xy.add(delta) {
				betweenMasks[sq][sqOf(xy)] = between
				between |= sqBit(sqOf(xy))
			}
		}
	}

	castleData := []struct {
		c            color
		ct           castleType
		empty        []XY
		unthreatened [3]XY
	}{
		{ColorBlack, castleTypeQueenside, []XY{{1, 0}, {2, 0}, {3, 0}}, [3]XY{{2, 0}, {3, 0}, {4, 0}}},
		{ColorBlack, castleTypeKingside, []XY{{5, 0}, {6, 0}}, [3]XY{{4, 0}, {5, 0}, {6, 0}}},
		{ColorWhite, castleTypeQueenside, []XY{{1, 7}, {2, 7}, {3, 7}}, [3]XY{{2, 7}, {3, 7}, {4, 7}}},
		{ColorWhite, castleTypeKingside, []XY{{5, 7}, {6, 7}}, [3]XY{{4, 7}, {5, 7}, {6, 7}}},
	}
	for _, cd := range castleData {
		for _, xy := range cd.empty {
			castleEmptyMasks[cd.c][cd.ct] |= sqBit(sqOf(xy))
		}
		for i, xy := range cd.unthreatened {
			castleUnthreatenedSqs[cd.c][cd.ct][i] = int8(sqOf(xy))
		}
	}
}

func positiveRayAttack(dir, sq int, occ uint64) uint64 {
	attacks := rayAttacks[dir][sq]
	if blockers := attacks & occ; blockers != 0 {
		attacks ^= rayAttacks[dir][bits.TrailingZeros64(blockers)]
	}
	return attacks
}

func negativeRayAttack(dir, sq int, occ uint64) uint64 {
	attacks := rayAttacks[dir][sq]
	if blockers := attacks & occ; blockers != 0 {
		attacks ^= rayAttacks[dir][63-bits.LeadingZeros64(blockers)]
	}
	return attacks
}

func rookAttacks(sq int, occ uint64) uint64 {
	return positiveRayAttack(dirS, sq, occ) | positiveRayAttack(dirE, sq, occ) |
		negativeRayAttack(dirN, sq, occ) | negativeRayAttack(dirW, sq, occ)
}

func bishopAttacks(sq int, occ uint64) uint64 {
	return positiveRayAttack(dirSE, sq, occ) | positiveRayAttack(dirSW, sq, occ) |
		negativeRayAttack(dirNE, sq, occ) | negativeRayAttack(dirNW, sq, occ)
}

func queenAttacks(sq int, occ uint64) uint64 {
	return rookAttacks(sq, occ) | bishopAttacks(sq, occ)
}

// squares[sq] encodes the piece at sq as pieceType<<1|color, or 0 if empty.

func (g *Game) setSq(c color, t PieceType, sq int) {
	b := sqBit(sq)
	g.bb[c][t] |= b
	g.occ[c] |= b
	g.squares[sq] = uint8(t)<<1 | uint8(c)
	if t == PieceKing {
		g.kingSq[c] = int8(sq)
	}
}

func (g *Game) clearSq(c color, t PieceType, sq int) {
	b := sqBit(sq)
	g.bb[c][t] &^= b
	g.occ[c] &^= b
	g.squares[sq] = 0
}

func (g Game) pieceAtSq(sq int) Piece {
	v := g.squares[sq]
	if v == 0 {
		return Piece{}
	}
	return Piece{PieceType: PieceType(v >> 1), Owner: color(v & 1), XY: xyOfSq(sq)}
}

func (g Game) hasPieceAt(c color, t PieceType, xy XY) bool {
	return g.bb[c][t]&sqBit(sqOf(xy)) != 0
}

func (g Game) occAll() uint64 {
	return g.occ[ColorBlack] | g.occ[ColorWhite]
}
