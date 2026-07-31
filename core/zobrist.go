package core

import "math/bits"

// Zobrist hashing: a position (placement + turn + castling rights + e.p. target)
// maps to a uint64 by xoring per-feature random keys. Used to detect repetitions.
// https://www.chessprogramming.org/Zobrist_Hashing

var (
	zobristPieces        [2][7][64]uint64
	zobristWhiteTurn     uint64
	zobristCastling      [4]uint64 // WK, WQ, BK, BQ
	zobristEnPassantFile [8]uint64
)

func init() {
	// SplitMix64 with a fixed seed: keys must be deterministic across runs so that
	// hashes seeded from previous positions (e.g. via the API) remain comparable.
	state := uint64(0x9E3779B97F4A7C15)
	next := func() uint64 {
		state += 0x9E3779B97F4A7C15
		z := state
		z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
		z = (z ^ (z >> 27)) * 0x94D049BB133111EB
		return z ^ (z >> 31)
	}
	for c := 0; c < 2; c++ {
		for t := PieceQueen; t <= PiecePawn; t++ {
			for sq := 0; sq < 64; sq++ {
				zobristPieces[c][t][sq] = next()
			}
		}
	}
	zobristWhiteTurn = next()
	for i := range zobristCastling {
		zobristCastling[i] = next()
	}
	for i := range zobristEnPassantFile {
		zobristEnPassantFile[i] = next()
	}
}

// positionHash returns the Zobrist hash of the position for repetition purposes:
// piece placement, side to move, castling rights and en passant target square.
func (g Game) positionHash() uint64 {
	h := uint64(0)
	for c := 0; c < 2; c++ {
		for occ := g.occ[c]; occ != 0; occ &= occ - 1 {
			sq := bits.TrailingZeros64(occ)
			h ^= zobristPieces[c][g.squares[sq]>>1][sq]
		}
	}
	if g.Turn() == ColorWhite {
		h ^= zobristWhiteTurn
	}
	if g.CanWhiteKingsideCastle {
		h ^= zobristCastling[0]
	}
	if g.CanWhiteQueensideCastle {
		h ^= zobristCastling[1]
	}
	if g.CanBlackKingsideCastle {
		h ^= zobristCastling[2]
	}
	if g.CanBlackQueensideCastle {
		h ^= zobristCastling[3]
	}
	if g.IsLastMoveEnPassant {
		h ^= zobristEnPassantFile[g.EnPassantTargetSquare.X]
	}
	return h
}
