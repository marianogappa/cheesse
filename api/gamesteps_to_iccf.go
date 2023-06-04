package api

import "fmt"

func gamestepsToICCF(gamesteps []gameStep) []gameStep {
	iccfGameSteps := []gameStep{}
	for _, gamestep := range gamesteps {
		gamestep.s = gamestepToICCF(gamestep)
		iccfGameSteps = append(iccfGameSteps, gamestep)
	}
	return iccfGameSteps
}

func gamestepToICCF(gamestep gameStep) string {
	iccfMove := fmt.Sprintf(
		"%v%v",
		gamestep.a.fromPiece.xy.toICCF(),
		gamestep.a.toXY.toICCF(),
	)
	// For promotion, a fifth digit is added to the move's notation: 1 for queen, 2 for rook, 3 for bishop, and 4 for knight.
	if gamestep.a.isPromotion {
		switch gamestep.a.promotionPieceType {
		case pieceQueen:
			iccfMove += "1"
		case pieceRook:
			iccfMove += "2"
		case pieceBishop:
			iccfMove += "3"
		case pieceKnight:
			iccfMove += "4"
		}
	}
	return iccfMove
}
