package ai

import (
	"github.com/marianogappa/cheesse/core"
)

// Piece values in centipawns (standard simplified-evaluation values).
func pieceValue(pt core.PieceType) int {
	switch pt {
	case core.PieceQueen:
		return 900
	case core.PieceRook:
		return 500
	case core.PieceBishop:
		return 330
	case core.PieceKnight:
		return 320
	case core.PiecePawn:
		return 100
	case core.PieceKing:
		return 20000
	}
	return 0
}

// Piece-square tables from the "simplified evaluation function" (chessprogramming
// wiki). Tables are laid out from White's perspective with rank 8 first, which
// matches this codebase's square indexing (sq = Y*8 + X, where Y=0 is rank 8).
// Black uses the vertically mirrored square.
var pstPawn = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pstKnight = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var pstBishop = [64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var pstRook = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var pstQueen = [64]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

var pstKingMidgame = [64]int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}

var pstKingEndgame = [64]int{
	-50, -40, -30, -20, -20, -30, -40, -50,
	-30, -20, -10, 0, 0, -10, -20, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 30, 40, 40, 30, -10, -30,
	-30, -10, 20, 30, 30, 20, -10, -30,
	-30, -30, 0, 0, 0, 0, -30, -30,
	-50, -30, -30, -30, -30, -30, -30, -50,
}

// endgameMaterialThreshold: when both sides' non-pawn, non-king material is at
// or below roughly a rook and a minor piece each, use the endgame king table.
const endgameMaterialThreshold = 2 * (500 + 330)

// evaluate returns the static score of the position in centipawns from the
// perspective of the side to move. Terminal positions (mate/draw) are handled
// by the search, not here.
func evaluate(g core.Game) int {
	whiteScore, whiteNonPawnMaterial := sideScore(g, core.ColorWhite)
	blackScore, blackNonPawnMaterial := sideScore(g, core.ColorBlack)

	endgame := whiteNonPawnMaterial+blackNonPawnMaterial <= endgameMaterialThreshold
	whiteScore += kingSquareScore(g, core.ColorWhite, endgame)
	blackScore += kingSquareScore(g, core.ColorBlack, endgame)

	score := whiteScore - blackScore
	if g.Turn() == core.ColorBlack {
		return -score
	}
	return score
}

// sideScore returns material + piece-square score for all non-king pieces of
// the given color, plus the side's non-pawn material (for endgame detection).
func sideScore(g core.Game, c core.Color) (score, nonPawnMaterial int) {
	for _, piece := range g.Pieces(c) {
		if piece.PieceType == core.PieceKing {
			continue
		}
		value := pieceValue(piece.PieceType)
		score += value + pstScore(piece.PieceType, c, piece.XY)
		if piece.PieceType != core.PiecePawn {
			nonPawnMaterial += value
		}
	}
	return score, nonPawnMaterial
}

func kingSquareScore(g core.Game, c core.Color, endgame bool) int {
	sq := pstIndex(c, g.King(c).XY)
	if endgame {
		return pstKingEndgame[sq]
	}
	return pstKingMidgame[sq]
}

func pstScore(pt core.PieceType, c core.Color, xy core.XY) int {
	sq := pstIndex(c, xy)
	switch pt {
	case core.PiecePawn:
		return pstPawn[sq]
	case core.PieceKnight:
		return pstKnight[sq]
	case core.PieceBishop:
		return pstBishop[sq]
	case core.PieceRook:
		return pstRook[sq]
	case core.PieceQueen:
		return pstQueen[sq]
	}
	return 0
}

// pstIndex maps a board coordinate to a piece-square table index. Tables are
// stored from White's perspective (Y=0 is rank 8); Black mirrors vertically.
func pstIndex(c core.Color, xy core.XY) int {
	if c == core.ColorWhite {
		return xy.Y*8 + xy.X
	}
	return (7-xy.Y)*8 + xy.X
}
