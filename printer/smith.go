package printer

import (
	"fmt"
	"strings"

	"github.com/marianogappa/cheesse/core"
)

type SmithPrinter struct{}

func (p SmithPrinter) PrintGame(gameSteps []core.GameStep, gameCharacteristics GameCharacteristics) ([]string, error) {
	return genericGamePrinter(gameSteps, gameCharacteristics, p)
}

func (p SmithPrinter) PrintAction(gameStep core.GameStep, gameCharacteristics GameCharacteristics) (string, error) {
	var (
		capture   = ""
		promotion = ""
	)
	if gameStep.StepAction.IsCapture {
		capture = gameStep.StepAction.CapturedPiece.PieceType.ToSmith()
		if gameStep.StepAction.IsEnPassant {
			capture = "E"
		}
	}
	if gameStep.StepAction.IsKingsideCastle {
		capture = "c"
	}
	if gameStep.StepAction.IsQueensideCastle {
		capture = "C"
	}
	if gameStep.StepAction.IsPromotion {
		promotion = strings.ToUpper(gameStep.StepAction.PromotionPieceType.ToSmith())
	}
	return fmt.Sprintf(
		"%v%v%v%v",
		gameStep.StepAction.FromPiece.XY.ToAlgebraic(),
		gameStep.StepAction.ToXY.ToAlgebraic(),
		capture,
		promotion,
	), nil
}
