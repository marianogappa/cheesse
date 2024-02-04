package printer

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

type ICCFPrinter struct{}

func (p ICCFPrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	return genericGamePrinter(gameSteps, gameCharacteristics, p)
}

func (p ICCFPrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	promotion := ""
	if gameStep.StepAction.IsPromotion {
		promotion = gameStep.StepAction.PromotionPieceType.ToICCF()
	}
	return fmt.Sprintf("%v%v%v", gameStep.StepAction.FromPiece.XY.ToICCF(), gameStep.StepAction.ToXY.ToICCF(), promotion), nil
}
