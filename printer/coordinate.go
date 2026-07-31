package printer

import (
	"fmt"

	"github.com/marianogappa/cheesse/core"
)

type CoordinatePrinter struct{}

func (p CoordinatePrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	return genericGamePrinter(gameSteps, gameCharacteristics, p)
}

func (p CoordinatePrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	gameCharacteristics = applyDefaultGameCharacteristics(gameCharacteristics)
	if gameStep.StepAction.IsCastle {
		return algCastle(gameStep, gameCharacteristics) + coordCheck(gameStep, gameCharacteristics), nil
	}
	if gameStep.StepAction.IsResign || gameStep.StepAction.IsDraw {
		return algResign(gameStep, gameCharacteristics), nil
	}
	delimiter := "-"
	if gameStep.StepAction.IsCapture {
		delimiter = "x"
	}
	promotion := ""
	if gameStep.StepAction.IsPromotion {
		promotion = fmt.Sprintf("=%v", gameStep.StepAction.PromotionPieceType.ToAlgebraic())
	}
	return fmt.Sprintf(
		"%v%v%v%v%v",
		gameStep.StepAction.FromPiece.XY.ToAlgebraic(),
		delimiter,
		gameStep.StepAction.ToXY.ToAlgebraic(),
		promotion,
		coordCheck(gameStep, gameCharacteristics),
	), nil
}

func coordCheck(gameStep core.GameStep, gameCharacteristics GameCharacteristics) string {
	return algCheck(gameStep, gameCharacteristics) + algCheckmate(gameStep, gameCharacteristics)
}
